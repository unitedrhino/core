package task

import (
	"context"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendLogic {
	return &SendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendLogic) Send(req *types.TimedTaskSendReq) (*types.TimedTaskWithTaskID, error) {
	l.Infof("req:%v", utils.Fmt(req))
	ret, err := l.svcCtx.TimedJob.TaskSend(l.ctx, ToSendDelayReqPb(req))
	if err != nil {
		return nil, err
	}
	return &types.TimedTaskWithTaskID{TaskID: ret.TaskID}, nil
}
