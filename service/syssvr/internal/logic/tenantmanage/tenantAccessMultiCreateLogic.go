package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/stores"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAccessMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAccessMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAccessMultiCreateLogic {
	return &TenantAccessMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAccessMultiCreateLogic) TenantAccessMultiCreate(in *sys.TenantAccessMultiSaveReq) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	uc := ctxs.GetUserCtx(l.ctx)
	uc.AllTenant = true
	defer func() { uc.AllTenant = false }()
	var datas []*relationDB.SysTenantAccess
	for _, v := range in.AccessCodes {
		datas = append(datas, &relationDB.SysTenantAccess{
			TenantCode: stores.TenantCode(in.Code),
			AccessCode: v,
		})
	}
	err := relationDB.NewTenantAccessRepo(l.ctx).MultiInsert(l.ctx, datas)
	return &sys.Empty{}, err
}
