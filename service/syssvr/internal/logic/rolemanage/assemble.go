package rolemanagelogic

import (
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
)

func ToRoleInfoRpc(in *relationDB.SysRoleInfo) *sys.RoleInfo {
	if in == nil {
		return nil
	}
	return &sys.RoleInfo{
		Id:          in.ID,
		Name:        in.Name,
		Desc:        in.Desc,
		CreatedTime: in.CreatedTime.Unix(),
		Status:      in.Status,
	}
}
func ToRoleInfosRpc(in []*relationDB.SysRoleInfo) []*sys.RoleInfo {
	var ret []*sys.RoleInfo
	for _, v := range in {
		ret = append(ret, ToRoleInfoRpc(v))
	}
	return ret
}
