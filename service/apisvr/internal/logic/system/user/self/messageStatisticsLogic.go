package self

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageStatisticsLogic {
	return &MessageStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageStatisticsLogic) MessageStatistics() (resp *types.UserMessageStatisticsResp, err error) {
	ret, err := l.svcCtx.UserRpc.UserMessageStatistics(l.ctx, &sys.Empty{})
	return utils.Copy[types.UserMessageStatisticsResp](ret), nil
}
