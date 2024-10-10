package notifymanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyChannelUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyChannelUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyChannelUpdateLogic {
	return &NotifyChannelUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyChannelUpdateLogic) NotifyChannelUpdate(in *sys.NotifyChannel) (*sys.Empty, error) {
	old, err := relationDB.NewNotifyChannelRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		old.Name = in.Name
	}
	if in.Email != nil {
		old.Email = utils.Copy[relationDB.SysTenantEmail](in.Email)
	}
	if in.Webhook != "" {
		old.WebHook = in.Webhook
	}
	if in.App != nil {
		old.App = utils.Copy[relationDB.SysTenantThird](in.App)
	}
	if in.Sms != nil {
		old.Sms = utils.Copy[relationDB.SysSms](in.Sms)
	}
	old.Desc = in.Desc
	err = relationDB.NewNotifyChannelRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
