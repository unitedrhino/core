package export

import (
	role "gitee.com/i-Things/core/service/syssvr/client/rolemanage"
	user "gitee.com/i-Things/core/service/syssvr/client/usermanage"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type CheckTokenWareMiddleware struct {
	UserRpc user.UserManage
	AuthRpc role.RoleManage
}

func NewCheckTokenWareMiddleware(UserRpc user.UserManage, AuthRpc role.RoleManage) *CheckTokenWareMiddleware {
	return &CheckTokenWareMiddleware{UserRpc: UserRpc, AuthRpc: AuthRpc}
}

func (m *CheckTokenWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.WithContext(r.Context()).Infof("%s.Lifecycle.Before", utils.FuncName())

		var (
			userCtx *ctxs.UserCtx
			err     error
		)

		//如果是用户请求
		//校验 Jwt Token
		userCtx, err = m.UserAuth(w, r)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("%s.UserAuth error=%s", utils.FuncName(), err)
			http.Error(w, "用户请求失败："+err.Error(), http.StatusUnauthorized)
			return
		}
		//注入 用户信息 到 ctx
		ctx2 := ctxs.SetUserCtx(r.Context(), userCtx)
		r = r.WithContext(ctx2)
		////校验 Casbin Rule
		_, err = m.AuthRpc.RoleApiAuth(r.Context(), &user.RoleApiAuthReq{
			RoleID: userCtx.RoleID,
			Path:   r.URL.Path,
			Method: r.Method,
		})
		if err != nil {
			logx.WithContext(r.Context()).Errorf("%s.AuthApiCheck error=%s", utils.FuncName(), err)
			//http.Error(w, "接口权限不足："+err.Error(), http.StatusUnauthorized)
			//return
		}
		next(w, r)
		logx.WithContext(r.Context()).Infof("%s.Lifecycle.After", utils.FuncName())
	}
}

func (m *CheckTokenWareMiddleware) UserAuth(w http.ResponseWriter, r *http.Request) (*ctxs.UserCtx, error) {
	strIP, _ := utils.GetIP(r)

	strToken := r.Header.Get(ctxs.UserTokenKey)
	if strToken == "" {
		logx.WithContext(r.Context()).Errorf("%s.CheckTokenWare ip=%s not find token",
			utils.FuncName(), strIP)
		return nil, errors.NotLogin
	}
	strProjectID := r.Header.Get(ctxs.UserProjectID)
	projectID := cast.ToInt64(strProjectID)
	if projectID == 0 {
		projectID = def.NotClassified
	}
	strRoleID := r.Header.Get(ctxs.UserRoleKey)
	roleID := cast.ToInt64(strRoleID)

	appCode := r.Header.Get(ctxs.UserAppCodeKey)

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
	logx.WithContext(r.Context()).Infof("%s.CheckTokenWare ip:%v in.token=%s roleID：%v checkResp:%v",
		utils.FuncName(), strIP, strToken, strRoleID, utils.Fmt(resp))
	if roleID != 0 { //如果传了角色
		if !utils.SliceIn(roleID, resp.RoleIDs...) {
			err := errors.Parameter.AddMsgf("所选角色无权限")
			return nil, err
		}
	} else {
		roleID = resp.RoleIDs[0]
	}
	return &ctxs.UserCtx{
		IsOpen:     false,
		TenantCode: resp.TenantCode,
		ProjectID:  projectID,
		AppCode:    appCode,
		UserID:     resp.UserID,
		RoleID:     roleID,
		IsAdmin:    resp.IsAdmin == def.True,
		IsAllData:  resp.IsAllData == def.True,
	}, nil
}
