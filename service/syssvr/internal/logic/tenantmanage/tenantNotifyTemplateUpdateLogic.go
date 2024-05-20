package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
)

type TenantNotifyTemplateUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantNotifyTemplateUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantNotifyTemplateUpdateLogic {
	return &TenantNotifyTemplateUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 租户通知配置
func (l *TenantNotifyTemplateUpdateLogic) TenantNotifyTemplateUpdate(in *sys.TenantNotify) (*sys.Empty, error) {
	po := relationDB.SysTenantNotifyTemplate{
		NotifyCode: in.NotifyCode,
		Type:       in.Type,
		TemplateID: in.TemplateID,
	}
	err := relationDB.NewTenantNotifyTemplateRepo(l.ctx).Save(l.ctx, &po)

	return &sys.Empty{}, err
}
