package middlewares

import (
	"bytes"
	"context"
	operLog "gitee.com/unitedrhino/core/service/syssvr/client/log"
	role "gitee.com/unitedrhino/core/service/syssvr/client/rolemanage"
	tenant "gitee.com/unitedrhino/core/service/syssvr/client/tenantmanage"
	user "gitee.com/unitedrhino/core/service/syssvr/client/usermanage"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/service/syssvr/sysdirect"
	"gitee.com/unitedrhino/core/share/domain/log"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"gitee.com/unitedrhino/share/tools"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"io"
	"net/http"
	"strings"
	"sync"
)

type CheckTokenWareMiddleware struct {
	UserRpc   user.UserManage
	AuthRpc   role.RoleManage
	TenantRpc tenant.TenantManage
	LogRpc    operLog.Log
}

var respPool sync.Pool
var bufferSize = 512

func init() {
	respPool.New = func() interface{} {
		return make([]byte, bufferSize)
	}
}

func NewCheckTokenWareMiddleware(UserRpc user.UserManage, AuthRpc role.RoleManage, TenantRpc tenant.TenantManage, LogRpc operLog.Log) *CheckTokenWareMiddleware {
	return &CheckTokenWareMiddleware{UserRpc: UserRpc, AuthRpc: AuthRpc, TenantRpc: TenantRpc, LogRpc: LogRpc}
}

func NewCheckTokenWareMiddleware2(SysRpc conf.RpcClientConf) *CheckTokenWareMiddleware {
	var (
		TenantRpc tenant.TenantManage
		LogRpc    operLog.Log
		UserRpc   user.UserManage
		AuthRpc   role.RoleManage
	)
	if SysRpc.Mode == conf.ClientModeDirect {
		TenantRpc = sysdirect.NewTenantManage(SysRpc.RunProxy)
		LogRpc = sysdirect.NewLog(SysRpc.RunProxy)
		UserRpc = sysdirect.NewUser(SysRpc.RunProxy)
		AuthRpc = sysdirect.NewRole(SysRpc.RunProxy)
	} else {
		TenantRpc = tenant.NewTenantManage(zrpc.MustNewClient(SysRpc.Conf))
		LogRpc = operLog.NewLog(zrpc.MustNewClient(SysRpc.Conf))
		UserRpc = user.NewUserManage(zrpc.MustNewClient(SysRpc.Conf))
		AuthRpc = role.NewRoleManage(zrpc.MustNewClient(SysRpc.Conf))
	}

	return &CheckTokenWareMiddleware{UserRpc: UserRpc, AuthRpc: AuthRpc, TenantRpc: TenantRpc, LogRpc: LogRpc}
}

func (m *CheckTokenWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			userCtx *ctxs.UserCtx
			err     error
			//isOpen   bool
			token    string
			strIP, _        = utils.GetIP(r)
			authType        = "user"
			appCode  string = ctxs.GetHandle(r, ctxs.UserAppCodeKey, ctxs.UserAppCodeKey2)
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
		if userCtx.AppCode != "" && userCtx.AppCode != appCode {
			result.HttpErr(w, r, http.StatusUnauthorized, errors.Permissions.AddMsg("认证失败,应用不一致"))
			return
		}
		userCtx.Os = ctxs.GetHandle(r, "User-Agent")
		userCtx.AcceptLanguage = ctxs.GetHandle(r, "Accept-Language")
		userCtx.Token = token
		strProjectID := ctxs.GetHandle(r, ctxs.UserProjectID, ctxs.UserProjectID2)
		projectID := cast.ToInt64(strProjectID)
		if projectID == 0 {
			projectID = def.NotClassified
		}
		if projectID > def.NotClassified && !userCtx.IsAdmin && userCtx.ProjectAuth[projectID] == nil {
			result.HttpErr(w, r, http.StatusOK, errors.Permissions.AddMsg("无所选项目的权限").AddDetailf(strProjectID))
			return
		}
		//注入 用户信息 到 ctx
		ctx2 := ctxs.SetUserCtx(r.Context(), userCtx)
		r = r.WithContext(ctx2)
		var apiRet *sys.RoleApiAuthResp
		////校验 Casbin Rule
		req := user.RoleApiAuthReq{
			Path:   r.URL.Path,
			Method: r.Method,
		}
		apiRet, err = m.AuthRpc.RoleApiAuth(r.Context(), &req)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("%s.AuthApiCheck error=%s", utils.FuncName(), err)
			http.Error(w, "接口权限不足："+err.Error(), http.StatusUnauthorized)
			//systems.SysNotify(fmt.Sprintf("接口权限不足userCtx:%v req:%v err:%s", utils.Fmt(userCtx), utils.Fmt(req), err))
			return
		}

		m.OperationLogRecord(next, w, r, apiRet)
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
		AppCode:      resp.AppCode,
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
		AppCode:      resp.AppCode,
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

