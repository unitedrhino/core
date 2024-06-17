package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyChannelReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyChannelReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyChannelReadLogic {
	return &NotifyChannelReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyChannelReadLogic) NotifyChannelRead(in *sys.WithID) (*sys.NotifyChannel, error) {
	po, err := relationDB.NewNotifyChannelRepo(l.ctx).FindOne(l.ctx, in.Id)
	return utils.Copy[sys.NotifyChannel](po), err
}
