package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserMessageStatisticsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserMessageStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserMessageStatisticsLogic {
	return &UserMessageStatisticsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserMessageStatisticsLogic) UserMessageStatistics(in *sys.Empty) (*sys.UserMessageStatisticsResp, error) {
	count, err := relationDB.NewUserMessageRepo(l.ctx).CountNotRead(l.ctx, ctxs.GetUserCtx(l.ctx).UserID)
	if err != nil {
		return nil, err
	}
	var list []*sys.UserMessageStatistics
	for k, v := range count {
		list = append(list, &sys.UserMessageStatistics{
			Group: k,
			Count: v,
		})
	}
	return &sys.UserMessageStatisticsResp{List: list}, nil
}