package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyConfigTemplateUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyConfigTemplateUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyConfigTemplateUpdateLogic {
	return &NotifyConfigTemplateUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 租户通知配置
func (l *NotifyConfigTemplateUpdateLogic) NotifyConfigTemplateUpdate(in *sys.NotifyConfigTemplate) (*sys.Empty, error) {
	po := relationDB.SysNotifyConfigTemplate{
		NotifyCode: in.NotifyCode,
		Type:       in.Type,
		TemplateID: in.TemplateID,
	}
	err := relationDB.NewNotifyConfigTemplateRepo(l.ctx).Save(l.ctx, &po)
	return &sys.Empty{}, err
}
