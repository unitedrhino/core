package usermanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/cache"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRoleMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleMultiUpdateLogic {
	return &UserRoleMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRoleMultiUpdateLogic) UserRoleMultiUpdate(in *sys.UserRoleMultiUpdateReq) (*sys.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	err := relationDB.NewUserRoleRepo(l.ctx).MultiUpdate(l.ctx, in.UserID, in.RoleIDs)
	if err == nil {
		l.svcCtx.UsersCache.SetData(l.ctx, in.UserID, nil)
		cache.ClearProjectAuth(in.UserID)
	}
	return &sys.Empty{}, err
}
