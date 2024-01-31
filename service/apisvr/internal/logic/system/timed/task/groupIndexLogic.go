package task

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/i-Things/core/shared/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupIndexLogic {
	return &GroupIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupIndexLogic) GroupIndex(req *types.TimedTaskGroupIndexReq) (resp *types.TimedTaskGroupIndexResp, err error) {
	l.Infof("req:%v", utils.Fmt(req))
	ret, err := l.svcCtx.TimedJob.TaskGroupIndex(l.ctx, &timedjob.TaskGroupIndexReq{Page: logic.ToTimedJobPageRpc(req.Page)})
	if err != nil {
		return nil, err
	}
	return &types.TimedTaskGroupIndexResp{List: ToTaskGroupsTypes(ret.List), Total: ret.Total}, nil
}
