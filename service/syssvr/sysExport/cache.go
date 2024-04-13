package sysExport

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/client/tenantmanage"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/domain/tenant"
	"gitee.com/i-Things/share/eventBus"
	"strings"
)

func NewTenantInfoCache(pm tenantmanage.TenantManage, fastEvent *eventBus.FastEvent) (*caches.Cache[tenant.Info], error) {
	return caches.NewCache(caches.CacheConfig[tenant.Info]{
		KeyType:   eventBus.ServerCacheKeySysTenantInfo,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key string) (*tenant.Info, error) {
			ret, err := pm.TenantInfoRead(ctx, &sys.WithIDCode{Code: key})
			return logic.RpcToTenantInfoCache(ret), err
		},
	})
}

func NewTenantOpenWebhookCache(pm tenantmanage.TenantManage, fastEvent *eventBus.FastEvent) (*caches.Cache[sys.TenantOpenWebHook], error) {
	return caches.NewCache(caches.CacheConfig[sys.TenantOpenWebHook]{
		KeyType:   eventBus.ServerCacheKeySysTenantOpenWebhook,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key string) (*sys.TenantOpenWebHook, error) {
			t := strings.Split(key, ":")
			ret, err := pm.TenantOpenWebHook(ctx, &sys.WithCode{Code: t[1]})
			return ret, err
		},
	})
}

func GenWebhookCacheKey(ctx context.Context, code string) string {
	tenantCode := ctxs.GetUserCtxNoNil(ctx).TenantCode
	return fmt.Sprintf("%s:%s", tenantCode, code)
}
