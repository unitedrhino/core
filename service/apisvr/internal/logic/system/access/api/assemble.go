package api

import (
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
)

func ToApiInfoPb(in *types.ApiInfo) *sys.ApiInfo {
	if in == nil {
		return nil
	}
	return &sys.ApiInfo{
		AccessCode: in.AccessCode,
		Id:         in.ID,
		Route:      in.Route,
		Method:     in.Method,
		Name:       in.Name,
		//BusinessType: in.BusinessType,
		//AuthType:     in.AuthType,
		Desc: in.Desc,
	}
}

func ToApiInfosTypes(in []*sys.ApiInfo) (ret []*types.ApiInfo) {
	for _, v := range in {
		ret = append(ret, ToApiInfoTypes(v))
	}
	return
}

func ToApiInfoTypes(in *sys.ApiInfo) *types.ApiInfo {
	if in == nil {
		return nil
	}
	return &types.ApiInfo{
		AccessCode: in.AccessCode,
		ID:         in.Id,
		Route:      in.Route,
		Method:     in.Method,
		Name:       in.Name,
		//BusinessType: in.BusinessType,
		//AuthType:     in.AuthType,
		Desc: in.Desc,
	}
}
