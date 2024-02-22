package startup

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func Init(svcCtx *svc.ServiceContext) {
	ctx := context.Background()
	utils.Go(ctx, func() {
		list, err := relationDB.NewTenantInfoRepo(ctx).FindByFilter(ctx, relationDB.TenantInfoFilter{}, nil)
		logx.Must(err)
		err = caches.InitTenant(ctx, logic.ToTenantInfoCaches(list)...)
		logx.Must(err)
	})
}
