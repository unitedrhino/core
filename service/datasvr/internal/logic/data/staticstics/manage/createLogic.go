package manage

import (
	"context"
	"gitee.com/i-Things/core/service/datasvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/datasvr/internal/svc"
	"gitee.com/i-Things/core/service/datasvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.DataStatisticsManage) (resp *types.WithID, err error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	po := utils.Copy[relationDB.DataStatisticsInfo](req)
	po.ID = 0
	err = relationDB.NewStatisticsInfoRepo(l.ctx).Insert(l.ctx, po)
	return &types.WithID{ID: po.ID}, err
}
