package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppModuleMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppModuleMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppModuleMultiCreateLogic {
	return &TenantAppModuleMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppModuleMultiCreateLogic) TenantAppModuleMultiCreate(in *sys.TenantAppInfo) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	uc := ctxs.GetUserCtx(l.ctx)
	uc.AllTenant = true
	defer func() { uc.AllTenant = false }()
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		for _, module := range in.Modules {
			err := ModuleCreate(l.ctx, tx, in.Code, in.AppCode, module.Code, module.MenuIDs)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return &sys.Empty{}, err
}
