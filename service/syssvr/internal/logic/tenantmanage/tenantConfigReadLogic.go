package tenantmanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantConfigReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantConfigReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantConfigReadLogic {
	return &TenantConfigReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantConfigReadLogic) TenantConfigRead(in *sys.WithCode) (*sys.TenantConfig, error) {
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	tenantCode := in.Code
	if tenantCode == "" {
		tenantCode = uc.TenantCode
	}
	if err := ctxs.IsRoot(l.ctx); err != nil {
		if uc.TenantCode != tenantCode {
			return nil, errors.Permissions
		}
	} else {
		uc.AllTenant = true
		defer func() {
			uc.AllTenant = false
		}()
	}
	po, err := relationDB.NewTenantConfigRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantConfigFilter{TenantCode: tenantCode})
	if err != nil {
		return nil, err
	}
	return ToTenantConfigPb(l.ctx, l.svcCtx, po), nil
}
