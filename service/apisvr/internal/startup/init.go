package startup

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/sysExport"
	"gitee.com/unitedrhino/core/share/domain/slot"
	ws "gitee.com/unitedrhino/share/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func Init(svcCtx *svc.ServiceContext) {

	ws.RegisterSubscribeCheck(func(ctx context.Context, in *ws.SubscribeInfo) error {
		ctx, _ = context.WithTimeout(ctx, 2*time.Second)
		sl, err := svcCtx.Slot.GetData(ctx, sysExport.GenSlotCacheKey(slot.CodeUserSubscribe, in.Code))
		if err != nil {
			return err
		}
		err = sl.Request(ctx, in, nil)
		return err
	})
	err := svcCtx.ServerMsg.Start()
	logx.Must(err)
}
