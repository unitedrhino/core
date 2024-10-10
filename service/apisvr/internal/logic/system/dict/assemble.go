package dict

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

func ToDetailPb(in *types.DictDetail) *sys.DictDetail {
	return utils.Copy[sys.DictDetail](in)
}

func ToInfoPb(in *types.DictInfo) *sys.DictInfo {
	return utils.Copy[sys.DictInfo](in)
}

func ToDetailTypes(in *sys.DictDetail) *types.DictDetail {
	return utils.Copy[types.DictDetail](in)
}

func ToDetailsTypes(in []*sys.DictDetail) (ret []*types.DictDetail) {
	for _, v := range in {
		ret = append(ret, ToDetailTypes(v))
	}
	return
}

func ToInfoTypes(in *sys.DictInfo) *types.DictInfo {
	return utils.Copy[types.DictInfo](in)
}
func ToInfosTypes(in []*sys.DictInfo) (ret []*types.DictInfo) {
	for _, v := range in {
		ret = append(ret, ToInfoTypes(v))
	}
	return
}
