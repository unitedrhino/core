package notifymanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/stores"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

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
	if err != nil {
		return nil, err
	}
	err = InitConfigEnableTypes(l.ctx, stores.GetTenantConn(l.ctx), in.NotifyCode)
	return &sys.Empty{}, err
}
