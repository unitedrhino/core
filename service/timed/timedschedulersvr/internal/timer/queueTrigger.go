package timer

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/timed/internal/domain"
	"gitee.com/i-Things/core/service/timed/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/i-Things/core/service/timed/timedschedulersvr/internal/svc"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type QueueTrigger struct {
	Subs []*nats.Subscription
	Task *relationDB.TimedTaskInfo
}

var (
	QueueMutex sync.Mutex
	QueueMap   = map[string]*QueueTrigger{}
)

// 定时任务检查
func QueueTaskCheck(svcCtx *svc.ServiceContext) {
	logx.Debug("QueueTaskCheck run")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	//dgsvr 订阅到了设备端数据，此时调用StartSpan方法，将订阅到的主题推送给jaeger
	//此时的ctx已经包含当前节点的span信息，会随着 handle(ctx).Publish 传递到下个节点
	ctx, span := ctxs.StartSpan(ctx, "timedSchedulersvr.queue.taskCheck", "")
	defer span.End()
	err := func() error {
		jobDB := stores.WithNoDebug(ctx, relationDB.NewTaskInfoRepo)
		//jobDB := relationDB.NewTaskInfoRepo(ctx)
		js, err := jobDB.FindByFilter(ctx, relationDB.TaskFilter{WithGroup: true,
			Status: []int64{def.StatusWaitStop, def.StatusWaitDelete, def.StatusWaitRun},
			Types:  []int64{domain.TaskTypeQueue}},
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
						return QueueTaskStatusRunCheck(ctx, svcCtx, &wait, t)
					case def.StatusWaitDelete, def.StatusWaitStop:
						return QueueTaskStatusStopCheck(ctx, svcCtx, &wait, t)
					}
					//其他状态不需要处理
					return nil
				}()
				if err != nil {
					logx.WithContext(ctx).Errorf("QueueTaskCheck.one  err:%+v , task:%+v", err, t)
				}
			})
		}
		wait.Wait()
		return nil
	}()
	if err != nil {
		logx.WithContext(ctx).Errorf("QueueTaskCheck  err:%v", err)
	}
}

func QueueTaskClose(ctx context.Context, taskCode string) {
	QueueMutex.Lock()
	defer QueueMutex.Unlock()
	t := QueueMap[taskCode]
	if t != nil {
		for _, sub := range t.Subs {
			err := sub.Unsubscribe()
			if err != nil {
				logx.WithContext(ctx).Error(err)
			}
		}
		delete(QueueMap, taskCode)
	}
}

// 需要检查任务是否启动,如果没有启动需要启动
func QueueTaskStatusRunCheck(ctx context.Context, svcCtx *svc.ServiceContext, wait *sync.WaitGroup, task *relationDB.TimedTaskInfo) error {
	defer wait.Done()
	var err error
	taskCode := getQueueTaskCode(task)
	QueueTaskClose(ctx, taskCode)
	var val = QueueTrigger{
		Task: task,
	}
	for _, topic := range task.Topics {
		sub, err := svcCtx.Queue.SubscribeWithConsumer(topic, fmt.Sprintf("%s_%s", svcCtx.Config.Name, taskCode), func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			return QueueSendTask(ctx, svcCtx, natsMsg.Subject, string(msg), task)
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("QueueTaskStatusRunCheck.QueueSubscribe err:%v", err)
			continue
		}
		val.Subs = append(val.Subs, sub)
	}
	func() {
		QueueMutex.Lock()
		defer QueueMutex.Unlock()
		QueueMap[taskCode] = &val
	}()
	jobDB := relationDB.NewTaskInfoRepo(ctx)
	task.Status = def.StatusRunning
	err = jobDB.Update(ctx, task)
	return err
}

func QueueSendTask(ctx context.Context, svcCtx *svc.ServiceContext, topic string, payload string, po *relationDB.TimedTaskInfo) error {
	do := relationDB.ToTaskInfoDo(po)
	req := timedjob.TaskSendReq{
		GroupCode:   po.GroupCode,
		Code:        po.Code,
		ParamQueue:  nil,
		ParamSql:    nil,
		ParamScript: nil,
	}
	if do.Queue != nil {
		req.ParamQueue = &timedjob.TaskParamQueue{
			Topic:   do.Queue.Topic,
			Payload: do.Queue.Payload,
		}
	}
	if do.Sql != nil {
		req.ParamSql = &timedjob.TaskParamSql{
			Sql: do.Sql.Param.Sql,
		}
	}
	if do.Script != nil {
		req.ParamScript = &timedjob.TaskParamScript{
			Param: do.Script.Param.Param,
		}
		if req.ParamScript.Param == nil {
			req.ParamScript.Param = map[string]string{}
		}
		//补充参数
		req.ParamScript.Param["topic"] = topic
		req.ParamScript.Param["payload"] = payload
	}
	_, err := svcCtx.TimedM.TaskSend(ctx, &req)
	return err
}

// 如果处于运行状态需要停止
func QueueTaskStatusStopCheck(ctx context.Context, svcCtx *svc.ServiceContext, wait *sync.WaitGroup, task *relationDB.TimedTaskInfo) error {
	defer wait.Done()
	var err error
	taskCode := getQueueTaskCode(task)
	QueueTaskClose(ctx, taskCode)
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
func getQueueTaskCode(j *relationDB.TimedTaskInfo) string {
	return fmt.Sprintf("queue:%s:%s", j.GroupCode, j.Code)
}
