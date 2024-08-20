package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAppUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAppUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAppUpdateLogic {
	return &TenantAppUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAppUpdateLogic) TenantAppUpdate(in *sys.TenantAppInfo) (*sys.Empty, error) {
	old, err := relationDB.NewTenantAppRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantAppFilter{TenantCode: in.Code, AppCodes: []string{in.AppCode}})
	if err != nil {
		return nil, err
	}
	if in.WxMini != nil {
		old.WxMini = utils.Copy[relationDB.SysTenantThird](in.WxMini)
	}
	if in.DingMini != nil {
		old.DingMini = utils.Copy[relationDB.SysTenantThird](in.DingMini)
	}
	if in.WxOpen != nil {
		old.WxOpen = utils.Copy[relationDB.SysTenantThird](in.WxOpen)
	}
	if in.LoginTypes != nil {
		old.LoginTypes = in.LoginTypes
	}
	err = relationDB.NewTenantAppRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
