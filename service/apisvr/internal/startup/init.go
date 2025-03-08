package startup

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/sysExport"
	"gitee.com/unitedrhino/core/share/domain/slot"
	ws "gitee.com/unitedrhino/core/share/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type retStu struct {
	List []map[string]any `json:"list"`
}

func Init(svcCtx *svc.ServiceContext) {

	ws.RegisterSubscribeCheck2(func(ctx context.Context, in *ws.SubscribeInfo) ([]map[string]any, error) {
		ctx, _ = context.WithTimeout(ctx, 2*time.Second)
		sl, err := svcCtx.Slot.GetData(ctx, sysExport.GenSlotCacheKey(slot.CodeUserSubscribe, in.Code))
		if err != nil {
			return nil, err
		}
		var ret retStu
		err = sl.Request(ctx, in, &ret)
		return ret.List, err
	})
	err := svcCtx.ServerMsg.Start()
	logx.Must(err)
}
