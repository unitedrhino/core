package info

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.WithCode) (resp *types.TenantAccessInfo, err error) {
	rst, err := l.svcCtx.TenantRpc.TenantAccessIndex(l.ctx, &sys.TenantAccessIndexReq{Code: req.Code})
	if err != nil {
		return nil, err
	}
	return &types.TenantAccessInfo{Code: req.Code, AccessCodes: rst.AccessCodes}, nil
}
