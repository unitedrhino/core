package app

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.TenantAppCreateReq) error {
	_, err := l.svcCtx.TenantRpc.TenantAppCreate(l.ctx, &sys.TenantAppSaveReq{
		Code:    req.Code,
		AppCode: req.AppCode,
		Modules: ToTenantAppModulesPb(req.Modules),
	})
	return err
}
