package tenantmanagelogic

import (
	"context"
	"fmt"

	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/share/caches"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/core/share/topics"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantInfoCreateLogic {
	return &TenantInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新增租户
func (l *TenantInfoCreateLogic) TenantInfoCreate(in *sys.TenantInfo) (*sys.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	if utils.SliceIn(in.Code, def.TenantCodeCommon, def.TenantCodeDefault) {
		return nil, errors.Parameter.AddMsgf("租户编码不能为内置的 %s或%s", def.TenantCodeCommon, def.TenantCodeDefault)
	}
	_, err := l.svcCtx.TenantCache.GetData(l.ctx, string(in.Code))
	if err == nil {
		return nil, errors.Duplicate.AddMsgf("租户编码已存在:%s", in.Code)
	}
	ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOne(l.ctx, in.AdminUserID)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.UnRegister.AddMsg("请先注册管理员账号")
		}
		return nil, err
	}

	projectPo := &relationDB.SysProjectInfo{
		TenantCode:   dataType.TenantCode(in.Code),
		ProjectID:    dataType.ProjectID(l.svcCtx.ProjectID.GetSnowflakeId()),
		ProjectName:  in.Name,
		IsSysCreated: def.True,
		AdminUserID:  ui.UserID,
	}

	po := logic.ToTenantInfoPo(in)
	if po.BackgroundImg != "" && in.IsUpdateBackgroundImg {
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessTenantManage, oss.SceneBackgroundImg,
			fmt.Sprintf("%s/%s", po.Code, oss.GetFileNameWithPath(po.BackgroundImg)))
		path, err := l.svcCtx.OssClient.PublicBucket().CopyFromTempBucket(po.BackgroundImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		po.BackgroundImg = path
	}
	if po.LogoImg != "" && in.IsUpdateLogoImg {
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessTenantManage, oss.SceneLogoImg,
			fmt.Sprintf("%s/%s", po.Code, oss.GetFileNameWithPath(po.LogoImg)))
		path, err := l.svcCtx.OssClient.PublicBucket().CopyFromTempBucket(po.LogoImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		po.LogoImg = path
	}
	err = TenantCreate(l.ctx, l.svcCtx, projectPo, po, ui, false)
	if err != nil {
		return nil, err
	}
	return &sys.WithID{Id: po.ID}, nil
}

func TenantCreate(ctx context.Context, svcCtx *svc.ServiceContext, projectPo *relationDB.SysProjectInfo,
	po *relationDB.SysTenantInfo, ui *relationDB.SysUserInfo, isCreateUser bool) error {
	err := stores.GetCommonConn(ctx).Transaction(func(tx *gorm.DB) error {
		ris := []*relationDB.SysRoleInfo{{TenantCode: po.Code, Name: "超级管理员", Code: "supper"},
			{TenantCode: po.Code, Name: "管理员", Code: "admin"},
			{TenantCode: po.Code, Name: "普通用户", Code: "client"}}
		err := relationDB.NewRoleInfoRepo(tx).MultiInsert(ctxs.WithRoot(ctx), ris)
		if err != nil {
			return err
		}
		if ui.TenantCode != def.TenantCodeCommon { //租户管理的账号是公共的
			err := relationDB.NewUserInfoRepo(tx).UpdateWithField(ctxs.WithRoot(ctx), relationDB.UserInfoFilter{UserID: ui.UserID}, map[string]any{"tenant_code": def.TenantCodeCommon})
			if err != nil {
				return err
			}
		}
		err = relationDB.NewProjectInfoRepo(tx).Insert(ctxs.WithRoot(ctx), projectPo)
		if err != nil {
			return err
		}
		po.DefaultProjectID = int64(projectPo.ProjectID)
		po.AdminUserID = ui.UserID
		po.AdminRoleID = ris[0].ID
		err = relationDB.NewTenantInfoRepo(tx).Insert(ctxs.WithRoot(ctx), po)
		if err != nil {
			return err
		}
		err = relationDB.NewTenantConfigRepo(tx).Insert(ctxs.WithRoot(ctx), &relationDB.SysTenantConfig{
			TenantCode:     po.Code,
			RegisterRoleID: ris[1].ID,
		})
		if err != nil {
			return err
		}
		err = relationDB.NewDataProjectRepo(tx).MultiInsert(ctxs.WithRoot(ctx), []*relationDB.SysDataProject{
			{
				TenantCode: po.Code,
				ProjectID:  po.DefaultProjectID,
				TargetType: def.TargetUser,
				TargetID:   po.AdminUserID,
				AuthType:   def.AuthAdmin,
			}, {
				TenantCode: po.Code,
				ProjectID:  po.DefaultProjectID,
				TargetType: def.TargetRole,
				TargetID:   po.AdminRoleID,
				AuthType:   def.AuthAdmin,
			},
		})
		if !isCreateUser && ui.TenantCode != def.TenantCodeCommon {
			err = relationDB.NewUserInfoRepo(tx).UpdateWithField(ctxs.WithRoot(ctx), relationDB.UserInfoFilter{UserID: ui.UserID}, map[string]any{"tenant_code": def.TenantCodeCommon})
			if err != nil {
				return err
			}
		}
		// 插入用户角色关联
		err = relationDB.NewUserRoleRepo(tx).Insert(ctxs.WithRoot(ctx), &relationDB.SysUserRole{
			TenantCode: po.Code,
			UserID:     ui.UserID,
			RoleID:     po.AdminRoleID,
		})
		if err != nil {
			return err
		}
		if isCreateUser {
			ui.TenantCode = def.TenantCodeCommon
			err = relationDB.NewUserInfoRepo(tx).Insert(ctxs.WithRoot(ctx), ui)
			if err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return err
	}
	err = caches.SetTenant(ctx, logic.ToTenantInfoCache(po))
	if err != nil {
		logx.WithContext(ctx).Error(err)
	}
	err = svcCtx.TenantCache.SetData(ctx, string(po.Code), logic.ToTenantInfoCache(po))
	if err != nil {
		logx.WithContext(ctx).Error(err)
	}
	svcCtx.FastEvent.Publish(ctx, topics.CoreTenantCreate, po.Code)
	return nil
}
