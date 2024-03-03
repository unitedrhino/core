package timer

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/core/service/timed/internal/domain"
	"gitee.com/i-Things/core/service/timed/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/internal/logic"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/internal/logic/processTask"
	timedmanagelogic "gitee.com/i-Things/core/service/timed/timedjobsvr/internal/logic/timedmanage"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type Timed struct {
	SvcCtx *svc.ServiceContext
}

func (t Timed) ProcessTask(ctx context.Context, Task *asynq.Task) error {
	go func() {
		defer utils.Recover(ctx)
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
		defer cancel()
		ctx, span := ctxs.StartSpan(ctx, "timedJob.Process", "")
		defer span.End()
		utils.Recover(ctx)
		err := func() error {
			var taskInfo domain.TaskInfo
			json.Unmarshal(Task.Payload(), &taskInfo)
			tr := stores.WithNoDebug(ctx, relationDB.NewTaskInfoRepo)
			task, err := tr.FindOneByFilter(ctx, relationDB.TaskFilter{
				IDs:       []int64{taskInfo.ID},
				WithGroup: true,
			})
			if err != nil {
				return err
			}
			if task.Type == domain.TaskTypeTiming && task.Status != def.StatusRunning { //如果没有处于运行中,任务不能执行
				return nil
			}
			err = logic.FillTaskInfoDo(&taskInfo, task)
			if err != nil {
				return err
			}
			logx.WithContext(ctx).Debug("timedJob Process task:%v", utils.Fmt(taskInfo))
			return processTask.NewProcessTask(ctx, t.SvcCtx, func(ctx context.Context, req *timedjob.TaskSendReq) error {
				_, err := timedmanagelogic.NewTaskSendLogic(ctx, t.SvcCtx).TaskSend(req)
				return err
			}).Process(ctx, taskInfo)
		}()
		if err != nil {
			logx.WithContext(ctx).Errorf("Process  task.Type:%v,task.Payload:%v err:%v", Task.Type(), string(Task.Payload()), err)
		}
	}()
	return nil
}
