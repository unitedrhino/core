package timedmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/timed/internal/repo/relationDB"
	"gitee.com/i-Things/share/stores"

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
	db := relationDB.NewJobLogRepo(l.ctx)
	f := relationDB.TaskLogFilter{
		GroupCode: in.GroupCode,
		TaskCode:  in.TaskCode,
	}
	total, err := db.CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	pos, err := db.FindByFilter(l.ctx, f, ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{
		Field: "createdTime",
		Sort:  2,
	}))
	if err != nil {
		return nil, err
	}
	var list []*timedjob.TaskLog
	for _, v := range pos {
		list = append(list, ToTaskLog(v))
	}
	return &timedjob.TaskLogIndexResp{Total: total, List: list}, nil
}
