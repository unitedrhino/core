package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyConfigUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyConfigUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyConfigUpdateLogic {
	return &NotifyConfigUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyConfigUpdateLogic) NotifyConfigUpdate(in *sys.NotifyConfig) (*sys.Empty, error) {
	db := relationDB.NewNotifyConfigRepo(l.ctx)
	old, err := db.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Group != "" {
		old.Group = in.Group
	}
	if in.Name != "" {
		old.Name = in.Name
	}
	if len(in.SupportTypes) != 0 {
		old.SupportTypes = in.SupportTypes
	}
	if in.Desc != "" {
		old.Desc = in.Desc
	}
	if in.IsRecord != 0 {
		old.IsRecord = in.IsRecord
	}
	if in.Params != nil {
		old.Params = in.Params
	}
	err = db.Update(l.ctx, old)
	return &sys.Empty{}, err
}
