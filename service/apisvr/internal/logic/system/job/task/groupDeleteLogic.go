package task

import (
	"context"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeleteLogic {
	return &GroupDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupDeleteLogic) GroupDelete(req *types.WithCode) error {
	l.Infof("req:%v", utils.Fmt(req))
	_, err := l.svcCtx.TimedJob.TaskGroupDelete(l.ctx, &timedjob.WithCode{Code: req.Code})
	return err
}
