package logic

import (
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
)

func ToModuleInfoPo(in *sys.ModuleInfo) *relationDB.SysModuleInfo {
	return utils.Copy[relationDB.SysModuleInfo](in)
}

func ToModuleInfoPb(in *relationDB.SysModuleInfo) *sys.ModuleInfo {
	return utils.Copy[sys.ModuleInfo](in)
}

func ToModuleInfosPb(in []*relationDB.SysModuleInfo) (ret []*sys.ModuleInfo) {
	for _, v := range in {
		ret = append(ret, ToModuleInfoPb(v))
	}
	return
}

//func ToTenantApiInfoPo(in *sys.TenantApiInfo) *relationDB.SysTenantAppApi {
//	if in == nil || in.Info == nil {
//		return nil
//	}
//	return &relationDB.SysTenantAppApi{
//		TempLateID:   in.TemplateID,
//		TenantCode:   stores.TenantCode(in.TemplateCode),
//		AppCode:      in.AppCode,
//		SysModuleApi: *ToApiInfoPo(in.Info),
//	}
//}

func ToMenuInfoPo(in *sys.MenuInfo) *relationDB.SysModuleMenu {
	if in == nil {
		return nil
	}
	return &relationDB.SysModuleMenu{
		ID:         in.Id,
		ModuleCode: in.ModuleCode,
		ParentID:   in.ParentID,
		Type:       in.Type,
		Order:      in.Order,
		Name:       in.Name,
		Path:       in.Path,
		Component:  in.Component,
		Icon:       in.Icon,
		Redirect:   in.Redirect,
		Body:       in.Body.GetValue(),
		HideInMenu: in.HideInMenu,
	}
}

func ToTenantAppMenuPo(in *sys.TenantAppMenu) *relationDB.SysTenantAppMenu {
	if in == nil || in.Info == nil {
		return nil
	}
	return &relationDB.SysTenantAppMenu{
		TempLateID:    in.TemplateID,
		TenantCode:    stores.TenantCode(in.Code),
		AppCode:       in.AppCode,
		SysModuleMenu: *ToMenuInfoPo(in.Info),
	}
}

func ToMenuInfoPb(in *relationDB.SysModuleMenu) *sys.MenuInfo {
	if in == nil {
		return nil
	}
	return &sys.MenuInfo{
		Id:         in.ID,
		ModuleCode: in.ModuleCode,
		ParentID:   in.ParentID,
		Type:       in.Type,
		Order:      in.Order,
		Name:       in.Name,
		Path:       in.Path,
		Component:  in.Component,
		Icon:       in.Icon,
		IsCommon:   in.IsCommon,
		Redirect:   in.Redirect,
		Body:       utils.ToRpcNullString(in.Body),
		HideInMenu: in.HideInMenu,
		CreateTime: in.CreatedTime.Unix(),
	}
}

func ToTenantAppMenuInfoPb(in *relationDB.SysTenantAppMenu) *sys.TenantAppMenu {
	if in == nil {
		return nil
	}
	return &sys.TenantAppMenu{
		TemplateID: in.TempLateID,
		Code:       string(in.TenantCode),
		AppCode:    in.AppCode,
		Info:       ToMenuInfoPb(&in.SysModuleMenu),
	}
}
