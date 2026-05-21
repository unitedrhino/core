// 租户配置与应用登录配置的组装逻辑
package config

import (
	tenantoauth "gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/tenant"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

// mergeAppLoginIntoConfig 将租户应用登录配置合并到租户配置响应
func mergeAppLoginIntoConfig(dst *types.TenantConfig, app *sys.TenantAppInfo) {
	if dst == nil || app == nil {
		return
	}
	dst.AppCode = app.AppCode
	if app.WxMini != nil && app.WxMini.AppID != "" {
		dst.WxMini = utils.Copy[types.ThirdAppConfig](app.WxMini)
	}
	if app.DingMini != nil && app.DingMini.AppID != "" {
		dst.DingMini = utils.Copy[types.ThirdAppConfig](app.DingMini)
	}
	if app.WxOpen != nil && app.WxOpen.AppID != "" {
		dst.WxOpen = utils.Copy[types.ThirdAppConfig](app.WxOpen)
	}
	if app.Huawei != nil && (app.Huawei.AppID != "" || app.Huawei.AppSecret != "") {
		dst.Huawei = utils.Copy[types.ThirdAppConfig](app.Huawei)
	}
	if app.Google != nil && (app.Google.AppID != "" || app.Google.AppSecret != "" || app.Google.AppKey != "") {
		dst.Google = utils.Copy[types.ThirdAppConfig](app.Google)
	}
	if app.Github != nil && (app.Github.AppID != "" || app.Github.AppSecret != "" || app.Github.AppKey != "") {
		dst.Github = utils.Copy[types.ThirdAppConfig](app.Github)
	}
	if app.Apple != nil && (app.Apple.AppID != "" || app.Apple.PrivateKey != "") {
		dst.Apple = utils.Copy[types.AppleAppConfig](app.Apple)
	}
	dst.LoginTypes = app.LoginTypes
	dst.IsAutoRegister = app.IsAutoRegister
	tenantoauth.FillTenantConfigOut(dst)
}

// hasAppLoginPayload 判断是否携带应用登录配置字段
func hasAppLoginPayload(req *types.TenantConfig) bool {
	if req == nil {
		return false
	}
	tenantoauth.NormalizeTenantConfigIn(req)
	return req.AppCode != "" || req.DingMini != nil || req.WxOpen != nil || req.WxMini != nil ||
		req.Huawei != nil || req.Google != nil || req.Github != nil || req.Apple != nil ||
		req.LoginTypes != nil || req.IsAutoRegister != 0
}

// toTenantAppInfo 从租户配置请求构建租户应用更新参数
func toTenantAppInfo(req *types.TenantConfig, tenantCode, appCode string) *sys.TenantAppInfo {
	tenantoauth.NormalizeTenantConfigIn(req)
	return &sys.TenantAppInfo{
		Code:           tenantCode,
		AppCode:        appCode,
		DingMini:       utils.Copy[sys.ThirdAppConfig](req.DingMini),
		WxOpen:         utils.Copy[sys.ThirdAppConfig](req.WxOpen),
		WxMini:         utils.Copy[sys.ThirdAppConfig](req.WxMini),
		Huawei:         utils.Copy[sys.ThirdAppConfig](req.Huawei),
		Google:         utils.Copy[sys.ThirdAppConfig](req.Google),
		Github:         utils.Copy[sys.ThirdAppConfig](req.Github),
		Apple:          utils.Copy[sys.AppleAppConfig](req.Apple),
		LoginTypes:     req.LoginTypes,
		IsAutoRegister: req.IsAutoRegister,
	}
}
