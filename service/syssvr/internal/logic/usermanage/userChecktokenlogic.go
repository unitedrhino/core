package usermanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/cache"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/users"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
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
		WithTTL(time.Minute * 1).
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
		po, err := relationDB.NewDataOpenAccessRepo(l.ctx).FindOneByFilter(ctxs.WithRoot(l.ctx), relationDB.DataOpenAccessFilter{
			TenantCode: claim.TenantCode,
			UserID:     claim.UserID,
			Code:       claim.Code,
		})
		if err != nil {
			return nil, err
		}
		if len(po.IpRange) != 0 { //ip校验
			var match bool
			for _, whiteIp := range po.IpRange {
				if utils.MatchIP(in.Ip, whiteIp) {
					match = true
					break
				}
			}
			if !match {
				return nil, errors.Permissions
			}
		}
		return []byte(po.AccessSecret), nil
	})
	if err != nil {
		l.Errorf("%s parse token fail err=%s", utils.FuncName(), err.Error())
		return nil, err
	}
	ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(ctxs.WithRoot(l.ctx), relationDB.UserInfoFilter{
		TenantCode: claim.TenantCode,
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
	ret := sys.UserCheckTokenResp{UserID: claim.UserID, IsAllData: ui.IsAllData, RoleIDs: rolses, RoleCodes: roleCodes,
		IsSuperAdmin: utils.SliceIn(def.RoleCodeSupper, roleCodes...) || (isAdmin == def.True),
		Account:      account, TenantCode: claim.TenantCode}
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
	if claim.DeviceID != "" && claim.DeviceID != in.DeviceID {
		return nil, errors.TokenInvalid.AddMsg("token不可以跨设备使用")
	}
	var token string

	ui, err := l.svcCtx.UsersCache.GetData(l.ctx, claim.UserID)
	if err != nil {
		l.Errorf("%s UsersCache.GetData fail err=%s", utils.FuncName(), err.Error())
		return nil, err
	}
	tc, err := l.svcCtx.TenantConfigCache.GetData(l.ctx, ui.TenantCode)
	if err != nil {
		l.Errorf("%s  err=%s", utils.FuncName(), err.Error())
		return nil, err
	}
	if tc.IsSsl == def.True {
		err := l.svcCtx.UserToken.CheckToken(l.ctx, claim)
		if err != nil {
			return nil, err
		}
	}
	if (claim.ExpiresAt.Unix()-time.Now().Unix())*2 < l.svcCtx.Config.UserToken.AccessExpire {
		token, _ = users.RefreshLoginToken(in.Token, l.svcCtx.Config.UserToken.AccessSecret, l.svcCtx.Config.UserToken.AccessExpire)
	}
	ti, err := l.svcCtx.TenantCache.GetData(l.ctx, ui.TenantCode)
	if err != nil {
		l.Errorf("%s  err=%s", utils.FuncName(), err.Error())
		return nil, err
	}
	ret := sys.UserCheckTokenResp{
		Token:        token,
		AppCode:      claim.AppCode,
		UserID:       claim.UserID,
		RoleIDs:      ui.RoleIDs,
		RoleCodes:    ui.RoleCodes,
		IsAllData:    ui.IsAllData,
		TenantCode:   ui.TenantCode,
		IsSuperAdmin: ti.AdminUserID == ui.UserID || utils.SliceIn(ti.AdminRoleID, ui.RoleIDs...),
		Account:      ui.Account,
	}
	ret.IsAdmin = utils.SliceIn(def.RoleCodeAdmin, ui.RoleCodes...) || ret.IsSuperAdmin
	projectAuth, err := cache.GetProjectAuth(ctxs.BindTenantCode(l.ctx, ui.TenantCode, 0), ret.UserID, ret.RoleIDs)
	if err != nil {
		return nil, err
	}
	ret.ProjectAuth = projectAuth

	return &ret, nil
}
