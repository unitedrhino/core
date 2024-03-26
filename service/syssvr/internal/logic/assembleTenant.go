package logic

import (
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/domain/tenant"
	"gitee.com/i-Things/share/utils"
)

func ToTenantInfoRpc(in *relationDB.SysTenantInfo) *sys.TenantInfo {
	return utils.Copy[sys.TenantInfo](in)
}
func ToTenantInfosRpc(in []*relationDB.SysTenantInfo) (ret []*sys.TenantInfo) {
	for _, v := range in {
		ret = append(ret, ToTenantInfoRpc(v))
	}
	return
}

func ToTenantInfoPo(in *sys.TenantInfo) *relationDB.SysTenantInfo {
	return utils.Copy[relationDB.SysTenantInfo](in)
}

func ToTenantInfoCaches(in []*relationDB.SysTenantInfo) (ret []*tenant.Info) {
	for _, v := range in {
		ret = append(ret, ToTenantInfoCache(v))
	}
	return ret
}

func ToTenantInfoCache(in *relationDB.SysTenantInfo) *tenant.Info {
	return utils.Copy[tenant.Info](in)
}
func CacheToTenantInfoRpc(in *tenant.Info) *sys.TenantInfo {
	return utils.Copy[sys.TenantInfo](in)
}

func RpcToTenantInfoCache(in *sys.TenantInfo) *tenant.Info {
	return utils.Copy[tenant.Info](in)
}
