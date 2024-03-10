package timedmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/timed/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskInfoDeleteLogic {
	return &TaskInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskInfoDeleteLogic) TaskInfoDelete(in *timedjob.WithGroupCode) (*timedjob.Response, error) {
	err := relationDB.NewTaskInfoRepo(l.ctx).DeleteByFilter(l.ctx,
		relationDB.TaskFilter{Codes: []string{in.Code}, GroupCode: in.GroupCode})
	if err != nil {
		return nil, err
	}
	return &timedjob.Response{}, nil
}
