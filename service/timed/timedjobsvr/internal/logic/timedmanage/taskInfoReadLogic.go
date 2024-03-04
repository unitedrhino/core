package timedmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/timed/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskInfoReadLogic {
	return &TaskInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TaskInfoReadLogic) TaskInfoRead(in *timedjob.WithCode) (*timedjob.TaskInfo, error) {
	po, err := relationDB.NewTaskInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TaskFilter{Codes: []string{in.Code}})
	if err != nil {
		return nil, err
	}
	return ToTaskInfoPb(po), nil
}
