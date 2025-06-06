package tenantmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantInfoIndexLogic {
	return &TenantInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取区域信息列表
func (l *TenantInfoIndexLogic) TenantInfoIndex(in *sys.TenantInfoIndexReq) (*sys.TenantInfoIndexResp, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		ret, err := relationDB.NewTenantInfoRepo(l.ctx).FindOneByFilter(l.ctx,
			relationDB.TenantInfoFilter{Code: ctxs.GetUserCtx(l.ctx).TenantCode})
		if err != nil {
			return nil, err
		}
		return &sys.TenantInfoIndexResp{List: []*sys.TenantInfo{logic.ToTenantInfoRpc(l.ctx, l.svcCtx, ret)}, Total: 1}, nil
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	f := relationDB.TenantInfoFilter{
		Code:    in.Code,
		Name:    in.Name,
		AppCode: in.AppCode,
	}
	list, err := relationDB.NewTenantInfoRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	total, err := relationDB.NewTenantInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	return &sys.TenantInfoIndexResp{List: logic.ToTenantInfosRpc(l.ctx, l.svcCtx, list), Total: total}, nil
}
