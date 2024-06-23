package startup

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	projectmanagelogic "gitee.com/i-Things/core/service/syssvr/internal/logic/projectmanage"
	usermanagelogic "gitee.com/i-Things/core/service/syssvr/internal/logic/usermanage"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/domain/tenant"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func Init(svcCtx *svc.ServiceContext) {
	ctx := context.Background()
	utils.Go(ctx, func() {
		list, err := relationDB.NewTenantInfoRepo(ctx).FindByFilter(ctx, relationDB.TenantInfoFilter{}, nil)
		logx.Must(err)
		err = caches.InitTenant(ctx, logic.ToTenantInfoCaches(list)...)
		logx.Must(err)
	})
	InitCache(svcCtx)
}

func InitCache(svcCtx *svc.ServiceContext) {
	tenantCache, err := caches.NewCache(caches.CacheConfig[tenant.Info, string]{
		KeyType:   eventBus.ServerCacheKeySysTenantInfo,
		FastEvent: svcCtx.ServerMsg,
		GetData: func(ctx context.Context, key string) (*tenant.Info, error) {
			db := relationDB.NewTenantInfoRepo(ctx)
			if key == "" {
				key = ctxs.GetUserCtxNoNil(ctx).TenantCode
			}
			pi, err := db.FindOneByFilter(ctx, relationDB.TenantInfoFilter{
				Codes: []string{key}})
			pb := logic.ToTenantInfoCache(pi)
			return pb, err
		},
		ExpireTime: 20 * time.Minute,
	})
	logx.Must(err)
	svcCtx.TenantCache = tenantCache

	userCache, err := caches.NewCache(caches.CacheConfig[sys.UserInfo, int64]{
		KeyType:   eventBus.ServerCacheKeySysUserInfo,
		FastEvent: svcCtx.ServerMsg,
		GetData: func(ctx context.Context, key int64) (*sys.UserInfo, error) {
			db := relationDB.NewUserInfoRepo(ctx)
			if key == 0 {
				key = ctxs.GetUserCtxNoNil(ctx).UserID
			}
			pi, err := db.FindOne(ctx, key)
			pb := usermanagelogic.UserInfoToPb(ctx, pi, svcCtx)
			return pb, err
		},
		ExpireTime: 20 * time.Minute,
	})
	logx.Must(err)
	svcCtx.UserCache = userCache

	projectCache, err := caches.NewCache(caches.CacheConfig[sys.ProjectInfo, int64]{
		KeyType:   eventBus.ServerCacheKeySysProjectInfo,
		FastEvent: svcCtx.ServerMsg,
		GetData: func(ctx context.Context, key int64) (*sys.ProjectInfo, error) {
			db := relationDB.NewProjectInfoRepo(ctx)
			if key == 0 {
				key = ctxs.GetUserCtxNoNil(ctx).ProjectID
			}
			pi, err := db.FindOne(ctx, key)
			pb := projectmanagelogic.ProjectInfoToPb(ctx, svcCtx, pi)
			return pb, err
		},
		ExpireTime: 20 * time.Minute,
	})
	logx.Must(err)
	svcCtx.ProjectCache = projectCache
}
