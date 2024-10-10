package tenantmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppReadLogic {
	return &TenantAppReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppReadLogic) TenantAppRead(in *sys.TenantAppWithIDOrCode) (*sys.TenantAppInfo, error) {
	if err := ctxs.IsRoot(l.ctx); err == nil && in.Code != "" {
		ctxs.GetUserCtx(l.ctx).AllTenant = true
		defer func() {
			ctxs.GetUserCtx(l.ctx).AllTenant = false
		}()
	}
	if in.AppCode == "" {
		in.AppCode = ctxs.GetUserCtxNoNil(l.ctx).AppCode
	}
	po, err := relationDB.NewTenantAppRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantAppFilter{TenantCode: in.Code, AppCodes: []string{in.AppCode}})
	if err != nil {
		return nil, err
	}
	return ToTenantApp(l.ctx, l.svcCtx, po), nil
}
