package tenantmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantOpenWebHookLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantOpenWebHookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantOpenWebHookLogic {
	return &TenantOpenWebHookLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantOpenWebHookLogic) TenantOpenWebHook(in *sys.WithCode) (*sys.TenantOpenWebHook, error) {
	po, err := relationDB.NewTenantOpenWebhookRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantOpenWebhookFilter{Code: in.Code})
	return utils.Copy[sys.TenantOpenWebHook](po), err
}
