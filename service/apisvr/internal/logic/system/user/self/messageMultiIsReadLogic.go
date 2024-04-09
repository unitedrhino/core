package self

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageMultiIsReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageMultiIsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageMultiIsReadLogic {
	return &MessageMultiIsReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageMultiIsReadLogic) MessageMultiIsRead(req *types.IDList) error {
	_, err := l.svcCtx.UserRpc.UserMessageMultiIsRead(l.ctx, utils.Copy[sys.IDList](req))
	return err
}
