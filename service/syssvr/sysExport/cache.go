package sysExport

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/syssvr/client/areamanage"
	"gitee.com/unitedrhino/core/service/syssvr/client/common"
	"gitee.com/unitedrhino/core/service/syssvr/client/projectmanage"
	"gitee.com/unitedrhino/core/service/syssvr/client/tenantmanage"
	"gitee.com/unitedrhino/core/service/syssvr/client/usermanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/domain/slot"
	"gitee.com/unitedrhino/share/domain/tenant"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/utils"
	"strings"
)

type AreaCacheT = *caches.Cache[areamanage.AreaInfo, int64]

func NewAreaInfoCache(pm areamanage.AreaManage, fastEvent *eventBus.FastEvent) (AreaCacheT, error) {
	return caches.NewCache(caches.CacheConfig[areamanage.AreaInfo, int64]{
		KeyType:   eventBus.ServerCacheKeySysAreaInfo,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key int64) (*areamanage.AreaInfo, error) {
			ret, err := pm.AreaInfoRead(ctx, &sys.AreaInfoReadReq{ProjectID: ctxs.GetUserCtxNoNil(ctx).ProjectID, AreaID: key})
			return ret, err
		},
	})
}

type ProjectCacheT = *caches.Cache[projectmanage.ProjectInfo, int64]

func NewProjectInfoCache(pm projectmanage.ProjectManage, fastEvent *eventBus.FastEvent) (ProjectCacheT, error) {
	return caches.NewCache(caches.CacheConfig[projectmanage.ProjectInfo, int64]{
		KeyType:   eventBus.ServerCacheKeySysProjectInfo,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key int64) (*projectmanage.ProjectInfo, error) {
			ret, err := pm.ProjectInfoRead(ctx, &sys.ProjectWithID{ProjectID: key})
			return ret, err
		},
	})
}

type UserCacheT = *caches.Cache[usermanage.UserInfo, int64]

func NewUserInfoCache(pm usermanage.UserManage, fastEvent *eventBus.FastEvent) (UserCacheT, error) {
	return caches.NewCache(caches.CacheConfig[usermanage.UserInfo, int64]{
		KeyType:   eventBus.ServerCacheKeySysUserInfo,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key int64) (*usermanage.UserInfo, error) {
			ret, err := pm.UserInfoRead(ctx, &sys.UserInfoReadReq{UserID: key})
			return ret, err
		},
	})
}

type TenantCacheT = *caches.Cache[tenant.Info, string]

func NewTenantInfoCache(pm tenantmanage.TenantManage, fastEvent *eventBus.FastEvent) (TenantCacheT, error) {
	return caches.NewCache(caches.CacheConfig[tenant.Info, string]{
		KeyType:   eventBus.ServerCacheKeySysTenantInfo,
		FastEvent: fastEvent,
		GetData: func(ctx context.Context, key string) (*tenant.Info, error) {
			ret, err := pm.TenantInfoRead(ctx, &sys.WithIDCode{Code: key})
			return logic.RpcToTenantInfoCache(ret), err
		},
	})
}

type WebHookCacheT = *caches.Cache[sys.TenantOpenWebHook, string]

func NewTenantOpenWebhookCache(pm tenantmanage.TenantManage, fastEvent *eventBus.FastEvent) (WebHookCacheT, error) {
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

type SlotCacheT = *caches.Cache[slot.Infos, string]

func NewSlotCache(pm common.Common) (SlotCacheT, error) {
	return caches.NewCache(caches.CacheConfig[slot.Infos, string]{
		KeyType: "slot",
		GetData: func(ctx context.Context, key string) (*slot.Infos, error) {
			t := strings.Split(key, ":")
			ret, err := pm.SlotInfoIndex(ctx, &sys.SlotInfoIndexReq{Code: t[0], SubCode: t[1]})
			slots := slot.Infos(utils.CopySlice[slot.Info](ret.List))
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

func GenApiCacheKey(method, route string) string {
	return fmt.Sprintf("%s:%s", method, route)
}
