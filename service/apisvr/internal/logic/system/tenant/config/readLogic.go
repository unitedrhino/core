package config

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	tenantoauth "gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/tenant"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.WithCode) (resp *types.TenantConfig, err error) {
	tenantCode := req.Code
	if tenantCode == "" {
		tenantCode = ctxs.GetUserCtxNoNil(l.ctx).TenantCode
	}
	ret, err := l.svcCtx.TenantRpc.TenantConfigRead(l.ctx, &sys.WithCode{Code: tenantCode})
	if err != nil {
		return nil, err
	}
	resp = utils.Copy[types.TenantConfig](ret)
	appCode := ctxs.GetUserCtxNoNil(l.ctx).AppCode
	app, err := l.loadTenantApp(tenantCode, appCode)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return resp, nil
		}
		return nil, err
	}
	mergeAppLoginIntoConfig(resp, app)
	tenantoauth.FillTenantConfigOut(resp)
	return resp, nil
}

// loadTenantApp 读取租户应用登录配置，优先指定 appCode，否则取该租户第一个应用
func (l *ReadLogic) loadTenantApp(tenantCode, appCode string) (*sys.TenantAppInfo, error) {
	if appCode != "" {
		return l.svcCtx.TenantRpc.TenantAppRead(l.ctx, &sys.TenantAppWithIDOrCode{
			Code:    tenantCode,
			AppCode: appCode,
		})
	}
	idx, err := l.svcCtx.TenantRpc.TenantAppIndex(l.ctx, &sys.TenantAppIndexReq{Code: tenantCode})
	if err != nil {
		return nil, err
	}
	if len(idx.List) == 0 {
		return nil, errors.NotFind
	}
	return idx.List[0], nil
}
