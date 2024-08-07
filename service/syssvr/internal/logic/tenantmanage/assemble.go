package tenantmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/oss/common"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func ToTenantConfigPb(ctx context.Context, svcCtx *svc.ServiceContext, po *relationDB.SysTenantConfig) *sys.TenantConfig {
	for _, p := range po.RegisterAutoCreateProject {
		for _, a := range p.Areas {
			if a.AreaImg != "" {
				var err error
				a.AreaImg, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, a.AreaImg, 24*60*60, common.OptionKv{})
				if err != nil {
					logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
				}
			}
		}
	}
	return utils.Copy[sys.TenantConfig](po)
}
