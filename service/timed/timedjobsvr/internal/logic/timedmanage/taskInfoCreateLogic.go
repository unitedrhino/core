package timedmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"

	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskInfoCreateLogic {
	return &TaskInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskInfoCreateLogic) TaskInfoCreate(in *timedjob.TaskInfo) (*timedjob.Response, error) {
	po := ToTaskInfoPo(in)
	err := relationDB.NewTaskInfoRepo(l.ctx).Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	return &timedjob.Response{}, nil
}
