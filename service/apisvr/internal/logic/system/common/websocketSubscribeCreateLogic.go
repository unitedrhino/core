package common

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/sysExport"
	"gitee.com/i-Things/share/domain/slot"
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

	sl, err := l.svcCtx.Slot.GetData(l.ctx, sysExport.GenSlotCacheKey(slot.CodeUserSubscribe, req.Code))
	if err != nil {
		return err
	}
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
