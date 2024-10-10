package menu

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

func ToMenuInfoRpc(in *types.MenuInfo) *sys.MenuInfo {
	if in == nil {
		return nil
	}
	return &sys.MenuInfo{
		Id:         in.ID,
		Name:       in.Name,
		ParentID:   in.ParentID,
		Type:       in.Type,
		Path:       in.Path,
		Component:  in.Component,
		Icon:       in.Icon,
		Redirect:   in.Redirect,
		Order:      in.Order,
		HideInMenu: in.HideInMenu,
		Body:       utils.ToRpcNullString(in.Body),
		ModuleCode: in.ModuleCode,
		IsCommon:   in.IsCommon,
	}
}
func ToMenuInfosRpc(in []*types.MenuInfo) (ret []*sys.MenuInfo) {
	for _, v := range in {
		ret = append(ret, ToMenuInfoRpc(v))
	}
	return
}
