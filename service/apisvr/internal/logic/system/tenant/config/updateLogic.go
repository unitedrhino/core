package config

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"

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

func (l *UpdateLogic) Update(req *types.TenantConfig) error {
	tenantCode := req.TenantCode
	if tenantCode == "" {
		tenantCode = ctxs.GetUserCtxNoNil(l.ctx).TenantCode
	}
	_, err := l.svcCtx.TenantRpc.TenantConfigUpdate(l.ctx, utils.Copy[sys.TenantConfig](req))
	if err != nil {
		return err
	}
	if !hasAppLoginPayload(req) {
		return nil
	}
	appCode := req.AppCode
	if appCode == "" {
		appCode = ctxs.GetUserCtxNoNil(l.ctx).AppCode
	}
	_, err = l.svcCtx.TenantRpc.TenantAppUpdate(l.ctx, toTenantAppInfo(req, tenantCode, appCode))
	return err
}
