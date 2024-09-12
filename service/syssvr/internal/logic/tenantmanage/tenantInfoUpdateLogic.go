package tenantmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/oss/common"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantInfoUpdateLogic {
	return &TenantInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新区域
func (l *TenantInfoUpdateLogic) TenantInfoUpdate(in *sys.TenantInfo) (*sys.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	if ctxs.IsRoot(l.ctx) == nil {
		ctxs.GetUserCtx(l.ctx).AllTenant = true
		defer func() {
			ctxs.GetUserCtx(l.ctx).AllTenant = false
		}()
	}

	repo := relationDB.NewTenantInfoRepo(l.ctx)
	old, err := repo.FindOneByFilter(l.ctx, relationDB.TenantInfoFilter{ID: in.Id, Code: in.Code})
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		old.Name = in.Name
	}
	if in.AdminUserID != 0 && in.AdminUserID != old.AdminUserID { //只有default的超管才有权限修改管理员
		err := logic.IsSupperAdmin(l.ctx, def.TenantCodeDefault)
		if err != nil {
			return nil, err
		}
		old.AdminUserID = in.AdminUserID
	}
	if in.BackgroundImg != "" && in.IsUpdateBackgroundImg {
		if old.BackgroundImg != "" {
			err := l.svcCtx.OssClient.PublicBucket().Delete(l.ctx, old.BackgroundImg, common.OptionKv{})
			if err != nil {
				l.Errorf("Delete file err path:%v,err:%v", old.BackgroundImg, err)
			}
		}
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessTenantManage, oss.SceneBackgroundImg,
			fmt.Sprintf("%s/%s", old.Code, oss.GetFileNameWithPath(in.BackgroundImg)))
		path, err := l.svcCtx.OssClient.PublicBucket().CopyFromTempBucket(in.BackgroundImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		old.BackgroundImg = path
	}
	if in.LogoImg != "" && in.IsUpdateLogoImg {
		if old.LogoImg != "" {
			err := l.svcCtx.OssClient.PublicBucket().Delete(l.ctx, old.LogoImg, common.OptionKv{})
			if err != nil {
				l.Errorf("Delete file err path:%v,err:%v", old.LogoImg, err)
			}
		}
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessTenantManage, oss.SceneLogoImg,
			fmt.Sprintf("%s/%s", old.Code, oss.GetFileNameWithPath(in.LogoImg)))
		path, err := l.svcCtx.OssClient.PublicBucket().CopyFromTempBucket(in.LogoImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		old.LogoImg = path
	}
	if in.Footer != "" {
		old.Footer = in.Footer
	}
	if in.Title != "" {
		old.Title = in.Title
	}
	if in.TitleEn != "" {
		old.TitleEn = in.TitleEn
	}
	if in.BackgroundDesc != nil {
		old.BackgroundDesc = in.BackgroundDesc.GetValue()
	}
	if in.BackgroundColour != "" {
		old.BackgroundColour = in.BackgroundColour
	}
	if in.Desc != nil {
		old.Desc = utils.ToEmptyString(in.Desc)
	}
	if in.Status != 0 {
		old.Status = in.Status
	}
	err = repo.Update(l.ctx, old)
	err = caches.SetTenant(l.ctx, logic.ToTenantInfoCache(old))
	if err != nil {
		l.Error(err)
	}
	err = l.svcCtx.TenantCache.SetData(l.ctx, old.Code, logic.ToTenantInfoCache(old))
	if err != nil {
		l.Error(err)
	}
	return &sys.Empty{}, err
}
