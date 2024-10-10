package timedmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"

	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskGroupReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskGroupReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskGroupReadLogic {
	return &TaskGroupReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskGroupReadLogic) TaskGroupRead(in *timedjob.WithCode) (*timedjob.TaskGroup, error) {
	po, err := relationDB.NewTaskGroupRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TaskGroupFilter{Codes: []string{in.Code}})
	if err != nil {
		return nil, err
	}
	return ToTaskGroupPb(po), nil
}
