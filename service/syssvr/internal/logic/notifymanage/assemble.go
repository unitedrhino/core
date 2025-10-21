package notifymanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
)

func ChannelPoToPb(ctx context.Context, in *relationDB.SysNotifyChannel) *sys.NotifyChannel {
	ret := utils.Copy[sys.NotifyChannel](in)
	uc := ctxs.GetUserCtxNoNil(ctx)
	if ctxs.IsRoot(ctx) == nil {
		return ret
	}
	if uc.TenantCode == ret.TenantCode && ctxs.IsAdmin(ctx) == nil {
		return ret
	}

	ret.Webhook = ""
	ret.Email = nil
	ret.App = nil
	ret.Sms = nil
	return ret
}
func ChannelsPoToPb(ctx context.Context, in []*relationDB.SysNotifyChannel) []*sys.NotifyChannel {
	ret := utils.CopySlice[sys.NotifyChannel](in)
	uc := ctxs.GetUserCtxNoNil(ctx)
	if ctxs.IsRoot(ctx) == nil {
		return ret
	}
	for _, v := range ret {
		if uc.TenantCode == v.TenantCode && ctxs.IsAdmin(ctx) == nil {
			continue
		}
		v.Webhook = ""
		v.Email = nil
		v.App = nil
		v.Sms = nil
	}
	return ret
}
