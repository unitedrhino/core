package user

import (
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/role"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"
	"github.com/golang/protobuf/ptypes/wrappers"
)

func GetNullVal(val *wrappers.StringValue) *string {
	if val == nil {
		return nil
	}
	return &val.Value
}

func UserInfoToApi(ui *sys.UserInfo, roles []*sys.RoleInfo, tenant *sys.TenantInfo) *types.UserInfo {
	if ui == nil {
		return nil
	}
	ret := utils.Copy[types.UserInfo](ui)
	ret.Roles = role.ToRoleInfosTypes(roles)
	ret.Tenant = system.ToTenantInfoTypes(tenant, nil, nil)
	return ret
}
func UserInfoToRpc(ui *types.UserInfo) *sys.UserInfo {
	return utils.Copy[sys.UserInfo](ui)
}
