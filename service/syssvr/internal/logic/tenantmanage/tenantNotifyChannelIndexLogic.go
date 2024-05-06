package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantNotifyChannelIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantNotifyChannelIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantNotifyChannelIndexLogic {
	return &TenantNotifyChannelIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantNotifyChannelIndexLogic) TenantNotifyChannelIndex(in *sys.TenantNotifyChannelIndexReq) (*sys.TenantNotifyChannelIndexResp, error) {
	db := relationDB.NewTenantNotifyChannelRepo(l.ctx)
	f := relationDB.TenantNotifyChannelFilter{
		Name: in.Name,
		Type: in.Type,
	}
	l.ctx = ctxs.WithCommonTenant(l.ctx)
	total, err := db.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := db.FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	return &sys.TenantNotifyChannelIndexResp{Total: total, List: utils.CopySlice[sys.TenantNotifyChannel](pos)}, nil
}
