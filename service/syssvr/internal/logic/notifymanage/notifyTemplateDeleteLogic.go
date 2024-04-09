package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyTemplateDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyTemplateDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyTemplateDeleteLogic {
	return &NotifyTemplateDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyTemplateDeleteLogic) NotifyTemplateDelete(in *sys.WithID) (*sys.Empty, error) {
	err := relationDB.NewNotifyTemplateRepo(l.ctx).Delete(l.ctx, in.Id)

	return &sys.Empty{}, err
}
