package timedmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskGroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskGroupCreateLogic {
	return &TaskGroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskGroupCreateLogic) TaskGroupCreate(in *timedjob.TaskGroup) (*timedjob.Response, error) {
	err := relationDB.NewTaskGroupRepo(l.ctx).Insert(l.ctx, ToTaskGroupPo(in))
	if err != nil {
		return nil, err
	}
	return &timedjob.Response{}, nil
}
