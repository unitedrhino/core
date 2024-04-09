package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessageInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageInfoUpdateLogic {
	return &MessageInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MessageInfoUpdateLogic) MessageInfoUpdate(in *sys.MessageInfo) (*sys.Empty, error) {
	db := relationDB.NewMessageInfoRepo(l.ctx)
	old, err := db.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Subject != "" {
		old.Subject = in.Subject
	}
	if in.Body != "" {
		old.Body = in.Body
	}
	if in.Str1 != "" {
		old.Str1 = in.Str1
	}
	if in.Str2 != "" {
		old.Str2 = in.Str2
	}
	if in.Str3 != "" {
		old.Str3 = in.Str3
	}
	err = db.Update(l.ctx, old)
	if err != nil {
		return nil, err
	}
	return &sys.Empty{}, nil
}
