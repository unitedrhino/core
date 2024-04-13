package startup

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
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
	tenantCache, err := caches.NewCache(caches.CacheConfig[tenant.Info]{
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
		ExpireTime: 10 * time.Minute,
	})
	logx.Must(err)
	svcCtx.TenantCache = tenantCache
}
