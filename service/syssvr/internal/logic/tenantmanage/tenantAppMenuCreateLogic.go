package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/shared/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppMenuCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppMenuCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppMenuCreateLogic {
	return &TenantAppMenuCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppMenuCreateLogic) TenantAppMenuCreate(in *sys.TenantAppMenu) (*sys.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	if err := CheckModule(l.ctx, in.Code, in.AppCode, in.Info.ModuleCode); err != nil {
		return nil, err
	}
	po := logic.ToTenantAppMenuPo(in)
	relationDB.NewTenantAppMenuRepo(l.ctx).Insert(l.ctx, po)
	return &sys.WithID{Id: po.ID}, nil
}
