package processTask

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/internal/domain"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/logic/processTask/sqlFunc"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessTask struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	TaskSend sqlFunc.TaskSendFunc
}

func NewProcessTask(ctx context.Context, svcCtx *svc.ServiceContext, TaskSend sqlFunc.TaskSendFunc) *ProcessTask {
	return &ProcessTask{
		ctx:      ctx,
		svcCtx:   svcCtx,
		Logger:   logx.WithContext(ctx),
		TaskSend: TaskSend,
	}
}

func (t ProcessTask) Process(ctx context.Context, taskInfo domain.TaskInfo) error {
	switch taskInfo.GroupType {
	case domain.TaskGroupTypeQueue:
		return t.Queue(ctx, &taskInfo)
	case domain.TaskGroupTypeSql:
		return t.SqlExec(ctx, &taskInfo)
	case domain.TaskGroupTypeScript:
		return t.ScriptExec(ctx, &taskInfo)
	default:
		logx.WithContext(ctx).Errorf("not support job type:%v", taskInfo.GroupType)
		return errors.Parameter.AddMsgf("not support job type:%v", taskInfo.GroupType)
	}
}
