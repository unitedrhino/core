package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyConfigTemplateDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyConfigTemplateDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyConfigTemplateDeleteLogic {
	return &NotifyConfigTemplateDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyConfigTemplateDeleteLogic) NotifyConfigTemplateDelete(in *sys.NotifyConfigTemplateDeleteReq) (*sys.Empty, error) {
	err := relationDB.NewNotifyConfigTemplateRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.NotifyConfigTemplateFilter{
		NotifyCode: in.NotifyCode,
		Type:       in.Type,
	})
	return &sys.Empty{}, err
}
