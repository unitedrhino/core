package timer

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/core/service/timed/internal/domain"
	"gitee.com/i-Things/core/service/timed/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/timed/timedschedulersvr/internal/svc"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

// 定时任务检查
func TimingTaskCheck(svcCtx *svc.ServiceContext) {
	logx.Debug("TimingTaskCheck run")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	//dgsvr 订阅到了设备端数据，此时调用StartSpan方法，将订阅到的主题推送给jaeger
	//此时的ctx已经包含当前节点的span信息，会随着 handle(ctx).Publish 传递到下个节点
	ctx, span := ctxs.StartSpan(ctx, "timedSchedulersvr.timing.taskCheck", "")
	defer span.End()
	err := func() error {
		jobDB := stores.WithNoDebug(ctx, relationDB.NewTaskInfoRepo)
		js, err := jobDB.FindByFilter(ctx, relationDB.TaskFilter{WithGroup: true,
			Status: []int64{def.StatusWaitStop, def.StatusWaitDelete, def.StatusWaitRun},
			Types:  []int64{domain.TaskTypeTiming}},
			&def.PageInfo{
				Orders: []def.OrderBy{{Filed: "priority", Sort: def.OrderDesc}},
			})
		if err != nil {
			return err
		}
		wait := sync.WaitGroup{}
		for _, j := range js {
			wait.Add(1)
			t := j
			utils.Go(ctx, func() {
				err := func() error {
					switch t.Status {
					case def.StatusWaitRun:
						return TimingTaskStatusRunCheck(ctx, svcCtx, &wait, t)
					case def.StatusWaitDelete, def.StatusWaitStop:
						return TimingTaskStatusStopCheck(ctx, svcCtx, &wait, t)
					}
					//其他状态不需要处理
					return nil
				}()
				if err != nil {
					logx.WithContext(ctx).Errorf("TimingTaskCheck.one  err:%+v , task:%+v", err, t)
				}
			})
		}
		wait.Wait()
		return nil
	}()
	if err != nil {
		logx.WithContext(ctx).Errorf("TimingTaskCheck  err:%v", err)
	}
}

// 需要检查任务是否启动,如果没有启动需要启动
func TimingTaskStatusRunCheck(ctx context.Context, svcCtx *svc.ServiceContext, wait *sync.WaitGroup, task *relationDB.TimedTaskInfo) error {
	defer wait.Done()
	taskCode := getTimingTaskCode(task)
	taskInfo := domain.TaskInfo{
		ID:     task.ID,
		Params: "",
	}
	payload, _ := json.Marshal(taskInfo)
	err := svcCtx.Scheduler.Register(task.CronExpr, taskCode, payload, asynq.Queue(domain.ToPriority(task.Priority)))
	if err != nil {
		logx.WithContext(ctx).Errorf("TimingTaskStatusRunCheck.Register err:%v task:%v", err, task)
		return errors.System.AddDetail(err)
	}
	jobDB := relationDB.NewTaskInfoRepo(ctx)
	task.Status = def.StatusRunning
	err = jobDB.Update(ctx, task)
	return err
}

// 如果处于运行状态需要停止
func TimingTaskStatusStopCheck(ctx context.Context, svcCtx *svc.ServiceContext, wait *sync.WaitGroup, task *relationDB.TimedTaskInfo) error {
	defer wait.Done()
	taskCode := getTimingTaskCode(task)
	err := svcCtx.Scheduler.Unregister(taskCode)
	if err != nil {
		logx.WithContext(ctx).Errorf("TimingTaskStatusStopCheck.Unregister err:%v task:%v", err, task)
		return errors.System.AddDetail(err)
	}
	jobDB := relationDB.NewTaskInfoRepo(ctx)
	switch task.Status {
	case def.StatusWaitDelete:
		err = jobDB.Delete(ctx, task.ID)
		if err != nil {
			return err
		}
	case def.StatusWaitStop:
		task.Status = def.StatusStopped
		err = jobDB.Update(ctx, task)
	}
	return err
}

func getTimingTaskCode(j *relationDB.TimedTaskInfo) string {
	return fmt.Sprintf("timing:%s:%s", j.GroupCode, j.Code)
}
