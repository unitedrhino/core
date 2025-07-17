package menu

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

func ToMenuInfoRpc(in *types.MenuInfo) *sys.MenuInfo {
	return utils.Copy[sys.MenuInfo](in)
}
func ToMenuInfosRpc(in []*types.MenuInfo) (ret []*sys.MenuInfo) {
	for _, v := range in {
		ret = append(ret, ToMenuInfoRpc(v))
	}
	return
}
