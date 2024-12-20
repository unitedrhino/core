package startup

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/event"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/repo/event/subscribe"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/timer"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func Init(svcCtx *svc.ServiceContext) error {
	Subscribe(svcCtx)
	return InitTimer(svcCtx)
}

func Subscribe(svcCtx *svc.ServiceContext) {
	subAppCli, err := subscribe.NewSubServer(svcCtx.Config.Event, svcCtx.NodeID)
	logx.Must(err)
	err = subAppCli.Subscribe(func(ctx context.Context) subscribe.ServerEvent {
		return event.NewEventServer(ctx, svcCtx)
	})
}

func InitTimer(svcCtx *svc.ServiceContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	//dgsvr 订阅到了设备端数据，此时调用StartSpan方法，将订阅到的主题推送给jaeger
	//此时的ctx已经包含当前节点的span信息，会随着 handle(ctx).Publish 传递到下个节点
	ctx, span := ctxs.StartSpan(ctx, "InitTimer", "")
	defer span.End()
	as := clients.NewAsynqServer(svcCtx.Config.CacheRedis)
	utils.Go(ctx, func() {
		as.Run(timer.Timed{SvcCtx: svcCtx})
	})
	return nil
}
