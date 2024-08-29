package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAccessMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAccessMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAccessMultiDeleteLogic {
	return &TenantAccessMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAccessMultiDeleteLogic) TenantAccessMultiDelete(in *sys.TenantAccessMultiSaveReq) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	uc := ctxs.GetUserCtx(l.ctx)
	uc.AllTenant = true
	defer func() { uc.AllTenant = false }()
	err := relationDB.NewTenantAccessRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.TenantAccessFilter{
		TenantCode:  in.Code,
		AccessCodes: in.AccessCodes,
	})
	return &sys.Empty{}, err
}
