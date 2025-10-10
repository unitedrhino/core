package info

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

func ToModuleInfoRpc(in *types.ModuleInfo) *sys.ModuleInfo {
	return utils.Copy[sys.ModuleInfo](in)
}
func ToModuleInfoApi(in *sys.ModuleInfo) *types.ModuleInfo {
	return utils.Copy[types.ModuleInfo](in)

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
