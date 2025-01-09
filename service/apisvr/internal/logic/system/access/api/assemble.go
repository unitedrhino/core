package api

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

func ToApiInfoPb(in *types.ApiInfo) *sys.ApiInfo {
	return utils.Copy[sys.ApiInfo](in)
}

func ToApiInfosTypes(in []*sys.ApiInfo) (ret []*types.ApiInfo) {
	for _, v := range in {
		ret = append(ret, ToApiInfoTypes(v))
	}
	return
}

func ToApiInfoTypes(in *sys.ApiInfo) *types.ApiInfo {
	return utils.Copy[types.ApiInfo](in)
}