// 接口操作日志记录
func (m *CheckTokenWareMiddleware) OperationLogRecord(next http.HandlerFunc, w http.ResponseWriter, r *http.Request, apiInfo *sys.RoleApiAuthResp) {
	ctx := ctxs.CopyCtx(r.Context())
	useCtx := ctxs.GetUserCtx(ctx)
	if useCtx.IsOpen || useCtx.UserID == 0 || apiInfo == nil || apiInfo.RecordLogMode == 3 || (apiInfo.RecordLogMode == 1 && apiInfo.BusinessType == log.OptQuery) || apiInfo.BusinessType == 0 {
		next(w, r)
		return
	}
	var reqBodyStr string
	if r.Body != nil {
		reqBody, _ := io.ReadAll(r.Body)                //读取 reqBody
		r.Body = io.NopCloser(bytes.NewReader(reqBody)) //重建 reqBody
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			if len(reqBody) > bufferSize {
				// 截断
				newBody := respPool.Get().([]byte)
				copy(newBody, reqBody)
				defer respPool.Put(newBody)
			}
		}
		var reqLen = len(reqBody)
		if reqLen > bufferSize {
			reqLen = bufferSize
		}
		reqBodyStr = string(reqBody[:reqLen])
	}

	respStatusCode := http.StatusOK
	respStatusMsg := ""
	respBodyStr := ""
	r = ctxs.NeedResp(r)
	next(w, r)
	resp := ctxs.GetResp(r)
	if resp != nil {
		respStatusCode = resp.StatusCode
		respStatusMsg = resp.Status
		if resp.Body != nil {
			respBody, _ := io.ReadAll(resp.Body) //读取 respBody
			var respLen = len(respBody)
			if respLen > bufferSize {
				respLen = bufferSize
			}
			respBodyStr = string(respBody[:respLen])
		}

	}

	uri := "https://"
	if !strings.Contains(r.Proto, "HTTPS") {
		uri = "http://"
	}

	ipAddr, err := utils.GetIP(r)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.GetIP is error : %s req:%v",
			utils.FuncName(), err.Error(), utils.Fmt(r))
		ipAddr = "0.0.0.0"
	}
	utils.Go(ctx, func() {
		_, err = m.LogRpc.OperLogCreate(ctx, &user.OperLogCreateReq{
			Uri:          uri + r.Host + r.RequestURI,
			Route:        r.RequestURI,
			OperName:     apiInfo.Name,
			BusinessType: apiInfo.BusinessType,
			OperIpAddr:   ipAddr,
			OperLocation: tools.GetCityByIp(ipAddr),
			Code:         int64(respStatusCode),
			Msg:          respStatusMsg,
			Req:          reqBodyStr,
			Resp:         respBodyStr,
			AppCode:      ctxs.GetUserCtx(r.Context()).AppCode,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.OperationLogRecord is error : %s",
				utils.FuncName(), err.Error())
		}
		return
	})

}
