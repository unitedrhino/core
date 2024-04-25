package logic

import (
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
)

func ToPageInfo(info *sys.PageInfo, defaultOrders ...def.OrderBy) *def.PageInfo {
	if info == nil {
		return nil
	}

	var orders = defaultOrders
	if infoOrders := info.GetOrders(); len(infoOrders) > 0 {
		orders = make([]def.OrderBy, 0, len(infoOrders))
		for _, infoOd := range infoOrders {
			if infoOd.GetFiled() != "" {
				orders = append(orders, def.OrderBy{utils.CamelCaseToUdnderscore(infoOd.GetFiled()), infoOd.GetSort()})
			}
		}
	}

	return &def.PageInfo{
		Page:   info.GetPage(),
		Size:   info.GetSize(),
		Orders: orders,
	}
}

func ToPageInfoWithDefault(info *sys.PageInfo, defau *def.PageInfo) *def.PageInfo {
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
		Id:         ui.ID,
		Body:       utils.ToRpcNullString(ui.Body),
		ModuleCode: ui.ModuleCode,
		Name:       ui.Name,
		ParentID:   ui.ParentID,
		Type:       ui.Type,
		Path:       ui.Path,
		Component:  ui.Component,
		Icon:       ui.Icon,
		Redirect:   ui.Redirect,
		CreateTime: ui.CreatedTime.Unix(),
		Order:      ui.Order,
		HideInMenu: ui.HideInMenu,
	}
}
