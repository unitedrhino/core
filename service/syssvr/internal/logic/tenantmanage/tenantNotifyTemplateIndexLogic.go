package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantNotifyTemplateIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantNotifyTemplateIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantNotifyTemplateIndexLogic {
	return &TenantNotifyTemplateIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantNotifyTemplateIndexLogic) TenantNotifyTemplateIndex(in *sys.TenantNotifyTemplateIndexReq) (*sys.TenantNotifyTemplateIndexResp, error) {
	db := relationDB.NewTenantNotifyTemplateRepo(l.ctx)
	pos, err := db.FindByFilter(l.ctx, relationDB.TenantNotifyTemplateFilter{
		NotifyCode: in.NotifyCode,
		Type:       in.Type,
	}, nil)
	return &sys.TenantNotifyTemplateIndexResp{List: utils.CopySlice[sys.TenantNotify](pos)}, err
}
