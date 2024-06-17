package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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
	err := relationDB.NewNotifyConfigRepo(l.ctx).Delete(l.ctx, in.Id)
	//todo 后续需要做好是否可以删除的检测
	return &sys.Empty{}, err
}
