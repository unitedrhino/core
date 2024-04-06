package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantNotifyMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantNotifyMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantNotifyMultiUpdateLogic {
	return &TenantNotifyMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 租户通知配置
func (l *TenantNotifyMultiUpdateLogic) TenantNotifyMultiUpdate(in *sys.TenantNotifyMultiUpdateReq) (*sys.Empty, error) {
	err := relationDB.NewTenantNotifyRepo(l.ctx).MultiInsert(l.ctx, utils.CopySlice[relationDB.SysTenantNotify](in.Notifies))

	return &sys.Empty{}, err
}
