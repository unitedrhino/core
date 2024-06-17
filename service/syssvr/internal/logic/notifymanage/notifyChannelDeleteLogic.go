package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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
	err := relationDB.NewNotifyChannelRepo(l.ctx).Delete(l.ctx, in.Id)
	return &sys.Empty{}, err
}
