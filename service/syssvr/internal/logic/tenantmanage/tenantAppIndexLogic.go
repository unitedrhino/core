package tenantmanagelogic

import (
	"context"
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
	f := relationDB.TenantAppFilter{TenantCode: in.Code, AppCodes: in.AppCodes, AppID: in.AppID, Type: in.Type, SubType: in.SubType}
	list, err := relationDB.NewTenantAppRepo(l.ctx).FindByFilter(l.ctx, f, nil)
	if err != nil {
		return nil, err
	}
	var retList []*sys.TenantAppInfo
	for _, v := range list {
		val := ToTenantApp(l.ctx, l.svcCtx, v)
		retList = append(retList, val)
	}
	return &sys.TenantAppIndexResp{List: retList, Total: int64(len(list))}, nil
}
