package logic

import (
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
)

func ToPageInfo(info *sys.PageInfo) *stores.PageInfo {
	return utils.Copy[stores.PageInfo](info)
}

func ToPageInfoWithDefault(info *sys.PageInfo, defau *stores.PageInfo) *stores.PageInfo {
	if page := ToPageInfo(info); page == nil {
		return defau
	} else {
		if page.Page == 0 {
			page.Page = defau.Page
		}
		if page.Size == 0 {
			page.Size = defau.Size
		}
		if len(page.Orders) == 0 {
			page.Orders = defau.Orders
		}
		return page
	}
}

func ToSysPoint(point stores.Point) *sys.Point {
	return &sys.Point{Longitude: point.Longitude, Latitude: point.Latitude}
}
func ToStorePoint(point *sys.Point) stores.Point {
	if point == nil {
		return stores.Point{Longitude: 0, Latitude: 0}
	}
	return stores.Point{Longitude: point.Longitude, Latitude: point.Latitude}
}

func MenuInfoToPb(ui *relationDB.SysModuleMenu) *sys.MenuInfo {
	return &sys.MenuInfo{
		Id:          ui.ID,
		Body:        utils.ToRpcNullString(ui.Body),
		ModuleCode:  ui.ModuleCode,
		Name:        ui.Name,
		ParentID:    ui.ParentID,
		Type:        ui.Type,
		Path:        ui.Path,
		Component:   ui.Component,
		Icon:        ui.Icon,
		Redirect:    ui.Redirect,
		CreatedTime: ui.CreatedTime.Unix(),
		Order:       ui.Order,
		HideInMenu:  ui.HideInMenu,
	}
}
