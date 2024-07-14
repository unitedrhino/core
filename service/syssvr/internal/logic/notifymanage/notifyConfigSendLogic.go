package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyConfigSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyConfigSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyConfigSendLogic {
	return &NotifyConfigSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyConfigSendLogic) NotifyConfigSend(in *sys.NotifyConfigSendReq) (*sys.Empty, error) {
	err := SendNotifyMsg(l.ctx, l.svcCtx, SendMsgConfig{
		UserIDs:    in.UserIDs,
		Accounts:   in.Accounts,
		NotifyCode: in.NotifyCode,
		TemplateID: in.TemplateID,
		Type:       in.Type,
		Params:     utils.ToStringMap(in.Params),
		Str1:       in.Str1,
		Str2:       in.Str2,
		Str3:       in.Str3,
	})

	return &sys.Empty{}, err
}
