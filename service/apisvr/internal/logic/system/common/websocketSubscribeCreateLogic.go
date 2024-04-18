package common

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/domain/slot"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type WebsocketSubscribeCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebsocketSubscribeCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsocketSubscribeCreateLogic {
	return &WebsocketSubscribeCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebsocketSubscribeCreateLogic) WebsocketSubscribeCreate(req *types.WebsocketSubscribeInfo) error {

	ret, err := l.svcCtx.Common.SlotInfoIndex(l.ctx, &sys.SlotInfoIndexReq{
		Code:    "userSubscribe",
		SubCode: req.Code,
	})
	if err != nil {
		return err
	}

	s := utils.CopySlice[slot.Info](ret.Slots)
	sl := slot.Infos(s)
	err = sl.Request(l.ctx, req, nil)
	if err != nil {
		return err
	}
	//err = l.svcCtx.UserSubscribe.Add(l.ctx, ctxs.GetUserCtx(l.ctx).UserID, &websockets.SubscribeInfo{
	//	Code:   req.Code,
	//	Params: req.Params,
	//})
	return err
}
