package messagemanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifySendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifySendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifySendLogic {
	return &NotifySendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifySendLogic) NotifySend(in *sys.NotifySendReq) (*sys.Empty, error) {
	err := SendNotifyMsg(l.ctx, l.svcCtx, SendMsgConfig{
		UserIDs:    in.UserIDs,
		NotifyCode: in.NotifyCode,
		Type:       in.Type,
		Params:     utils.ToStringMap(in.Params),
		Str1:       in.Str1,
		Str2:       in.Str2,
		Str3:       in.Str3,
	})

	return &sys.Empty{}, err
}
