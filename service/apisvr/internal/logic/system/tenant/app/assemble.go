package app

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
)

func ToTenantAppModulePb(in *types.TenantAppModule) *sys.TenantAppModule {
	return &sys.TenantAppModule{
		Code:    in.Code,
		MenuIDs: in.MenuIDs,
	}
}
func ToTenantAppModulesPb(in []*types.TenantAppModule) (ret []*sys.TenantAppModule) {
	for _, v := range in {
		ret = append(ret, ToTenantAppModulePb(v))
	}
	return
}
