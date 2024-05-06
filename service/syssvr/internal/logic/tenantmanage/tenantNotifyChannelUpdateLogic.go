package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantNotifyChannelUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantNotifyChannelUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantNotifyChannelUpdateLogic {
	return &TenantNotifyChannelUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantNotifyChannelUpdateLogic) TenantNotifyChannelUpdate(in *sys.TenantNotifyChannel) (*sys.Empty, error) {
	old, err := relationDB.NewTenantNotifyChannelRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		old.Name = in.Name
	}
	if in.Email != nil {
		old.Email = utils.Copy[relationDB.SysTenantEmail](in.Email)
	}
	if in.Webhook != "" {
		old.WebHook = in.Webhook
	}
	old.Desc = in.Desc
	err = relationDB.NewTenantNotifyChannelRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
