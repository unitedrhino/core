package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyInfoReadLogic {
	return &NotifyInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通知配置信息
func (l *NotifyInfoReadLogic) NotifyInfoRead(in *sys.WithIDCode) (*sys.NotifyInfo, error) {
	po, err := relationDB.NewNotifyInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.NotifyInfoFilter{
		Code: in.Code,
		ID:   in.Id,
	})
	return utils.Copy[sys.NotifyInfo](po), err
}
