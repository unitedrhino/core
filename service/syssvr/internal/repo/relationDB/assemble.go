package relationDB

import (
	"gitee.com/i-Things/share/domain/slot"
	"gitee.com/i-Things/share/utils"
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
