package task

import (
	"context"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoCreateLogic {
	return &InfoCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoCreateLogic) InfoCreate(req *types.TimedTaskInfo) error {
	l.Infof("req:%v", utils.Fmt(req))
	_, err := l.svcCtx.TimedJob.TaskInfoCreate(l.ctx, ToTaskInfoPb(req))
	return err
}
