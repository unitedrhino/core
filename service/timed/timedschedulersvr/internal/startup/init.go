package startup

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/timedschedulersvr/internal/svc"
	"gitee.com/unitedrhino/core/service/timed/timedschedulersvr/internal/timer"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func Init(svcCtx *svc.ServiceContext) error {
	utils.Go(context.Background(), func() {
		utils.SingletonRun(context.Background(), svcCtx.Store, "svr:timedschedulersvr", func(ctx2 context.Context) {
			svcCtx.SchedulerRun = true
			logx.Info("timedschedulersvr 开始运行")
			timer.Run(svcCtx)
		})
	})
	return nil
}
