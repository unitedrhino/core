// 前端 OAuth 字段别名兼容（googleConfig/appleConfig 等）
package tenant

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

// NormalizeTenantAppInfoIn 将前端字段名映射为后端标准字段（写入前）
func NormalizeTenantAppInfoIn(req *types.TenantAppInfo) {
	if req == nil {
		return
	}
	if req.Google == nil && req.GoogleConfig != nil {
		req.Google = req.GoogleConfig
	}
	if req.Github == nil && req.GithubConfig != nil {
		req.Github = req.GithubConfig
	}
	if req.Apple == nil && req.AppleConfig != nil {
		req.Apple = req.AppleConfig
	}
	normalizeAppleIn(req.Apple)
}

// NormalizeTenantConfigIn 租户配置页写入前字段归一
func NormalizeTenantConfigIn(req *types.TenantConfig) {
	if req == nil {
		return
	}
	if req.Google == nil && req.GoogleConfig != nil {
		req.Google = req.GoogleConfig
	}
	if req.Github == nil && req.GithubConfig != nil {
		req.Github = req.GithubConfig
	}
	if req.Apple == nil && req.AppleConfig != nil {
		req.Apple = req.AppleConfig
	}
	normalizeAppleIn(req.Apple)
}

// FillTenantAppInfoOut 响应时同步前端使用的别名字段（读出后）
func FillTenantAppInfoOut(req *types.TenantAppInfo) {
	if req == nil {
		return
	}
	req.GoogleConfig = req.Google
	req.GithubConfig = req.Github
	req.AppleConfig = fillAppleOut(req.Apple)
}

// FillTenantAppOut 列表项响应别名
func FillTenantAppOut(req *types.TenantApp) {
	if req == nil {
		return
	}
	req.GoogleConfig = req.Google
	req.GithubConfig = req.Github
	req.AppleConfig = fillAppleOut(req.Apple)
}

// FillTenantConfigOut 租户配置读响应别名
func FillTenantConfigOut(req *types.TenantConfig) {
	if req == nil {
		return
	}
	req.GoogleConfig = req.Google
	req.GithubConfig = req.Github
	req.AppleConfig = fillAppleOut(req.Apple)
}

// normalizeAppleIn bundleID 兼容写入 appID
func normalizeAppleIn(a *types.AppleAppConfig) {
	if a == nil {
		return
	}
	if a.AppID == "" && a.BundleID != "" {
		a.AppID = a.BundleID
	}
}

// fillAppleOut 读出时同步 bundleID 供前端回显
func fillAppleOut(a *types.AppleAppConfig) *types.AppleAppConfig {
	if a == nil {
		return nil
	}
	out := *a
	if out.BundleID == "" && out.AppID != "" {
		out.BundleID = out.AppID
	}
	return &out
}

// ToSysTenantAppInfo 归一化后转 RPC 请求
func ToSysTenantAppInfo(req *types.TenantAppInfo) *sys.TenantAppInfo {
	NormalizeTenantAppInfoIn(req)
	return utils.Copy[sys.TenantAppInfo](req)
}
