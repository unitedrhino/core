package tenantmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/ctxs"

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
			TenantCode: dataType.TenantCode(in.Code),
			AccessCode: v,
		})
	}
	err := relationDB.NewTenantAccessRepo(l.ctx).MultiInsert(l.ctx, datas)
	return &sys.Empty{}, err
}
