package task

import (
	"context"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoReadLogic {
	return &InfoReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoReadLogic) InfoRead(req *types.WithGroupCode) (resp *types.TimedTaskInfo, err error) {
	l.Infof("req:%v", utils.Fmt(req))
	ret, err := l.svcCtx.TimedJob.TaskInfoRead(l.ctx, &timedjob.WithGroupCode{Code: req.Code, GroupCode: req.GroupCode})
	return ToTaskInfoTypes(ret), err
}
