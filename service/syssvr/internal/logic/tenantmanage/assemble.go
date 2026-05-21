package tenantmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/share/utils"
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

// tenantThirdToPb 将数据库第三方配置转为 proto（utils.Copy 无法跨类型拷贝嵌套结构）
func tenantThirdToPb(src *relationDB.SysTenantThird) *sys.ThirdAppConfig {
	if src == nil {
		return nil
	}
	if src.AppID == "" && src.AppKey == "" && src.AppSecret == "" {
		return nil
	}
	return &sys.ThirdAppConfig{
		AppID:     src.AppID,
		AppKey:    src.AppKey,
		AppSecret: src.AppSecret,
	}
}

// tenantAppleToPb 将数据库 Apple 配置转为 proto
func tenantAppleToPb(src *relationDB.SysTenantAppleConfig) *sys.AppleAppConfig {
	if src == nil {
		return nil
	}
	if src.AppID == "" && src.TeamID == "" && src.KeyID == "" && src.PrivateKey == "" && src.RedirectURI == "" {
		return nil
	}
	return &sys.AppleAppConfig{
		AppID:       src.AppID,
		TeamID:      src.TeamID,
		KeyID:       src.KeyID,
		PrivateKey:  src.PrivateKey,
		RedirectURI: src.RedirectURI,
	}
}

func ToTenantApp(ctx context.Context, svcCtx *svc.ServiceContext, po *relationDB.SysTenantApp) *sys.TenantAppInfo {
	ret := &sys.TenantAppInfo{
		Id:             po.ID,
		Code:           string(po.TenantCode),
		AppCode:        po.AppCode,
		LoginTypes:     po.LoginTypes,
		IsAutoRegister: po.IsAutoRegister,
		Config:         po.Config,
		DingMini:       tenantThirdToPb(po.DingMini),
		WxMini:         tenantThirdToPb(po.WxMini),
		WxOpen:         tenantThirdToPb(po.WxOpen),
		Huawei:         tenantThirdToPb(po.Huawei),
		Google:         tenantThirdToPb(po.Google),
		Github:         tenantThirdToPb(po.Github),
		Apple:          tenantAppleToPb(po.Apple),
	}
	if po.Android != nil {
		ret.Android = utils.Copy[sys.ThirdApp](po.Android)
		if ret.Android != nil && ret.Android.FilePath != "" {
			var err error
			ret.Android.FilePath, err = svcCtx.OssClient.PublicBucket().GetUrl(ret.Android.FilePath, false)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
			}
		}
	}
	return ret
}
