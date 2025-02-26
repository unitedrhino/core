package task

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupReadLogic {
	return &GroupReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupReadLogic) GroupRead(req *types.WithCode) (resp *types.TimedTaskGroup, err error) {
	l.Infof("req:%v", utils.Fmt(req))
	ret, err := l.svcCtx.TimedJob.TaskGroupRead(l.ctx, &timedjob.WithCode{Code: req.Code})
	return ToGroupTypes(ret), err
}
