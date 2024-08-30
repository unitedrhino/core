package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"

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
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	po, err := relationDB.NewTenantConfigRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantConfigFilter{TenantCode: in.Code})
	if err != nil {
		return nil, err
	}
	return ToTenantConfigPb(l.ctx, l.svcCtx, po), nil
}
