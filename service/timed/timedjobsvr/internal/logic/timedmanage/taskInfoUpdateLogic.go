package timedmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"

	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskInfoUpdateLogic {
	return &TaskInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskInfoUpdateLogic) TaskInfoUpdate(in *timedjob.TaskInfo) (*timedjob.Response, error) {
	repo := relationDB.NewTaskInfoRepo(l.ctx)
	oldPo, err := repo.FindOneByFilter(l.ctx, relationDB.TaskFilter{Codes: []string{in.Code}})
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		oldPo.Name = in.Name
	}
	if in.Priority != 0 {
		oldPo.Priority = in.Priority
	}
	if in.Params != "" {
		oldPo.Params = in.Params
	}
	if in.CronExpr != "" {
		oldPo.CronExpr = in.CronExpr
	}
	if in.Status != 0 {
		oldPo.Status = in.Status
	}
	if in.Topics != nil {
		oldPo.Topics = in.Topics
	}
	if in.Priority != 0 {
		oldPo.Priority = in.Priority
	}
	err = repo.Update(l.ctx, oldPo)
	if err != nil {
		return nil, err
	}
	return &timedjob.Response{}, nil
}
