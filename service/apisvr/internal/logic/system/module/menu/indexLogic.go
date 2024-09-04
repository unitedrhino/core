package menu

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system"
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

func (l *IndexLogic) Index(req *types.MenuInfoIndexReq) (resp *types.MenuInfoIndexResp, err error) {
	info, err := l.svcCtx.ModuleRpc.ModuleMenuIndex(l.ctx, &sys.MenuInfoIndexReq{
		ModuleCode: req.ModuleCode,
		IsRetTree:  req.IsRetTree,
		IsCommon:   req.IsCommon,
	})
	if err != nil {
		return nil, err
	}

	return &types.MenuInfoIndexResp{List: system.ToMenuInfosApi(info.List)}, nil
}
