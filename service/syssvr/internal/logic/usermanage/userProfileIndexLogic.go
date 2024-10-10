package usermanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserProfileIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileIndexLogic {
	return &UserProfileIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserProfileIndexLogic) UserProfileIndex(in *sys.UserProfileIndexReq) (*sys.UserProfileIndexResp, error) {
	ret, err := relationDB.NewUserProfileRepo(l.ctx).FindByFilter(l.ctx, relationDB.UserProfileFilter{
		Codes:  in.Codes,
		UserID: ctxs.GetUserCtxNoNil(l.ctx).UserID,
	}, nil)
	return &sys.UserProfileIndexResp{Profiles: utils.CopySlice[sys.UserProfile](ret)}, err
}
