package system

import (
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/area/info"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"
)

func ProjectInfoToApi(pb *sys.ProjectInfo, ui *sys.UserInfo) *types.ProjectInfo {
	return &types.ProjectInfo{
		CreatedTime:   pb.CreatedTime,
		ProjectID:     pb.ProjectID,
		ProjectName:   pb.ProjectName,
		AdminUserID:   pb.AdminUserID,
		ProjectImg:    pb.ProjectImg,
		Desc:          utils.ToNullString(pb.Desc),
		Position:      logic.ToSysPointApi(pb.Position),
		IsSysCreated:  pb.IsSysCreated,
		AreaCount:     pb.AreaCount,
		AdminUserInfo: utils.Copy[types.UserCore](ui),
		Area:          utils.ToNullFloat32(pb.Area),
		Areas:         info.ToAreaInfosTypes(pb.Areas),
	}
}
func ProjectInfosToApi(pb []*sys.ProjectInfo) (ret []*types.ProjectInfo) {
	for _, v := range pb {
		ret = append(ret, ProjectInfoToApi(v, nil))
	}
	return
}

func ToMenuInfoApi(i *sys.MenuInfo) *types.MenuInfo {
	return &types.MenuInfo{
		ModuleCode: i.ModuleCode,
		ID:         i.Id,
		Name:       i.Name,
		ParentID:   i.ParentID,
		Type:       i.Type,
		Path:       i.Path,
		Component:  i.Component,
		Icon:       i.Icon,
		Redirect:   i.Redirect,
		CreateTime: i.CreateTime,
		Order:      i.Order,
		HideInMenu: i.HideInMenu,
		Body:       utils.ToNullString(i.Body),
		Children:   ToMenuInfosApi(i.Children),
	}
}
func ToMenuInfosApi(i []*sys.MenuInfo) (ret []*types.MenuInfo) {
	if i == nil {
		return nil
	}
	for _, v := range i {
		ret = append(ret, ToMenuInfoApi(v))
	}
	return
}

func ToTenantAppMenuApi(i *sys.TenantAppMenu) *types.TenantAppMenu {
	if i == nil {
		return nil
	}
	return &types.TenantAppMenu{
		TemplateID: i.TemplateID,
		Code:       i.Code,
		AppCode:    i.AppCode,
		MenuInfo:   *ToMenuInfoApi(i.Info),
		Children:   ToTenantAppMenusApi(i.Children),
	}
}
func ToTenantAppMenusApi(i []*sys.TenantAppMenu) (ret []*types.TenantAppMenu) {
	for _, v := range i {
		ret = append(ret, ToTenantAppMenuApi(v))
	}
	return
}

func ToSysWithIDCode(in *types.WithIDOrCode) *sys.WithIDCode {
	return &sys.WithIDCode{
		Id:   in.ID,
		Code: in.Code,
	}
}

func ToTenantInfoRpc(in *types.TenantInfo) *sys.TenantInfo {
	return utils.Copy[sys.TenantInfo](in)
}

func ToTenantInfoTypes(in *sys.TenantInfo, user *sys.UserInfo) *types.TenantInfo {
	if in == nil {
		return nil
	}
	ret := utils.Copy[types.TenantInfo](in)
	ret.AdminUserInfo = utils.Copy[types.UserCore](user)
	return ret
}

func ToTenantCoreTypes(in *sys.TenantInfo) *types.TenantCore {
	return utils.Copy[types.TenantCore](in)
}

func ToTenantInfosTypes(in []*sys.TenantInfo, userMap map[int64]*sys.UserInfo) []*types.TenantInfo {
	var ret []*types.TenantInfo
	for _, v := range in {
		ti := ToTenantInfoTypes(v, userMap[v.AdminUserID])
		ret = append(ret, ti)
	}
	return ret
}
