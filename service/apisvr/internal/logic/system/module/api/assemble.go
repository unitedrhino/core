package api

import (
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
)

func ToApiInfoRpc(in *types.ApiInfo) *sys.ApiInfo {
	if in == nil {
		return nil
	}
	return &sys.ApiInfo{
		Id:           in.ID,
		BusinessType: in.BusinessType,
		ModuleCode:   in.ModuleCode,
		Route:        in.Route,
		Method:       in.Method,
		Group:        in.Group,
		Name:         in.Name,
		IsNeedAuth:   in.IsNeedAuth,
		Desc:         in.Desc,
	}
}

func ToApiInfoTypes(in *sys.ApiInfo) *types.ApiInfo {
	if in == nil {
		return nil
	}
	return &types.ApiInfo{
		ID:           in.Id,
		BusinessType: in.BusinessType,
		ModuleCode:   in.ModuleCode,
		Route:        in.Route,
		Method:       in.Method,
		Group:        in.Group,
		Name:         in.Name,
	}
}
