package info

import (
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"
)

func ToAppInfoRpc(in *types.AppInfo) *sys.AppInfo {
	return utils.Copy[sys.AppInfo](in)
}

func ToAppInfoTypes(in *sys.AppInfo) *types.AppInfo {
	return utils.Copy[types.AppInfo](in)
}

func ToAppInfosTypes(in []*sys.AppInfo) []*types.AppInfo {
	var ret []*types.AppInfo
	for _, v := range in {
		ret = append(ret, ToAppInfoTypes(v))
	}
	return ret
}
