package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserMessageMultiIsReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserMessageMultiIsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserMessageMultiIsReadLogic {
	return &UserMessageMultiIsReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserMessageMultiIsReadLogic) UserMessageMultiIsRead(in *sys.IDList) (*sys.Empty, error) {
	err := relationDB.NewUserMessageRepo(l.ctx).MultiIsRead(l.ctx, ctxs.GetUserCtxNoNil(l.ctx).UserID, in.Ids)
	return &sys.Empty{}, err
}
