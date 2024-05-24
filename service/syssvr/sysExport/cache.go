package sysExport

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/client/common"
	"gitee.com/i-Things/core/service/syssvr/client/projectmanage"
	"gitee.com/i-Things/core/service/syssvr/client/tenantmanage"
	"gitee.com/i-Things/core/service/syssvr/client/usermanage"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/domain/slot"
	"gitee.com/i-Things/share/domain/tenant"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/utils"
	"strings"
)

func NewProjectInfoCache(pm projectmanage.ProjectManage, fastEvent *eventBus.FastEvent) (*caches.Cache[projectmanage.ProjectInfo, int64], error) {
	return caches.NewCache(caches.CacheConfig[projectmanage.ProjectInfo, int64]{
		KeyType:   eventBus.ServerCacheKeySysProjectInfo,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key int64) (*projectmanage.ProjectInfo, error) {
			ret, err := pm.ProjectInfoRead(ctx, &sys.ProjectWithID{ProjectID: key})
			return ret, err
		},
	})
}

func NewUserInfoCache(pm usermanage.UserManage, fastEvent *eventBus.FastEvent) (*caches.Cache[usermanage.UserInfo, int64], error) {
	return caches.NewCache(caches.CacheConfig[usermanage.UserInfo, int64]{
		KeyType:   eventBus.ServerCacheKeySysUserInfo,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key int64) (*usermanage.UserInfo, error) {
			ret, err := pm.UserInfoRead(ctx, &sys.UserInfoReadReq{UserID: key})
			return ret, err
		},
	})
}

func NewTenantInfoCache(pm tenantmanage.TenantManage, fastEvent *eventBus.FastEvent) (*caches.Cache[tenant.Info, string], error) {
	return caches.NewCache(caches.CacheConfig[tenant.Info, string]{
		KeyType:   eventBus.ServerCacheKeySysTenantInfo,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key string) (*tenant.Info, error) {
			ret, err := pm.TenantInfoRead(ctx, &sys.WithIDCode{Code: key})
			return logic.RpcToTenantInfoCache(ret), err
		},
	})
}

func NewTenantOpenWebhookCache(pm tenantmanage.TenantManage, fastEvent *eventBus.FastEvent) (*caches.Cache[sys.TenantOpenWebHook, string], error) {
	return caches.NewCache(caches.CacheConfig[sys.TenantOpenWebHook, string]{
		KeyType:   eventBus.ServerCacheKeySysTenantOpenWebhook,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key string) (*sys.TenantOpenWebHook, error) {
			t := strings.Split(key, ":")
			ret, err := pm.TenantOpenWebHook(ctx, &sys.WithCode{Code: t[1]})
			return ret, err
		},
	})
}

func NewSlotCache(pm common.Common) (*caches.Cache[slot.Infos, string], error) {
	return caches.NewCache(caches.CacheConfig[slot.Infos, string]{
		KeyType: "slot",
		GetData: func(ctx context.Context, key string) (*slot.Infos, error) {
			t := strings.Split(key, ":")
			ret, err := pm.SlotInfoIndex(ctx, &sys.SlotInfoIndexReq{Code: t[0], SubCode: t[1]})
			slots := slot.Infos(utils.CopySlice[slot.Info](ret.Slots))
			return &slots, err
		},
	})
}

func GenSlotCacheKey(code string, subCode string) string {
	return fmt.Sprintf("%s:%s", code, subCode)
}

func GenWebhookCacheKey(ctx context.Context, code string) string {
	tenantCode := ctxs.GetUserCtxNoNil(ctx).TenantCode
	return fmt.Sprintf("%s:%s", tenantCode, code)
}
