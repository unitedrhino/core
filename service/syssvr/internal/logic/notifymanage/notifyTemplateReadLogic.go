package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyTemplateReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyTemplateReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyTemplateReadLogic {
	return &NotifyTemplateReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通知模版
func (l *NotifyTemplateReadLogic) NotifyTemplateRead(in *sys.WithID) (*sys.NotifyTemplate, error) {
	l.ctx = ctxs.WithCommonTenant(l.ctx)
	po, err := relationDB.NewNotifyTemplateRepo(l.ctx).FindOne(l.ctx, in.Id)
	return utils.Copy[sys.NotifyTemplate](po), err
}
