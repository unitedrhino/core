package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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

func (l *TenantAppMenuUpdateLogic) TenantAppMenuUpdate(in *sys.TenantAppMenu) (*sys.Response, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	old, err := relationDB.NewTenantAppMenuRepo(l.ctx).FindOne(l.ctx, in.Info.Id)
	old.Type = in.Info.Type
	old.Order = in.Info.Order
	old.Name = in.Info.Name
	old.Path = in.Info.Path
	old.Component = in.Info.Component
	old.Icon = in.Info.Icon
	old.Redirect = in.Info.Redirect
	old.Body = in.Info.Body.Value
	old.HideInMenu = in.Info.HideInMenu
	err = relationDB.NewTenantAppMenuRepo(l.ctx).Update(l.ctx, old)
	return &sys.Response{}, err
}
