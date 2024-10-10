package api

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

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

func (l *IndexLogic) Index(req *types.ApiInfoIndexReq) (resp *types.ApiInfoIndexResp, err error) {
	info, err := l.svcCtx.ModuleRpc.ModuleApiIndex(l.ctx, &sys.ApiInfoIndexReq{
		Page:       logic.ToSysPageRpc(req.Page),
		Route:      req.Route,
		Method:     req.Method,
		Group:      req.Group,
		Name:       req.Name,
		IsNeedAuth: req.IsNeedAuth,
	})
	if err != nil {
		return nil, err
	}

	var total int64
	total = info.Total
	var apiInfo []*types.ApiInfo
	apiInfo = make([]*types.ApiInfo, 0, len(apiInfo))
	for _, i := range info.List {
		apiInfo = append(apiInfo, ToApiInfoTypes(i))
	}
	return &types.ApiInfoIndexResp{List: apiInfo, Total: total}, nil
}
