package task

import (
	"context"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoDeleteLogic {
	return &InfoDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoDeleteLogic) InfoDelete(req *types.WithCode) error {
	l.Infof("req:%v", utils.Fmt(req))
	_, err := l.svcCtx.TimedJob.TaskInfoDelete(l.ctx, &timedjob.WithCode{Code: req.Code})
	return err
}
