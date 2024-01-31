package tenantmanagelogic

import (
	"context"
	appmanagelogic "gitee.com/i-Things/core/service/syssvr/internal/logic/appmanage"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppIndexLogic {
	return &TenantAppIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppIndexLogic) TenantAppIndex(in *sys.TenantAppIndexReq) (*sys.TenantAppIndexResp, error) {
	if err := ctxs.IsRoot(l.ctx); err == nil && in.Code != "" {
		ctxs.GetUserCtx(l.ctx).AllTenant = true
		defer func() {
			ctxs.GetUserCtx(l.ctx).AllTenant = false
		}()
	}
	f := relationDB.TenantAppFilter{TenantCode: in.Code, AppCodes: in.AppCodes}
	list, err := relationDB.NewTenantAppRepo(l.ctx).FindByFilter(l.ctx, f, nil)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return &sys.TenantAppIndexResp{List: []*sys.AppInfo{}, Total: 0}, nil
	}
	appCodes := make([]string, 0)
	codeIDMap := make(map[string]int64)
	for _, v := range list {
		appCodes = append(appCodes, v.AppCode)
		codeIDMap[v.AppCode] = v.ID
	}
	apps, err := relationDB.NewAppInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.AppInfoFilter{Codes: appCodes}, nil)
	if err != nil {
		return nil, err
	}
	for _, v := range apps {
		v.ID = codeIDMap[v.Code] //修正为关联的id
	}
	return &sys.TenantAppIndexResp{List: appmanagelogic.ToAppInfosPb(apps), Total: int64(len(list))}, nil
}
