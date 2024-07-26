package manage

import (
	"context"
	"gitee.com/i-Things/core/service/datasvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/datasvr/internal/svc"
	"gitee.com/i-Things/core/service/datasvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.WithID) error {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return err
	}
	err := relationDB.NewStatisticsInfoRepo(l.ctx).Delete(l.ctx, req.ID)

	return err
}
