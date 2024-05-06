package module

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/tenant/app"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiCreateLogic {
	return &MultiCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiCreateLogic) MultiCreate(req *types.TenantAppCreateReq) error {
	_, err := l.svcCtx.TenantRpc.TenantAppModuleMultiCreate(l.ctx, &sys.TenantAppSaveReq{
		Code:    req.Code,
		AppCode: req.AppCode,
		Modules: app.ToTenantAppModulesPb(req.Modules),
	})
	return err
}
