package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyConfigReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyConfigReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyConfigReadLogic {
	return &NotifyConfigReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通知配置信息
func (l *NotifyConfigReadLogic) NotifyConfigRead(in *sys.WithIDCode) (*sys.NotifyConfig, error) {
	po, err := relationDB.NewNotifyConfigRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.NotifyConfigFilter{
		Code: in.Code,
		ID:   in.Id,
	})
	return utils.Copy[sys.NotifyConfig](po), err
}
