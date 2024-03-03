package task

import (
	"context"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/client/timedmanage"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLogic {
	return &CancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelLogic) Cancel(req *types.TimedTaskWithTaskID) error {
	_, err := l.svcCtx.TimedJob.TaskCancel(l.ctx, &timedmanage.TaskWithTaskID{TaskID: req.TaskID})
	if err != nil {
		return err
	}
	return nil
}
