package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyInfoCreateLogic {
	return &NotifyInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyInfoCreateLogic) NotifyInfoCreate(in *sys.NotifyInfo) (*sys.WithID, error) {
	po := utils.Copy[relationDB.SysNotifyInfo](in)
	po.ID = 0
	err := relationDB.NewNotifyInfoRepo(l.ctx).Insert(l.ctx, po)

	return &sys.WithID{Id: po.ID}, err
}
