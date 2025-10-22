package cache

import (
	"context"
	"time"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/domain/tenant"
	"gitee.com/unitedrhino/core/share/topics"
	"gitee.com/unitedrhino/core/share/users"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
)

type UserCache struct {
	*caches.Cache[users.UserInfo, users.UserTenantCore]
}

func NewUserCache(FastEvent *eventBus.FastEvent, tenantCache *caches.Cache[tenant.Info, string], userCache *caches.Cache[sys.UserInfo, int64]) (*UserCache, error) {
	c, err := caches.NewCache(caches.CacheConfig[users.UserInfo, users.UserTenantCore]{
		KeyType:   topics.ServerCacheKeySysUserTokenInfo,
		FastEvent: FastEvent,
		GetData: func(ctx context.Context, key users.UserTenantCore) (*users.UserInfo, error) {
			ui, err := userCache.GetData(ctx, key.UserID)
			if err != nil {
				return nil, err
			}
			var ut = ui.Tenants[0]
			for _, v := range ui.Tenants {
				if v.TenantCode == key.TenantCode {
					ut = v
					break
				}
			}
			roles, err := relationDB.NewUserRoleRepo(ctx).FindByFilter(ctx,
				relationDB.UserRoleFilter{TenantCode: key.TenantCode, UserID: key.UserID, WithRole: true}, nil)
			if err != nil {
				return nil, err
			}
			var rolses []int64
			var roleCodes []string
			var isAdmin int64 = def.False
			for _, v := range roles {
				rolses = append(rolses, v.RoleID)
				if v.Role != nil && v.Role.Code != "" {
					roleCodes = append(roleCodes, v.Role.Code)
				}
			}
			Tenant, err := tenantCache.GetData(ctx, ut.TenantCode)
			if err != nil {
				return nil, err
			}
			if Tenant != nil && (utils.SliceIn(Tenant.AdminRoleID, rolses...) || Tenant.AdminUserID == ui.UserID) {
				isAdmin = def.True
			}
			var account = ui.UserName
			if account == "" {
				account = ui.Phone.GetValue()
			}
			if account == "" {
				account = ui.Email.GetValue()
			}
			if account == "" {
				account = cast.ToString(ui.UserID)
			}
			uii := users.UserInfo{
				UserID:      ui.UserID,
				LastTokenID: ui.LastTokenID,
				Account:     account,
				RoleIDs:     rolses,
				RoleCodes:   roleCodes,
				TenantCode:  string(ut.TenantCode),
				IsAdmin:     isAdmin,
			}
			return &uii, nil
		},
		ExpireTime: 20 * time.Minute,
	})
	if err != nil {
		return nil, err
	}
	return &UserCache{
		Cache: c,
	}, nil
}
