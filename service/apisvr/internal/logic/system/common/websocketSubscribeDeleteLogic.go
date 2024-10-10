package common

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebsocketSubscribeDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebsocketSubscribeDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsocketSubscribeDeleteLogic {
	return &WebsocketSubscribeDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebsocketSubscribeDeleteLogic) WebsocketSubscribeDelete(req *types.WebsocketSubscribeInfo) error {
	//err := l.svcCtx.UserSubscribe.Del(l.ctx, ctxs.GetUserCtx(l.ctx).UserID, &websockets.SubscribeInfo{
	//	Code:   req.Code,
	//	Params: req.Params,
	//})
	return nil
}
