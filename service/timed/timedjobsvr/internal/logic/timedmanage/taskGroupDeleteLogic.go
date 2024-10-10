package timedmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"

	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskGroupDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskGroupDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskGroupDeleteLogic {
	return &TaskGroupDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskGroupDeleteLogic) TaskGroupDelete(in *timedjob.WithCode) (*timedjob.Response, error) {
	err := relationDB.NewTaskGroupRepo(l.ctx).DeleteByFilter(l.ctx,
		relationDB.TaskGroupFilter{Codes: []string{in.Code}})
	if err != nil {
		return nil, err
	}
	return &timedjob.Response{}, nil
}
