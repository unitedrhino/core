package info

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

func (l *IndexLogic) Index(req *types.ModuleInfoIndexReq) (resp *types.ModuleInfoIndexResp, err error) {
	ret, err := l.svcCtx.ModuleRpc.ModuleInfoIndex(l.ctx, &sys.ModuleInfoIndexReq{
		Name:    req.Name,
		Page:    logic.ToSysPageRpc(req.Page),
		Code:    req.Code,
		Codes:   req.Codes,
		AppCode: req.AppCode,
		Type:    req.Type,
	})
	if err != nil {
		return nil, err
	}

	return &types.ModuleInfoIndexResp{
		List:     ToModuleInfosApi(ret.List),
		PageResp: logic.ToPageResp(req.Page, ret.Total),
	}, nil
}
