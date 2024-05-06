package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantNotifyChannelReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantNotifyChannelReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantNotifyChannelReadLogic {
	return &TenantNotifyChannelReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantNotifyChannelReadLogic) TenantNotifyChannelRead(in *sys.WithID) (*sys.TenantNotifyChannel, error) {
	po, err := relationDB.NewTenantNotifyChannelRepo(l.ctx).FindOne(l.ctx, in.Id)
	return utils.Copy[sys.TenantNotifyChannel](po), err
}
