package workOrder

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"
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

func (l *IndexLogic) Index(req *types.OpsWorkOrderIndexReq) (resp *types.OpsWorkOrderIndexResp, err error) {
	ret, err := l.svcCtx.Ops.OpsWorkOrderIndex(l.ctx, &sys.OpsWorkOrderIndexReq{
		Page:   logic.ToSysPageRpc(req.Page),
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &types.OpsWorkOrderIndexResp{
		Total: ret.Total,
		List:  utils.CopySlice[types.OpsWorkOrder](ret.List),
	}, nil
}
