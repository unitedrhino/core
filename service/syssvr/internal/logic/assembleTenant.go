package logic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/domain/tenant"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func ToTenantInfoRpc(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.SysTenantInfo) *sys.TenantInfo {
	if in.BackgroundImg != "" {
		var err error
		in.BackgroundImg, err = svcCtx.OssClient.PublicBucket().SignedGetUrl(ctx, in.BackgroundImg, 24*60*60, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}
	if in.LogoImg != "" {
		var err error
		in.LogoImg, err = svcCtx.OssClient.PublicBucket().SignedGetUrl(ctx, in.LogoImg, 24*60*60, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}
	return utils.Copy[sys.TenantInfo](in)
}
func ToTenantInfosRpc(ctx context.Context, svcCtx *svc.ServiceContext, in []*relationDB.SysTenantInfo) (ret []*sys.TenantInfo) {
	for _, v := range in {
		ret = append(ret, ToTenantInfoRpc(ctx, svcCtx, v))
	}
	return
}

func ToTenantInfoPo(in *sys.TenantInfo) *relationDB.SysTenantInfo {
	return utils.Copy[relationDB.SysTenantInfo](in)
}

func ToTenantInfoCaches(in []*relationDB.SysTenantInfo) (ret []*tenant.Info) {
	for _, v := range in {
		ret = append(ret, ToTenantInfoCache(v))
	}
	return ret
}

func ToTenantInfoCache(in *relationDB.SysTenantInfo) *tenant.Info {
	return utils.Copy[tenant.Info](in)
}

//func CacheToTenantInfoRpc(ctx context.Context, svcCtx *svc.ServiceContext, in *tenant.Info) *sys.TenantInfo {
//	if in.BackgroundImg != "" {
//		var err error
//		in.BackgroundImg, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, in.BackgroundImg, 24*60*60, common.OptionKv{})
//		if err != nil {
//			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
//		}
//	}
//	if in.LogoImg != "" {
//		var err error
//		in.LogoImg, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, in.LogoImg, 24*60*60, common.OptionKv{})
//		if err != nil {
//			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
//		}
//	}
//	return utils.Copy[sys.TenantInfo](in)
//}

func RpcToTenantInfoCache(in *sys.TenantInfo) *tenant.Info {
	return utils.Copy[tenant.Info](in)
}
