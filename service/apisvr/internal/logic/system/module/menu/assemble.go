package menu

import (
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/core/shared/utils"
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
	}
}
func ToMenuInfosRpc(in []*types.MenuInfo) (ret []*sys.MenuInfo) {
	for _, v := range in {
		ret = append(ret, ToMenuInfoRpc(v))
	}
	return
}
