package info

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/access/info"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TreeLogic {
	return &TreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TreeLogic) Tree(req *types.WithCode) (resp *types.TenantAccessInfoTreeResp, err error) {
	rst, err := l.svcCtx.TenantRpc.TenantAccessIndex(l.ctx, &sys.TenantAccessIndexReq{Code: req.Code})
	if err != nil {
		return nil, err
	}
	ais, err := l.svcCtx.AccessRpc.AccessInfoIndex(l.ctx, &sys.AccessInfoIndexReq{Codes: rst.AccessCodes})
	if err != nil {
		return nil, err
	}
	return &types.TenantAccessInfoTreeResp{
		List:  info.ToAccessModuleInfoTypes(ais.List),
		Total: ais.Total,
	}, nil
}
