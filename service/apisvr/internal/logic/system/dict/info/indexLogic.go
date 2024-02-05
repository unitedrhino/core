package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/dict"
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

func (l *IndexLogic) Index(req *types.DictInfoIndexReq) (resp *types.DictInfoIndexResp, err error) {
	ret, err := l.svcCtx.DictM.DictInfoIndex(l.ctx, &sys.DictInfoIndexReq{
		Page:        logic.ToSysPageRpc(req.Page),
		Name:        req.Name,
		Type:        req.Type,
		Status:      req.Status,
		WithDetails: req.WithDetails,
	})
	if err != nil {
		return nil, err
	}

	return &types.DictInfoIndexResp{
		Total: ret.Total,
		List:  dict.ToInfosTypes(ret.List),
	}, nil
}
