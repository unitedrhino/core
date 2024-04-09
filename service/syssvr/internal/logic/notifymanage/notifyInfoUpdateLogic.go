package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyInfoUpdateLogic {
	return &NotifyInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyInfoUpdateLogic) NotifyInfoUpdate(in *sys.NotifyInfo) (*sys.Empty, error) {
	db := relationDB.NewNotifyInfoRepo(l.ctx)
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
	if in.DefaultSubject != "" {
		old.DefaultSubject = in.DefaultSubject
	}
	if in.DefaultBody != "" {
		old.DefaultBody = in.DefaultBody
	}
	if in.DefaultTemplateCode != "" {
		old.DefaultTemplateCode = in.DefaultTemplateCode
	}
	if in.DefaultSignName != "" {
		old.DefaultSignName = in.DefaultSignName
	}
	if in.Params != nil {
		old.Params = in.Params
	}
	err = db.Update(l.ctx, old)
	return &sys.Empty{}, err
}
