package relationDB

import (
	"gitee.com/unitedrhino/core/share/domain/slot"
	"gitee.com/unitedrhino/share/utils"
)

func ToSlotDo(in *SysSlotInfo) *slot.Info {
	return utils.Copy[slot.Info](in)
}
func ToSlotsDo(in []*SysSlotInfo) (ret slot.Infos) {
	for _, v := range in {
		ret = append(ret, ToSlotDo(v))
	}
	return
}
