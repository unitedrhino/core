package appmanagelogic

import (
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

func ToAppInfoPo(in *sys.AppInfo) *relationDB.SysAppInfo {
	return utils.Copy[relationDB.SysAppInfo](in)
}

func ToAppInfoPb(in *relationDB.SysAppInfo) *sys.AppInfo {
	return utils.Copy[sys.AppInfo](in)
}

func ToAppInfosPb(in []*relationDB.SysAppInfo) []*sys.AppInfo {
	var ret []*sys.AppInfo
	for _, v := range in {
		ret = append(ret, ToAppInfoPb(v))
	}
	return ret
}
