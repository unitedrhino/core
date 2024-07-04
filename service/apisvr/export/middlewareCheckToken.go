package export

import (
	"context"
	role "gitee.com/i-Things/core/service/syssvr/client/rolemanage"
	tenant "gitee.com/i-Things/core/service/syssvr/client/tenantmanage"
	user "gitee.com/i-Things/core/service/syssvr/client/usermanage"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"strings"
)

type CheckTokenWareMiddleware struct {
	UserRpc   user.UserManage
	AuthRpc   role.RoleManage
	TenantRpc tenant.TenantManage
}

func NewCheckTokenWareMiddleware(UserRpc user.UserManage, AuthRpc role.RoleManage, TenantRpc tenant.TenantManage) *CheckTokenWareMiddleware {
	return &CheckTokenWareMiddleware{UserRpc: UserRpc, AuthRpc: AuthRpc, TenantRpc: TenantRpc}
}

func (m *CheckTokenWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			userCtx  *ctxs.UserCtx
			err      error
			isOpen   bool
			token    string
			strIP, _ = utils.GetIP(r)
			authType = "user"
		)
		authHeader := ctxs.GetHandle(r, "Authorization")
		// 检查"Authorization"字段是否存在并且以"Bearer "为前缀
		if strings.HasPrefix(authHeader, "Bearer ") {
			authType = "open"
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			token = ctxs.GetHandle(r, ctxs.UserTokenKey, ctxs.UserToken2Key)
			if token == "" {
				logx.WithContext(r.Context()).Errorf("%s.CheckTokenWare ip=%s not find token",
					utils.FuncName(), strIP)
				result.HttpErr(w, r, http.StatusUnauthorized, errors.NotLogin.AddMsg("用户请求失败"))
				return
			}
			authType = "user"
		}
		userCtx, err = m.Auth(r.Context(), w, token, strIP, authType)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("%s.UserAuth error=%s", utils.FuncName(), err)
			result.HttpErr(w, r, http.StatusUnauthorized, errors.Fmt(err).AddMsg("认证失败"))
			return
		}
		//注入 用户信息 到 ctx
		ctx2 := ctxs.SetUserCtx(r.Context(), userCtx)
		r = r.WithContext(ctx2)
		if !isOpen {
			////校验 Casbin Rule
			_, err = m.AuthRpc.RoleApiAuth(r.Context(), &user.RoleApiAuthReq{
				//RoleID: userCtx.RoleID, todo
				Path:   r.URL.Path,
				Method: r.Method,
			})
			if err != nil {
				logx.WithContext(r.Context()).Errorf("%s.AuthApiCheck error=%s", utils.FuncName(), err)
				//http.Error(w, "接口权限不足："+err.Error(), http.StatusUnauthorized)
				//return
			}
		}

		next(w, r)
	}
}

func (m *CheckTokenWareMiddleware) OpenAuth(r *http.Request, token string) (*ctxs.UserCtx, error) {
	strIP, _ := utils.GetIP(r)
	resp, err := m.TenantRpc.TenantOpenCheckToken(r.Context(), &sys.TenantOpenCheckTokenReq{
		Token: token,
		Ip:    strIP,
	})
	if err != nil {
		return nil, err
	}
	return &ctxs.UserCtx{
		IsOpen:     true,
		TenantCode: resp.TenantCode,
		UserID:     resp.UserID,
		IsAdmin:    resp.IsAdmin == def.True,
		IsAllData:  true,
		Account:    resp.UserName,
	}, nil
}
func (m *CheckTokenWareMiddleware) Auth(ctx context.Context, w http.ResponseWriter, strToken string, strIP string, authType string) (*ctxs.UserCtx, error) {
	resp, err := m.UserRpc.UserCheckToken(ctx, &user.UserCheckTokenReq{
		Ip:       strIP,
		Token:    strToken,
		AuthType: authType,
	})
	if err != nil {
		er := errors.Fmt(err)
		logx.WithContext(ctx).Errorf("%s.CheckTokenWare ip=%s token=%s return=%s",
			utils.FuncName(), strIP, strToken, err)
		return nil, er
	}

	if resp.Token != "" {
		w.Header().Set("Access-Control-Expose-Headers", ctxs.UserSetTokenKey)
		w.Header().Set(ctxs.UserSetTokenKey, resp.Token)
	}
	logx.WithContext(ctx).Debugf("%s.CheckTokenWare ip:%v in.token=%s  checkResp:%v",
		utils.FuncName(), strIP, strToken, utils.Fmt(resp))
	return &ctxs.UserCtx{
		IsOpen:       authType == "open",
		TenantCode:   resp.TenantCode,
		UserID:       resp.UserID,
		RoleIDs:      resp.RoleIDs,
		RoleCodes:    resp.RoleCodes,
		IsAdmin:      resp.IsAdmin || resp.IsSuperAdmin,
		IsSuperAdmin: resp.IsSuperAdmin,
		IsAllData:    resp.IsAllData == def.True,
		Account:      resp.Account,
		ProjectAuth:  utils.CopyMap[ctxs.ProjectAuth](resp.ProjectAuth),
	}, nil
}

func (m *CheckTokenWareMiddleware) UserAuth(w http.ResponseWriter, r *http.Request) (*ctxs.UserCtx, error) {
	strIP, _ := utils.GetIP(r)

	strToken := ctxs.GetHandle(r, ctxs.UserTokenKey, ctxs.UserToken2Key)
	if strToken == "" {
		logx.WithContext(r.Context()).Errorf("%s.CheckTokenWare ip=%s not find token",
			utils.FuncName(), strIP)
		return nil, errors.NotLogin
	}
	resp, err := m.UserRpc.UserCheckToken(r.Context(), &user.UserCheckTokenReq{
		Ip:    strIP,
		Token: strToken,
	})
	if err != nil {
		er := errors.Fmt(err)
		logx.WithContext(r.Context()).Errorf("%s.CheckTokenWare ip=%s token=%s return=%s",
			utils.FuncName(), strIP, strToken, err)
		return nil, er
	}

	if resp.Token != "" {
		w.Header().Set("Access-Control-Expose-Headers", ctxs.UserSetTokenKey)
		w.Header().Set(ctxs.UserSetTokenKey, resp.Token)
	}
	logx.WithContext(r.Context()).Debugf("%s.CheckTokenWare ip:%v in.token=%s  checkResp:%v",
		utils.FuncName(), strIP, strToken, utils.Fmt(resp))
	return &ctxs.UserCtx{
		IsOpen:       false,
		TenantCode:   resp.TenantCode,
		UserID:       resp.UserID,
		RoleIDs:      resp.RoleIDs,
		RoleCodes:    resp.RoleCodes,
		IsAdmin:      resp.IsAdmin || resp.IsSuperAdmin,
		IsSuperAdmin: resp.IsSuperAdmin,
		IsAllData:    resp.IsAllData == def.True,
		Account:      resp.Account,
		ProjectAuth:  utils.CopyMap[ctxs.ProjectAuth](resp.ProjectAuth),
	}, nil
}
