package tenantmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAgreementCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAgreementCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAgreementCreateLogic {
	return &TenantAgreementCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAgreementCreateLogic) TenantAgreementCreate(in *sys.TenantAgreement) (*sys.WithID, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	in.Id = 0
	po := utils.Copy[relationDB.SysTenantAgreement](in)
	err := relationDB.NewTenantAgreementRepo(l.ctx).Insert(l.ctx, po)
	return &sys.WithID{Id: po.ID}, err
}
