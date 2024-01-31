package timedmanagelogic

import (
	"context"

	"gitee.com/i-Things/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskLogIndexLogic {
	return &TaskLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskLogIndexLogic) TaskLogIndex(in *timedjob.TaskLogIndexReq) (*timedjob.TaskLogIndexResp, error) {
	// todo: add your logic here and delete this line

	return &timedjob.TaskLogIndexResp{}, nil
}
