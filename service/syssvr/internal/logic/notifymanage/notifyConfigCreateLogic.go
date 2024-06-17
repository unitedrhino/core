package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyConfigCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyConfigCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyConfigCreateLogic {
	return &NotifyConfigCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyConfigCreateLogic) NotifyConfigCreate(in *sys.NotifyConfig) (*sys.WithID, error) {
	po := utils.Copy[relationDB.SysNotifyConfig](in)
	po.ID = 0
	err := relationDB.NewNotifyConfigRepo(l.ctx).Insert(l.ctx, po)

	return &sys.WithID{Id: po.ID}, err
}
