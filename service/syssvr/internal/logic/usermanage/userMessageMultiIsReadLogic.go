package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserMessageMultiIsReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserMessageMultiIsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserMessageMultiIsReadLogic {
	return &UserMessageMultiIsReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserMessageMultiIsReadLogic) UserMessageMultiIsRead(in *sys.IDList) (*sys.Empty, error) {
	var userID = ctxs.GetUserCtxNoNil(l.ctx).UserID
	err := relationDB.NewUserMessageRepo(l.ctx).MultiIsRead(l.ctx, userID, in.Ids)
	if err != nil {
		return nil, err
	}
	UpdateUserNotRead(l.ctx, userID)
	return &sys.Empty{}, err
}

func UpdateUserNotRead(ctx context.Context, userID int64) (err error) {
	var count = map[string]int64{}
	count, err = relationDB.NewUserMessageRepo(ctx).CountNotRead(ctx, userID)
	if err != nil {
		return err
	}
	err = relationDB.NewUserInfoRepo(ctx).UpdateMessageNotRead(ctx, userID, count)
	if err != nil {
		return err
	}
	return err
}
