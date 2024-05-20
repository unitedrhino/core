package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantNotifyTemplateMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantNotifyTemplateMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantNotifyTemplateMultiUpdateLogic {
	return &TenantNotifyTemplateMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 租户通知配置
func (l *TenantNotifyTemplateMultiUpdateLogic) TenantNotifyTemplateMultiUpdate(in *sys.TenantNotifyTemplateMultiUpdateReq) (*sys.Empty, error) {
	return &sys.Empty{}, nil //暂停使用
	err := relationDB.NewTenantNotifyTemplateRepo(l.ctx).MultiUpdate(l.ctx, utils.CopySlice[relationDB.SysTenantNotifyTemplate](in.Notifies))
	return &sys.Empty{}, err
}
