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

type TenantAppModuleCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppModuleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppModuleCreateLogic {
	return &TenantAppModuleCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppModuleCreateLogic) TenantAppModuleCreate(in *sys.TenantModuleCreateReq) (*sys.Response, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	conn := stores.GetTenantConn(l.ctx)
	err := conn.Transaction(func(tx *gorm.DB) error {
		err := ModuleCreate(l.ctx, tx, in.Code, in.AppCode, in.ModuleCode, in.MenuIDs)
		return err
	})
	return &sys.Response{}, err
}
