package notifymanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

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
	po, err := relationDB.NewNotifyTemplateRepo(l.ctx).FindOne(l.ctx, in.Id)
	return utils.Copy[sys.NotifyTemplate](po), err
}
