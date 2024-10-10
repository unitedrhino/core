package rolemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleMenuMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	RmDB *relationDB.RoleMenuRepo
}

func NewRoleMenuMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleMenuMultiUpdateLogic {
	return &RoleMenuMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		RmDB:   relationDB.NewRoleMenuRepo(ctx),
	}
}

func (l *RoleMenuMultiUpdateLogic) RoleMenuMultiUpdate(in *sys.RoleMenuMultiUpdateReq) (*sys.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	err := l.RmDB.MultiUpdate(l.ctx, in.Id, in.AppCode, in.ModuleCode, in.MenuIDs)
	if err != nil {
		return nil, err
	}
	return &sys.Empty{}, nil
}
