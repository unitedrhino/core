package tenantmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppModuleIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppModuleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppModuleIndexLogic {
	return &TenantAppModuleIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppModuleIndexLogic) TenantAppModuleIndex(in *sys.TenantModuleIndexReq) (*sys.TenantModuleIndexResp, error) {
	if err := ctxs.IsRoot(l.ctx); err == nil && in.Code != "" {
		ctxs.GetUserCtx(l.ctx).AllTenant = true
		defer func() {
			ctxs.GetUserCtx(l.ctx).AllTenant = false
		}()
	}
	ret, err := relationDB.NewTenantAppModuleRepo(l.ctx).FindByFilter(l.ctx,
		relationDB.TenantAppModuleFilter{TenantCode: in.Code, AppCodes: []string{in.AppCode}, ModuleCodes: in.ModuleCodes, WithModule: true}, nil)
	if err != nil {
		return nil, err
	}
	var modules []*relationDB.SysModuleInfo
	for _, v := range ret {
		if v.Module == nil {
			continue
		}
		modules = append(modules, v.Module)
	}
	return &sys.TenantModuleIndexResp{List: logic.ToModuleInfosPb(modules)}, nil
}
