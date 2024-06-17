package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/cache"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/users"
	"gitee.com/i-Things/share/utils"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserCheckTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckTokenLogic {
	return &CheckTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *CheckTokenLogic) UserCheckToken(in *sys.UserCheckTokenReq) (*sys.UserCheckTokenResp, error) {
	var claim users.LoginClaims
	err := users.ParseToken(&claim, in.Token, l.svcCtx.Config.UserToken.AccessSecret)
	if err != nil {
		l.Errorf("%s parse token fail err=%s", utils.FuncName(), err.Error())
		return nil, err
	}
	var token string

	if (claim.ExpiresAt.Unix()-time.Now().Unix())*2 < l.svcCtx.Config.UserToken.AccessExpire {
		token, _ = users.RefreshLoginToken(in.Token, l.svcCtx.Config.UserToken.AccessSecret, time.Now().Unix()+l.svcCtx.Config.UserToken.AccessExpire)
	}
	ret := sys.UserCheckTokenResp{
		Token:        token,
		UserID:       claim.UserID,
		RoleIDs:      claim.RoleIDs,
		RoleCodes:    claim.RoleCodes,
		IsAllData:    claim.IsAllData,
		TenantCode:   claim.TenantCode,
		IsSuperAdmin: utils.SliceIn(def.RoleCodeSupper, claim.RoleCodes...),
		Account:      claim.Account,
	}
	ret.IsAdmin = utils.SliceIn(def.RoleCodeAdmin, claim.RoleCodes...) || ret.IsSuperAdmin
	if !ret.IsAdmin {
		projectAuth, err := cache.GetProjectAuth(l.ctx, ret.UserID, ret.RoleIDs)
		if err != nil {
			return nil, err
		}
		ret.ProjectAuth = projectAuth
	}

	return &ret, nil
}
