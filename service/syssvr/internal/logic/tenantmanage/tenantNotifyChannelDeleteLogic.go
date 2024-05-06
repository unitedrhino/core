package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantNotifyChannelDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantNotifyChannelDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantNotifyChannelDeleteLogic {
	return &TenantNotifyChannelDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantNotifyChannelDeleteLogic) TenantNotifyChannelDelete(in *sys.WithID) (*sys.Empty, error) {
	err := relationDB.NewTenantNotifyChannelRepo(l.ctx).Delete(l.ctx, in.Id)
	return &sys.Empty{}, err
}
