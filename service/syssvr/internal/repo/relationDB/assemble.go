package relationDB

import "gitee.com/i-Things/core/service/syssvr/domain/slot"

func ToSlotDo(in *SysSlotInfo) *slot.Info {
	if in == nil {
		return nil
	}
	return &slot.Info{
		Code:     in.Code,
		SlotCode: in.SlotCode,
		Method:   in.Method,
		Uri:      in.Uri,
		Hosts:    in.Hosts,
		Body:     in.Body,
		Handler:  in.Handler,
		AuthType: in.AuthType,
	}
}
func ToSlotsDo(in []*SysSlotInfo) (ret slot.Infos) {
	for _, v := range in {
		ret = append(ret, ToSlotDo(v))
	}
	return
}
