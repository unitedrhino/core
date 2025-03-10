package workOrder

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
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
	ret, err := l.svcCtx.Ops.OpsWorkOrderIndex(l.ctx, utils.Copy[sys.OpsWorkOrderIndexReq](req))
	if err != nil {
		return nil, err
	}
	return &types.OpsWorkOrderIndexResp{
		PageResp: logic.ToPageResp(req.Page, ret.Total),
		List:     utils.CopySlice[types.OpsWorkOrder](ret.List),
	}, nil
}
