package tenantmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppMenuUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppMenuUpdateLogic {
	return &TenantAppMenuUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppMenuUpdateLogic) TenantAppMenuUpdate(in *sys.TenantAppMenu) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	old, err := relationDB.NewTenantAppMenuRepo(l.ctx).FindOne(l.ctx, in.Info.Id)
	if err != nil {
		return nil, err
	}
	if in.Info.Order != 0 {
		old.Order = in.Info.Order
	}
	if in.Info.Name != "" {
		old.Name = in.Info.Name
	}
	if in.Info.HideInMenu != 0 {
		old.HideInMenu = in.Info.HideInMenu
	}
	if in.Info.IsCommon != 0 {
		old.IsCommon = in.Info.IsCommon
	}
	err = relationDB.NewTenantAppMenuRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
