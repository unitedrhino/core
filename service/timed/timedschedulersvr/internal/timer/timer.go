package timer

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/internal/domain"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/timed/timedschedulersvr/internal/svc"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"time"
)

func Run(svcCtx *svc.ServiceContext) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	//dgsvr 订阅到了设备端数据，此时调用StartSpan方法，将订阅到的主题推送给jaeger
	//此时的ctx已经包含当前节点的span信息，会随着 handle(ctx).Publish 传递到下个节点
	ctx, span := ctxs.StartSpan(ctx, "timedSchedulersvr.taskRun", "")
	defer span.End()
	{ //先初始化数据库状态
		msg := "初始化数据库执行错误"
		jobDB := relationDB.NewTaskInfoRepo(ctx)
		//将运行中的任务修改为等待运行
		err := jobDB.UpdateByFilter(ctx, &relationDB.TimedTaskInfo{Status: def.StatusWaitRun},
			relationDB.TaskFilter{Status: []int64{def.StatusRunning}, Types: []int64{domain.TaskTypeTiming, domain.TaskTypeQueue}})
		errors.Must(err, msg)
		//将等待暂停的任务调整为已暂停
		err = jobDB.UpdateByFilter(ctx, &relationDB.TimedTaskInfo{Status: def.StatusStopped},
			relationDB.TaskFilter{Status: []int64{def.StatusWaitStop}, Types: []int64{domain.TaskTypeTiming, domain.TaskTypeQueue}})
		errors.Must(err, msg)
		//删除等待删除的任务
		err = jobDB.DeleteByFilter(ctx, relationDB.TaskFilter{Status: []int64{def.StatusWaitDelete}, Types: []int64{domain.TaskTypeTiming, domain.TaskTypeQueue}})
		errors.Must(err, msg)
	}
	utils.Go(ctx, func() {
		ctx := context.Background()
		TimingTaskCheck(svcCtx)
		utils.Go(ctx, func() {
			err := svcCtx.Scheduler.Run()
			errors.Must(err, "Scheduler.Run")
		})
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				TimingTaskCheck(svcCtx)
			}
		}
	})
	utils.Go(ctx, func() {
		QueueTaskCheck(svcCtx)
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				QueueTaskCheck(svcCtx)
			}
		}
	})

}
