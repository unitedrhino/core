package usermanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserProfileReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileReadLogic {
	return &UserProfileReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserProfileReadLogic) UserProfileRead(in *sys.WithCode) (*sys.UserProfile, error) {
	ret, err := relationDB.NewUserProfileRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserProfileFilter{
		Code:   in.Code,
		UserID: ctxs.GetUserCtxNoNil(l.ctx).UserID,
	})
	if errors.Cmp(err, errors.NotFind) {
		return &sys.UserProfile{
			Code:   in.Code,
			Params: "",
		}, nil
	}
	return utils.Copy[sys.UserProfile](ret), err
}
