package timer

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/core/service/timed/internal/domain"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/logic"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/logic/processTask"
	timedmanagelogic "gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/logic/timedmanage"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/trace"
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
		var span trace.Span
		ctx, span = ctxs.StartSpan(ctx, "timedJob.Process", "")
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
