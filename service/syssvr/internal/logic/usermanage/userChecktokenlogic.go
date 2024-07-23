package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/cache"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/users"
	"gitee.com/i-Things/share/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/maypok86/otter"
	"github.com/spf13/cast"
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
	switch in.AuthType {

	case "open":
		return l.openCheckToken(in)
	default:
		return l.userCheckToken(in)
	}

	return nil, errors.Parameter.AddMsg(in.AuthType)
}

var (
	openAuthCache otter.Cache[string, *sys.UserCheckTokenResp]
)

func init() {
	cache, err := otter.MustBuilder[string, *sys.UserCheckTokenResp](10_000).
		CollectStats().
		Cost(func(key string, value *sys.UserCheckTokenResp) uint32 {
			return 1
		}).
		WithTTL(time.Minute * 10).
		Build()
	logx.Must(err)
	openAuthCache = cache
}

func (l *CheckTokenLogic) openCheckToken(in *sys.UserCheckTokenReq) (*sys.UserCheckTokenResp, error) {
	v, ok := openAuthCache.Get(in.Token)
	if ok {
		return v, nil
	}
	var claim users.OpenClaims
	err := users.ParseTokenWithFunc(&claim, in.Token, func(token *jwt.Token) (interface{}, error) {
		if claim.TenantCode == "" || claim.UserID == 0 || claim.Code == "" {
			return nil, errors.TokenInvalid
		}
		po, err := relationDB.NewTenantOpenRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantOpenFilter{
			TenantCode: claim.TenantCode,
			UserID:     claim.UserID,
			Code:       claim.Code,
		})
		if err != nil {
			return nil, err
		}
		return []byte(po.AccessSecret), nil
	})
	if err != nil {
		l.Errorf("%s parse token fail err=%s", utils.FuncName(), err.Error())
		return nil, err
	}
	ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{
		UserIDs:    []int64{claim.UserID},
		WithRoles:  true,
		WithTenant: true,
	})
	if err != nil {
		l.Errorf("%s  err=%s", utils.FuncName(), err.Error())
		return nil, err
	}
	var rolses []int64
	var roleCodes []string
	var isAdmin int64 = def.False
	for _, v := range ui.Roles {
		rolses = append(rolses, v.RoleID)
		if v.Role != nil && v.Role.Code != "" {
			roleCodes = append(roleCodes, v.Role.Code)
		}
	}

	if ui.Tenant != nil && (utils.SliceIn(ui.Tenant.AdminRoleID, rolses...) || ui.Tenant.AdminUserID == ui.UserID) {
		isAdmin = def.True
	}
	var account = ui.UserName.String
	if account == "" {
		account = ui.Phone.String
	}
	if account == "" {
		account = ui.Email.String
	}
	if account == "" {
		account = cast.ToString(ui.UserID)
	}
	ret := sys.UserCheckTokenResp{UserID: claim.UserID, IsAllData: ui.IsAllData, RoleIDs: rolses, RoleCodes: roleCodes, IsSuperAdmin: utils.SliceIn(def.RoleCodeSupper, roleCodes...), IsAdmin: isAdmin == def.True, Account: account, TenantCode: claim.TenantCode}
	ret.IsAdmin = utils.SliceIn(def.RoleCodeAdmin, roleCodes...) || ret.IsSuperAdmin
	projectAuth, err := cache.GetProjectAuth(l.ctx, ret.UserID, ret.RoleIDs)
	if err != nil {
		return nil, err
	}
	ret.ProjectAuth = projectAuth
	openAuthCache.Set(in.Token, &ret)
	return &ret, nil
}

func (l *CheckTokenLogic) userCheckToken(in *sys.UserCheckTokenReq) (*sys.UserCheckTokenResp, error) {
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
	ui, err := l.svcCtx.UserTokenInfo.GetData(l.ctx, claim.UserID)
	if err != nil {
		l.Errorf("%s UserTokenInfo.GetData fail err=%s", utils.FuncName(), err.Error())
		return nil, err
	}
	ret := sys.UserCheckTokenResp{
		Token:        token,
		UserID:       claim.UserID,
		RoleIDs:      ui.RoleIDs,
		RoleCodes:    ui.RoleCodes,
		IsAllData:    ui.IsAllData,
		TenantCode:   ui.TenantCode,
		IsSuperAdmin: utils.SliceIn(def.RoleCodeSupper, ui.RoleCodes...),
		Account:      ui.Account,
	}
	ret.IsAdmin = utils.SliceIn(def.RoleCodeAdmin, ui.RoleCodes...) || ret.IsSuperAdmin
	projectAuth, err := cache.GetProjectAuth(l.ctx, ret.UserID, ret.RoleIDs)
	if err != nil {
		return nil, err
	}
	ret.ProjectAuth = projectAuth

	return &ret, nil
}
