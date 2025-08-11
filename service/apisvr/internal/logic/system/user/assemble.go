package user

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/role"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
	"github.com/golang/protobuf/ptypes/wrappers"
)

func GetNullVal(val *wrappers.StringValue) *string {
	if val == nil {
		return nil
	}
	return &val.Value
}

type UserOpt struct {
	Roles  []*sys.RoleInfo
	Tenant *sys.TenantInfo
	Depts  []*sys.DeptInfo
}

func UserInfoToApi(ui *sys.UserInfo, opt UserOpt) *types.UserInfo {
	if ui == nil {
		return nil
	}
	ret := utils.Copy[types.UserInfo](ui)
	ret.Password = "xxxx"
	ret.Roles = role.ToRoleInfosTypes(opt.Roles)
	ret.Tenant = system.ToTenantInfoTypes(opt.Tenant, nil, nil)
	ret.Depts = utils.CopySlice[types.DeptInfo](opt.Depts)
	return ret
}
func UserInfoToRpc(ui *types.UserInfo) *sys.UserInfo {
	return utils.Copy[sys.UserInfo](ui)
}
