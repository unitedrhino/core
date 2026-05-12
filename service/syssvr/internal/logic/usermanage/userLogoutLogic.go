package usermanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/users"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogoutLogic {
	return &UserLogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserLogoutLogic) UserLogout(in *sys.UserLogoutReq) (*sys.Empty, error) {
	var claim users.LoginClaims
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if in.Token == "" {
		in.Token = uc.Token
	}
	err := users.ParseToken(&claim, in.Token, l.svcCtx.Config.UserToken.AccessSecret)
	if err != nil {
		l.Errorf("%s parse token fail err=%s", utils.FuncName(), err.Error())
		return nil, err
	}

	err = l.svcCtx.UserToken.Logout(l.ctx, claim)
	if err != nil {
		return nil, err
	}
	l.svcCtx.UserCache.SetData(l.ctx, claim.UserID, nil)
	l.svcCtx.UsersCache.SetData(l.ctx, claim.UserID, nil)

	return &sys.Empty{}, nil
}
