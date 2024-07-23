package cache

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/users"
	"gitee.com/i-Things/share/utils"
	"github.com/spf13/cast"
	"time"
)

type UserToken struct {
	*caches.Cache[users.UserInfo, int64]
}

func NewUserToken(FastEvent *eventBus.FastEvent) (*UserToken, error) {
	c, err := caches.NewCache(caches.CacheConfig[users.UserInfo, int64]{
		KeyType:   eventBus.ServerCacheKeySysUserTokenInfo,
		FastEvent: FastEvent,
		GetData: func(ctx context.Context, key int64) (*users.UserInfo, error) {
			ui, err := relationDB.NewUserInfoRepo(ctx).FindOneByFilter(ctx, relationDB.UserInfoFilter{
				UserIDs:    []int64{key},
				WithRoles:  true,
				WithTenant: true,
			})
			if err != nil {
				return nil, err
			}
			var rolses []int64
			var roleCodes []string
			var isAdmin int64 = def.False
			for _, v := range ui.Roles {
				rolses = append(rolses, v.RoleID)
				if v.Role != nil && v.Role.Code != "" {
					roleCodes = append(roleCodes, v.Role.Code)
				}
			}

			if ui.Tenant != nil && (utils.SliceIn(ui.Tenant.AdminRoleID, rolses...) || ui.Tenant.AdminUserID == ui.UserID) {
				isAdmin = def.True
			}
			var account = ui.UserName.String
			if account == "" {
				account = ui.Phone.String
			}
			if account == "" {
				account = ui.Email.String
			}
			if account == "" {
				account = cast.ToString(ui.UserID)
			}
			uii := users.UserInfo{
				UserID:     ui.UserID,
				Account:    account,
				RoleIDs:    rolses,
				RoleCodes:  roleCodes,
				TenantCode: string(ui.TenantCode),
				IsAdmin:    isAdmin,
				IsAllData:  ui.IsAllData,
			}
			return &uii, nil
		},
		ExpireTime: 10 * time.Minute,
	})
	if err != nil {
		return nil, err
	}
	return &UserToken{
		Cache: c,
	}, nil
}
