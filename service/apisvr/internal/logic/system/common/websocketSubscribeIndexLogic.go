package common

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

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
	list, err := l.svcCtx.UserSubscribe.Index(l.ctx, ctxs.GetUserCtx(l.ctx).UserID)
	if err != nil {
		return nil, err
	}
	return &types.WebsocketSubscribeIndexResp{List: utils.CopySlice[types.WebsocketSubscribeInfo](list)}, nil
}
