package system

import (
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"
)

func ProjectInfoToApi(pb *sys.ProjectInfo) *types.ProjectInfo {
	return &types.ProjectInfo{
		CreatedTime: pb.CreatedTime,
		ProjectID:   pb.ProjectID,
		ProjectName: pb.ProjectName,
		AdminUserID: pb.AdminUserID,
		Desc:        utils.ToNullString(pb.Desc),
		Position:    logic.ToSysPointApi(pb.Position),
		AreaCount:   pb.AreaCount,
	}
}
func ProjectInfosToApi(pb []*sys.ProjectInfo) (ret []*types.ProjectInfo) {
	for _, v := range pb {
		ret = append(ret, ProjectInfoToApi(v))
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

func ToTenantInfoTypes(in *sys.TenantInfo) *types.TenantInfo {
	return utils.Copy[types.TenantInfo](in)
}

func ToTenantCoreTypes(in *sys.TenantInfo) *types.TenantCore {
	return utils.Copy[types.TenantCore](in)
}

func ToTenantInfosTypes(in []*sys.TenantInfo, userMap map[int64]*sys.UserInfo) []*types.TenantInfo {
	var ret []*types.TenantInfo
	for _, v := range in {
		ti := ToTenantInfoTypes(v)
		ti.AdminUserInfo = utils.Copy[types.UserCore](userMap[v.AdminUserID])
		ret = append(ret, ti)
	}
	return ret
}
