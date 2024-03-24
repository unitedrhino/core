package timedmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/timed/internal/domain"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskCancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskCancelLogic {
	return &TaskCancelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var (
	retryKeys []string
)

func init() {
	for _, priority := range domain.Priorities {
		retryKeys = append(retryKeys, fmt.Sprintf("asynq:{%s}:retry", priority))
	}
}

func (l *TaskCancelLogic) TaskCancel(in *timedjob.TaskWithTaskID) (*timedjob.Response, error) {
	for _, priority := range domain.Priorities {
		err := l.svcCtx.AsynqInspector.DeleteTask(priority, in.TaskID)
		if err == nil {
			break
		}
	}
	return &timedjob.Response{}, nil
}
