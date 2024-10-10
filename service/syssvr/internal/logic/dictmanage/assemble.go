package dictmanagelogic

import (
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

func ToDictInfoPb(in *relationDB.SysDictInfo) *sys.DictInfo {
	return utils.Copy[sys.DictInfo](in)
}

func ToDictInfosPb(in []*relationDB.SysDictInfo) (list []*sys.DictInfo) {
	for _, v := range in {
		list = append(list, ToDictInfoPb(v))
	}
	return
}
func ToDictDetailsPb(in []*relationDB.SysDictDetail) []*sys.DictDetail {
	return utils.CopySlice[sys.DictDetail](in)

}
