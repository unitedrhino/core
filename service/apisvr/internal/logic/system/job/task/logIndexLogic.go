package task

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogIndexLogic {
	return &LogIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogIndexLogic) LogIndex(req *types.TimedTaskLogIndexReq) (resp *types.TimedTaskLogIndexResp, err error) {
	if req.Page.Size > 200 || req.Page.Size == 0 {
		req.Page.Size = 200
	}
	ret, err := l.svcCtx.TimedJob.TaskLogIndex(l.ctx, &timedjob.TaskLogIndexReq{
		Page:      logic.ToTimedJobPageRpc(req.Page),
		GroupCode: req.GroupCode,
		TaskCode:  req.TaskCode,
	})
	if err != nil {
		return nil, err
	}
	var list []*types.TimedTaskLog
	for _, v := range ret.List {
		list = append(list, ToTaskLog(v))
	}
	return &types.TimedTaskLogIndexResp{
		List:  list,
		Total: ret.Total,
	}, nil
}
