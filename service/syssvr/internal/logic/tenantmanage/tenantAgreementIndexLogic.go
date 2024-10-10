package tenantmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantAgreementIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantAgreementIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantAgreementIndexLogic {
	return &TenantAgreementIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantAgreementIndexLogic) TenantAgreementIndex(in *sys.TenantAgreementIndexReq) (*sys.TenantAgreementIndexResp, error) {
	f := relationDB.TenantAgreementFilter{}
	pos, err := relationDB.NewTenantAgreementRepo(l.ctx).FindByFilter(l.ctx, f, utils.Copy[stores.PageInfo](in))
	if err != nil {
		return nil, err
	}
	total, err := relationDB.NewTenantAgreementRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	return &sys.TenantAgreementIndexResp{
		List:  utils.CopySlice[sys.TenantAgreement](pos),
		Total: total,
	}, nil
}
