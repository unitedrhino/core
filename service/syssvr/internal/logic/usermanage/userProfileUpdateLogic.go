package usermanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserProfileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileUpdateLogic {
	return &UserProfileUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserProfileUpdateLogic) UserProfileUpdate(in *sys.UserProfile) (*sys.Empty, error) {
	userID := ctxs.GetUserCtxNoNil(l.ctx).UserID
	old, err := relationDB.NewUserProfileRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserProfileFilter{
		Code:   in.Code,
		UserID: userID,
	})
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
		old = &relationDB.SysUserProfile{UserID: userID, Code: in.Code}
	}
	old.Params = in.Params
	err = relationDB.NewUserProfileRepo(l.ctx).Update(l.ctx, old)
	return &sys.Empty{}, err
}
