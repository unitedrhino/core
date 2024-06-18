package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyTemplateUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyTemplateUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyTemplateUpdateLogic {
	return &NotifyTemplateUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyTemplateUpdateLogic) NotifyTemplateUpdate(in *sys.NotifyTemplate) (*sys.Empty, error) {
	db := relationDB.NewNotifyTemplateRepo(l.ctx)
	old, err := db.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		old.Name = in.Name
	}
	if in.Type != "" {
		old.Type = in.Type
	}
	if in.SignName != "" {
		old.SignName = in.SignName
	}
	if in.TemplateCode != "" {
		old.TemplateCode = in.TemplateCode
	}
	if in.Subject != "" {
		old.Subject = in.Subject
	}
	if in.Body != "" {
		old.Body = in.Body
	}
	if in.Desc != "" {
		old.Desc = in.Desc
	}
	if in.ChannelID != 0 {
		old.ChannelID = in.ChannelID
	}
	err = db.Update(l.ctx, old)
	return &sys.Empty{}, err
}
