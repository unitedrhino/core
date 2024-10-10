package info

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

func ToModuleInfoRpc(in *types.ModuleInfo) *sys.ModuleInfo {
	return &sys.ModuleInfo{
		Id:         in.ID,
		Code:       in.Code,
		Name:       in.Name,
		Type:       in.Type,
		SubType:    in.SubType,
		Path:       in.Path,
		Desc:       utils.ToRpcNullString(in.Desc),
		Icon:       in.Icon,
		Url:        in.Url,
		Order:      in.Order,
		HideInMenu: in.HideInMenu,
		Body:       utils.ToRpcNullString(in.Body),
	}
}
func ToModuleInfoApi(in *sys.ModuleInfo) *types.ModuleInfo {
	if in == nil {
		return nil
	}
	return &types.ModuleInfo{
		ID:         in.Id,
		Code:       in.Code,
		Name:       in.Name,
		Type:       in.Type,
		SubType:    in.SubType,
		Path:       in.Path,
		Desc:       utils.ToNullString(in.Desc),
		Icon:       in.Icon,
		Url:        in.Url,
		Order:      in.Order,
		HideInMenu: in.HideInMenu,
		Body:       utils.ToNullString(in.Body),
	}

}
func ToModuleInfosApi(in []*sys.ModuleInfo) (ret []*types.ModuleInfo) {
	for _, v := range in {
		v1 := ToModuleInfoApi(v)
		if v1 != nil {
			ret = append(ret, v1)
		}
	}
	return
}
