package usermanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

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
	var userID = ctxs.GetUserCtxNoNil(l.ctx).UserID
	err := relationDB.NewUserMessageRepo(l.ctx).MultiIsRead(l.ctx, userID, in.Ids)
	if err != nil {
		return nil, err
	}
	return &sys.Empty{}, err
}
