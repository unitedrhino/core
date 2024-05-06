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

func (l *TenantAppUpdateLogic) TenantAppUpdate(in *sys.TenantAppSaveReq) (*sys.Empty, error) {
	old, err := relationDB.NewTenantAppRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantAppFilter{TenantCode: in.Code, AppCodes: []string{in.AppCode}})
	if err != nil {
		return nil, err
	}
	if in.MiniWx != nil {
		old.MiniWx = utils.Copy[relationDB.SysTenantThird](in.MiniWx)
	}
	if in.MiniDing != nil {
		old.MiniDing = utils.Copy[relationDB.SysTenantThird](in.MiniDing)
	}
	err = relationDB.NewTenantAppRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
