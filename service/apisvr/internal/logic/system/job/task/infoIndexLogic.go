package task

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoIndexLogic {
	return &InfoIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoIndexLogic) InfoIndex(req *types.TimedTaskInfoIndexReq) (resp *types.TimedTaskInfoIndexResp, err error) {
	l.Infof("req:%v", utils.Fmt(req))
	ret, err := l.svcCtx.TimedJob.TaskInfoIndex(l.ctx, &timedjob.TaskInfoIndexReq{Page: logic.ToTimedJobPageRpc(req.Page), GroupCode: req.GroupCode})
	if err != nil {
		return nil, err
	}
	return &types.TimedTaskInfoIndexResp{List: ToTaskInfosTypes(ret.List), Total: ret.Total}, nil
}
