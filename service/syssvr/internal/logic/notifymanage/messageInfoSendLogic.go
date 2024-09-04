package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"time"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageInfoSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessageInfoSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageInfoSendLogic {
	return &MessageInfoSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 站内信
func (l *MessageInfoSendLogic) MessageInfoSend(in *sys.MessageInfoSendReq) (*sys.WithID, error) {
	ni, err := relationDB.NewNotifyConfigRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.NotifyConfigFilter{Code: in.NotifyCode})
	if err != nil {
		return nil, err
	}

	po := relationDB.SysMessageInfo{
		Group:      ni.Group,
		NotifyCode: ni.Code,
		Subject:    in.Subject,
		Body:       in.Body,
		Str1:       in.Str1,
		Str2:       in.Str2,
		Str3:       in.Str3,
		IsGlobal:   in.IsGlobal,
	}
	if in.NotifyTime != 0 {
		po.NotifyTime = time.Unix(in.NotifyTime, 0)
	}

	if in.IsGlobal == def.True {
		err := relationDB.NewMessageInfoRepo(l.ctx).Insert(l.ctx, &po)
		if err != nil {
			return nil, err
		}
		return &sys.WithID{Id: po.ID}, nil
	}
	if len(in.UserIDs) == 0 {
		return nil, errors.Parameter.AddMsg("need userIDs")
	}
	userNum, err := relationDB.NewUserInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.UserInfoFilter{UserIDs: in.UserIDs})
	if err != nil {
		return nil, err
	}
	if userNum != int64(len(in.UserIDs)) {
		return nil, errors.Parameter.AddMsg("需要填写正确的用户ID")
	}

	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewMessageInfoRepo(l.ctx).Insert(l.ctx, &po)
		if err != nil {
			return err
		}
		var users []*relationDB.SysUserMessage
		for _, v := range in.UserIDs {
			users = append(users, &relationDB.SysUserMessage{
				UserID:    v,
				Group:     po.Group,
				MessageID: po.ID,
				IsRead:    def.False,
			})
		}
		err = relationDB.NewUserMessageRepo(l.ctx).MultiInsert(l.ctx, users)
		return err
	})
	return &sys.WithID{Id: po.ID}, err
}
