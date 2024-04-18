package startup

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/domain/slot"
	"gitee.com/i-Things/share/utils"
	ws "gitee.com/i-Things/share/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

func Init(svcCtx *svc.ServiceContext) {

	ws.RegisterSubscribeCheck(func(ctx context.Context, in *ws.SubscribeInfo) error {
		ret, err := svcCtx.Common.SlotInfoIndex(ctx, &sys.SlotInfoIndexReq{
			Code:    "userSubscribe",
			SubCode: in.Code,
		})
		if err != nil {
			return err
		}

		s := utils.CopySlice[slot.Info](ret.Slots)
		sl := slot.Infos(s)
		err = sl.Request(ctx, in, nil)
		return err
	})
	err := svcCtx.ServerMsg.Start()
	logx.Must(err)
}
