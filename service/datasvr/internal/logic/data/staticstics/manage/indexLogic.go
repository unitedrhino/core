package manage

import (
	"context"
	"gitee.com/i-Things/core/service/datasvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/datasvr/internal/svc"
	"gitee.com/i-Things/core/service/datasvr/internal/types"

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

func (l *IndexLogic) Index(req *types.DataStatisticsManageIndexReq) (resp *types.DataStatisticsManageIndexResp, err error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	f := relationDB.StatisticsInfoFilter{}
	total, err := relationDB.NewStatisticsInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := relationDB.NewStatisticsInfoRepo(l.ctx).FindByFilter(l.ctx, f, utils.Copy[stores.PageInfo](req.Page))
	if err != nil {
		return nil, err
	}
	return &types.DataStatisticsManageIndexResp{
		List:  utils.CopySlice[types.DataStatisticsManage](pos),
		Total: total,
	}, err
}
