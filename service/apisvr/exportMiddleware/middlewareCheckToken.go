package exportMiddleware

import (
	"bytes"
	"context"
	"fmt"
	operLog "gitee.com/i-Things/core/service/syssvr/client/log"
	role "gitee.com/i-Things/core/service/syssvr/client/rolemanage"
	tenant "gitee.com/i-Things/core/service/syssvr/client/tenantmanage"
	user "gitee.com/i-Things/core/service/syssvr/client/usermanage"
	"gitee.com/i-Things/core/service/syssvr/domain/log"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
	"gitee.com/i-Things/share/utils"
	"github.com/gogf/gf/v2/encoding/gcharset"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/zeromicro/go-zero/core/logx"
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
var bufferSize = 1024

func init() {
	respPool.New = func() interface{} {
		return make([]byte, bufferSize)
	}
}

func NewCheckTokenWareMiddleware(UserRpc user.UserManage, AuthRpc role.RoleManage, TenantRpc tenant.TenantManage, LogRpc operLog.Log) *CheckTokenWareMiddleware {
	return &CheckTokenWareMiddleware{UserRpc: UserRpc, AuthRpc: AuthRpc, TenantRpc: TenantRpc, LogRpc: LogRpc}
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
			req := user.RoleApiAuthReq{
				Path:   r.URL.Path,
				Method: r.Method,
			}
			ret, err := m.AuthRpc.RoleApiAuth(r.Context(), &req)
			if err != nil {
				logx.WithContext(r.Context()).Errorf("%s.AuthApiCheck error=%s", utils.FuncName(), err)
				//http.Error(w, "接口权限不足："+err.Error(), http.StatusUnauthorized)
				clients.SysNotify(fmt.Sprintf("接口权限不足userCtx:%v req:%v err:%s", utils.Fmt(userCtx), utils.Fmt(req), err))
				//return
			} else if ret.BusinessType != log.OptQuery {
				m.OperationLogRecord(r.Context(), r, ret)
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

// 接口操作日志记录
func (m *CheckTokenWareMiddleware) OperationLogRecord(ctx context.Context, r *http.Request, apiInfo *sys.RoleApiAuthResp) {
	ctx = ctxs.CopyCtx(ctx)
	useCtx := ctxs.GetUserCtx(ctx)
	if useCtx.IsOpen || useCtx.UserID == 0 {
		return
	}
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
	reqBodyStr := string(reqBody)

	respStatusCode := http.StatusOK
	respStatusMsg := ""
	respBodyStr := ""

	if r.Response != nil {
		respStatusCode = r.Response.StatusCode
		respStatusMsg = r.Response.Status
		respBody, _ := io.ReadAll(r.Response.Body)                //读取 respBody
		r.Response.Body = io.NopCloser(bytes.NewReader(respBody)) //重建 respBody
		respBodyStr = string(respBody)
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
			OperLocation: m.GetCityByIp(ipAddr),
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

// 获取ip所属城市
func (m *CheckTokenWareMiddleware) GetCityByIp(ip string) string {
	if ip == "" {
		return ""
	}
	if ip == "[::1]" || ip == "127.0.0.1" {
		return "内网IP"
	}

	url := "http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip
	bytes := g.Client().GetBytes(context.TODO(), url)
	src := string(bytes)
	srcCharset := "GBK"

	tmp, _ := gcharset.ToUTF8(srcCharset, src)
	json, err := gjson.DecodeToJson(tmp)
	if err != nil {
		return ""
	}

	if json.Get("code").Int() == 0 {
		city := fmt.Sprintf("%s %s", json.Get("pro").String(), json.Get("city").String())
		return city
	} else {
		return ""
	}
}
