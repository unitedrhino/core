package tenantmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppModuleDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppModuleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppModuleDeleteLogic {
	return &TenantAppModuleDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppModuleDeleteLogic) TenantAppModuleDelete(in *sys.TenantModuleWithIDOrCode) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	f := relationDB.TenantAppModuleFilter{ID: in.Id, TenantCode: in.Code}
	if in.AppCode != "" {
		f.AppCodes = []string{in.AppCode}
	}
	if in.ModuleCode != "" {
		f.ModuleCodes = []string{in.ModuleCode}
	}
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewTenantAppModuleRepo(tx).DeleteByFilter(l.ctx, f)
		if err != nil {
			return err
		}
		err = relationDB.NewTenantAppMenuRepo(tx).DeleteByFilter(l.ctx, relationDB.TenantAppMenuFilter{
			ModuleCode: in.ModuleCode,
			TenantCode: in.Code,
			AppCode:    in.AppCode,
		})
		if err != nil {
			return err
		}
		return nil
	})

	return &sys.Empty{}, err
}
