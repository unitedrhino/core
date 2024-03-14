package sysExport

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/client/tenantmanage"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/domain/tenant"
	"gitee.com/i-Things/share/eventBus"
)

func NewTenantInfoCache(pm tenantmanage.TenantManage, fastEvent *eventBus.FastEvent) (*caches.Cache[tenant.Info], error) {
	return caches.NewCache(caches.CacheConfig[tenant.Info]{
		KeyType:   eventBus.ServerCacheKeySysTenant,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key string) (*tenant.Info, error) {
			ret, err := pm.TenantInfoRead(ctx, &sys.WithIDCode{Code: key})
			return logic.RpcToTenantInfoCache(ret), err
		},
	})
}
