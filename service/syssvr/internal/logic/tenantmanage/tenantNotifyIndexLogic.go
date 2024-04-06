package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantNotifyIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantNotifyIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantNotifyIndexLogic {
	return &TenantNotifyIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantNotifyIndexLogic) TenantNotifyIndex(in *sys.TenantNotifyIndexReq) (*sys.TenantNotifyIndexResp, error) {
	db := relationDB.NewTenantNotifyRepo(l.ctx)
	pos, err := db.FindByFilter(l.ctx, relationDB.TenantNotifyConfigFilter{
		NotifyCode: in.NotifyCode,
		Type:       in.Type,
	}, nil)
	return &sys.TenantNotifyIndexResp{List: utils.CopySlice[sys.TenantNotify](pos)}, err
}
