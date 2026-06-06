package usermanagelogic

import (
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

func userMessagesToProto(pos []*relationDB.SysUserMessage) []*sys.UserMessage {
	ret := make([]*sys.UserMessage, 0, len(pos))
	for _, po := range pos {
		ret = append(ret, userMessageToProto(po))
	}
	return ret
}

func userMessageToProto(po *relationDB.SysUserMessage) *sys.UserMessage {
	if po == nil {
		return nil
	}
	um := &sys.UserMessage{
		Id:     po.ID,
		UserID: po.UserID,
		IsRead: po.IsRead,
	}
	if po.Message != nil {
		um.Message = messageInfoToProto(po.Message)
	}
	return um
}

func messageInfoToProto(mi *relationDB.SysMessageInfo) *sys.MessageInfo {
	if mi == nil {
		return nil
	}
	info := utils.Copy[sys.MessageInfo](mi)
	if mi.TriggerType != "" || mi.TriggerUserID > 0 || mi.TriggerUserNick != "" || mi.TriggerUserAccount != "" {
		info.TriggerUser = &sys.MessageTriggerUser{
			UserId:      mi.TriggerUserID,
			NickName:    mi.TriggerUserNick,
			Account:     mi.TriggerUserAccount,
			TriggerType: mi.TriggerType,
		}
	}
	return info
}
