package modulemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleMenuUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleMenuUpdateLogic {
	return &ModuleMenuUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleMenuUpdateLogic) ModuleMenuUpdate(in *sys.MenuInfo) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	old, err := relationDB.NewMenuInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	var tenantMenu relationDB.SysTenantAppMenu
	if in.Type != 0 && in.Type != old.Type {
		old.Type = in.Type
		tenantMenu.Type = in.Type
	}
	if in.Order != 0 && in.Order != old.Order {
		old.Order = in.Order
		tenantMenu.Order = in.Order
	}
	if in.Name != "" && in.Name != old.Name {
		old.Name = in.Name
		tenantMenu.Name = in.Name
	}
	if in.Path != "" && in.Path != old.Path {
		old.Path = in.Path
		tenantMenu.Path = in.Path
	}
	if in.Component != "" && in.Component != old.Component {
		old.Component = in.Component
		tenantMenu.Component = in.Component
	}
	if in.Icon != "" && in.Icon != old.Icon {
		old.Icon = in.Icon
		tenantMenu.Icon = in.Icon
	}
	if in.Redirect != "" && in.Redirect != old.Redirect {
		old.Redirect = in.Redirect
		tenantMenu.Redirect = in.Redirect
	}
	if in.Body != nil && in.Body.Value != old.Body {
		old.Body = in.Body.GetValue()
		tenantMenu.Body = in.Body.GetValue()
	}
	if in.IsCommon != 0 {
		old.IsCommon = in.IsCommon
		tenantMenu.IsCommon = in.IsCommon
	}
	if in.HideInMenu != 0 && in.HideInMenu != old.HideInMenu {
		old.HideInMenu = in.HideInMenu
		tenantMenu.HideInMenu = in.HideInMenu
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	conn := stores.GetTenantConn(l.ctx)
	err = conn.Transaction(func(tx *gorm.DB) error {
		err = relationDB.NewMenuInfoRepo(tx).Update(l.ctx, old)
		if err != nil {
			return err
		}
		err = relationDB.NewTenantAppMenuRepo(tx).UpdateByFilter(l.ctx, &tenantMenu, relationDB.TenantAppMenuFilter{
			TempLateID: in.Id,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return &sys.Empty{}, err
}
