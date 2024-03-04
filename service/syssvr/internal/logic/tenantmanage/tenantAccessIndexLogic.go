package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAccessIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAccessIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAccessIndexLogic {
	return &TenantAccessIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAccessIndexLogic) TenantAccessIndex(in *sys.TenantAccessIndexReq) (*sys.TenantAccessIndexResp, error) {
	if err := ctxs.IsRoot(l.ctx); err == nil && in.Code != "" {
		ctxs.GetUserCtx(l.ctx).AllTenant = true
		defer func() {
			ctxs.GetUserCtx(l.ctx).AllTenant = false
		}()
	}
	tas, err := relationDB.NewTenantAccessRepo(l.ctx).FindByFilter(l.ctx, relationDB.TenantAccessFilter{TenantCode: in.Code}, nil)
	if err != nil {
		return nil, err
	}
	var accessCodes []string
	for _, v := range tas {
		accessCodes = append(accessCodes, v.AccessCode)
	}

	return &sys.TenantAccessIndexResp{AccessCodes: accessCodes}, nil
}
