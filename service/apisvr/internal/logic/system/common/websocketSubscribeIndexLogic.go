package common

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebsocketSubscribeIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebsocketSubscribeIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsocketSubscribeIndexLogic {
	return &WebsocketSubscribeIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebsocketSubscribeIndexLogic) WebsocketSubscribeIndex() (resp *types.WebsocketSubscribeIndexResp, err error) {
	return
}
