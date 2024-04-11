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
	ret.Tenant = system.ToTenantInfoTypes(tenant)
	return ret
}
func UserInfoToRpc(ui *types.UserInfo) *sys.UserInfo {
	if ui == nil {
		return nil
	}
	return &sys.UserInfo{
		UserID:          ui.UserID,
		UserName:        ui.UserName,
		Email:           utils.ToRpcNullString(ui.Email),
		Phone:           utils.ToRpcNullString(ui.Phone),
		LastIP:          ui.LastIP,
		RegIP:           ui.RegIP,
		Role:            ui.Role,
		NickName:        ui.NickName,
		Sex:             ui.Sex,
		IsAllData:       ui.IsAllData,
		City:            ui.City,
		Country:         ui.Country,
		Province:        ui.Province,
		Language:        ui.Language,
		HeadImg:         ui.HeadImg,
		IsUpdateHeadImg: ui.IsUpdateHeadImg,
		Password:        ui.Password,
		CreatedTime:     ui.CreatedTime,
	}
}
