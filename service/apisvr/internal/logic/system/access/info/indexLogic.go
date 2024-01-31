package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
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

func (l *IndexLogic) Index(req *types.AccessIndexReq) (resp *types.AccessIndexResp, err error) {
	rst, err := l.svcCtx.AccessRpc.AccessInfoIndex(l.ctx, &sys.AccessInfoIndexReq{
		Page:       logic.ToSysPageRpc(req.Page),
		Group:      req.Group,
		Code:       req.Code,
		Name:       req.Name,
		IsNeedAuth: req.IsNeedAuth,
		WithApis:   req.WithApis,
	})
	if err != nil {
		return nil, err
	}
	return &types.AccessIndexResp{
		List:  ToAccessInfosTypes(rst.List),
		Total: rst.Total,
	}, nil
}
