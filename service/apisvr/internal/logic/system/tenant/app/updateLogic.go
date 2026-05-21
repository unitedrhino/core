package app

import (
	"context"

	tenantoauth "gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/tenant"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.TenantAppInfo) error {
	_, err := l.svcCtx.TenantRpc.TenantAppUpdate(l.ctx, tenantoauth.ToSysTenantAppInfo(req))
	return err
}
