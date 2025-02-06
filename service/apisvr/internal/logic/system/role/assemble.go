package role

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

func ToRoleInfoTypes(in *sys.RoleInfo) *types.RoleInfo {
	return utils.Copy[types.RoleInfo](in)
}
func ToRoleInfosTypes(in []*sys.RoleInfo) (ret []*types.RoleInfo) {
	for _, v := range in {
		ret = append(ret, ToRoleInfoTypes(v))
	}
	return
}
