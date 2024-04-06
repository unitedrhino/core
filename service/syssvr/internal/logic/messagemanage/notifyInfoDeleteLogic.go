package messagemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyInfoDeleteLogic {
	return &NotifyInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyInfoDeleteLogic) NotifyInfoDelete(in *sys.WithID) (*sys.Empty, error) {
	err:=relationDB.NewNotifyInfoRepo(l.ctx).Delete(l.ctx,in.Id)
	//todo 后续需要做好是否可以删除的检测
	return &sys.Empty{}, err
}
