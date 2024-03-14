package logic

import (
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/domain/tenant"
	"gitee.com/i-Things/share/utils"
)

func ToTenantInfoRpc(in *relationDB.SysTenantInfo) *sys.TenantInfo {
	if in == nil {
		return nil
	}
	return &sys.TenantInfo{
		Id:               in.ID,
		Code:             in.Code,
		Name:             in.Name,
		AdminUserID:      in.AdminUserID,
		AdminRoleID:      in.AdminRoleID,
		BaseUrl:          in.BaseUrl,
		LogoUrl:          in.LogoUrl,
		Desc:             utils.ToRpcNullString(in.Desc),
		DefaultProjectID: in.DefaultProjectID,
		ProjectLimit:     in.ProjectLimit,
	}
}
func ToTenantInfosRpc(in []*relationDB.SysTenantInfo) (ret []*sys.TenantInfo) {
	for _, v := range in {
		ret = append(ret, ToTenantInfoRpc(v))
	}
	return
}

func ToTenantInfoPo(in *sys.TenantInfo) *relationDB.SysTenantInfo {
	if in == nil {
		return nil
	}
	return &relationDB.SysTenantInfo{
		ID:               in.Id,
		Code:             in.Code,
		Name:             in.Name,
		AdminUserID:      in.AdminUserID,
		AdminRoleID:      in.AdminRoleID,
		BaseUrl:          in.BaseUrl,
		LogoUrl:          in.LogoUrl,
		Desc:             utils.ToEmptyString(in.Desc),
		DefaultProjectID: in.DefaultProjectID,
		ProjectLimit:     in.ProjectLimit,
	}
}

func ToTenantInfoCaches(in []*relationDB.SysTenantInfo) (ret []*tenant.Info) {
	for _, v := range in {
		ret = append(ret, ToTenantInfoCache(v))
	}
	return ret
}

func ToTenantInfoCache(in *relationDB.SysTenantInfo) *tenant.Info {
	if in == nil {
		return nil
	}
	return &tenant.Info{
		ID:               in.ID,
		Code:             in.Code,
		Name:             in.Name,
		AdminUserID:      in.AdminUserID,
		AdminRoleID:      in.AdminRoleID,
		BaseUrl:          in.BaseUrl,
		LogoUrl:          in.LogoUrl,
		Desc:             in.Desc,
		CreatedTime:      in.CreatedTime.Unix(),
		DefaultProjectID: in.DefaultProjectID,
		ProjectLimit:     in.ProjectLimit,
	}
}
func CacheToTenantInfoRpc(in *tenant.Info) *sys.TenantInfo {
	if in == nil {
		return nil
	}
	return &sys.TenantInfo{
		Id:               in.ID,
		Code:             in.Code,
		Name:             in.Name,
		AdminUserID:      in.AdminUserID,
		AdminRoleID:      in.AdminRoleID,
		BaseUrl:          in.BaseUrl,
		LogoUrl:          in.LogoUrl,
		DefaultProjectID: in.DefaultProjectID,
		ProjectLimit:     in.ProjectLimit,
		Desc:             utils.ToRpcNullString(in.Desc),
	}
}

func RpcToTenantInfoCache(in *sys.TenantInfo) *tenant.Info {
	if in == nil {
		return nil
	}
	return &tenant.Info{
		ID:               in.Id,
		Code:             in.Code,
		Name:             in.Name,
		AdminUserID:      in.AdminUserID,
		AdminRoleID:      in.AdminRoleID,
		BaseUrl:          in.BaseUrl,
		LogoUrl:          in.LogoUrl,
		DefaultProjectID: in.DefaultProjectID,
		ProjectLimit:     in.ProjectLimit,
		Desc:             utils.ToEmptyString(in.Desc),
	}
}
