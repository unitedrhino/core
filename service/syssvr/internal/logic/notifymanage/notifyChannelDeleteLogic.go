package notifymanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyChannelDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyChannelDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyChannelDeleteLogic {
	return &NotifyChannelDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyChannelDeleteLogic) NotifyChannelDelete(in *sys.WithID) (*sys.Empty, error) {
	old, err := relationDB.NewNotifyChannelRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if !ctxs.CanHandTenant(l.ctx, old.TenantCode) || ctxs.IsAdmin(l.ctx) != nil {
		return &sys.Empty{}, errors.Permissions
	}

	err = relationDB.NewNotifyChannelRepo(l.ctx).Delete(l.ctx, in.Id)
	return &sys.Empty{}, err
}
