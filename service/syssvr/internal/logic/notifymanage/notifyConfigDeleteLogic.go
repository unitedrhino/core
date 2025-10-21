package notifymanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyConfigDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyConfigDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyConfigDeleteLogic {
	return &NotifyConfigDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyConfigDeleteLogic) NotifyConfigDelete(in *sys.WithID) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	err := relationDB.NewNotifyConfigRepo(l.ctx).Delete(l.ctx, in.Id)
	return &sys.Empty{}, err
}
