package notifymanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessageInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageInfoDeleteLogic {
	return &MessageInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MessageInfoDeleteLogic) MessageInfoDelete(in *sys.WithID) (*sys.Empty, error) {
	if in.Id == 0 {
		return nil, errors.Parameter.AddMsg("id 必填")
	}
	err := stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewMessageInfoRepo(l.ctx).Delete(l.ctx, in.Id)
		if err != nil {
			return err
		}
		err = relationDB.NewUserMessageRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.UserMessageFilter{MessageID: in.Id})
		return err
	})
	return &sys.Empty{}, err
}
