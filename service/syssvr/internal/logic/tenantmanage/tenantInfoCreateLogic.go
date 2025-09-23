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

	projectPo := relationDB.SysProjectInfo{
		TenantCode:  dataType.TenantCode(in.Code),
		ProjectID:   dataType.ProjectID(l.svcCtx.ProjectID.GetSnowflakeId()),
		ProjectName: in.Name,
		//CompanyName: utils.ToEmptyString(in.CompanyName),
		AdminUserID: ui.UserID,
		//Region:      utils.ToEmptyString(in.Region),
		//Address:     utils.ToEmptyString(in.Address),
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
	err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		ris := []*relationDB.SysRoleInfo{{TenantCode: dataType.TenantCode(in.Code), Name: "超级管理员", Code: "supper"},
			{TenantCode: dataType.TenantCode(in.Code), Name: "管理员", Code: "admin"},
			{TenantCode: dataType.TenantCode(in.Code), Name: "普通用户", Code: "client"}}
		err = relationDB.NewRoleInfoRepo(tx).MultiInsert(l.ctx, ris)
		if err != nil {
			return err
		}
		err := relationDB.NewUserRoleRepo(tx).Insert(l.ctx, &relationDB.SysUserRole{
			TenantCode: dataType.TenantCode(in.Code),
			UserID:     ui.UserID,
			RoleID:     ris[0].ID,
		})
		if err != nil {
			return err
		}
		err = relationDB.NewUserTenantRepo(tx).Insert(l.ctx, &relationDB.SysUserTenant{
			TenantCode: dataType.TenantCode(in.Code),
			UserID:     ui.UserID,
			Status:     def.True,
		})
		if err != nil {
			return err
		}
		err = relationDB.NewProjectInfoRepo(tx).Insert(l.ctx, &projectPo)
		if err != nil {
			return err
		}
		po.DefaultProjectID = int64(projectPo.ProjectID)
		po.AdminUserID = ui.UserID
		po.AdminRoleID = ris[0].ID
		err = relationDB.NewTenantInfoRepo(tx).Insert(l.ctx, po)
		if err != nil {
			return err
		}
		err = relationDB.NewTenantConfigRepo(tx).Insert(l.ctx, &relationDB.SysTenantConfig{
			TenantCode:     dataType.TenantCode(in.Code),
			RegisterRoleID: ris[1].ID,
		})
		if err != nil {
			return err
		}
		err = relationDB.NewDataProjectRepo(tx).MultiInsert(l.ctx, []*relationDB.SysDataProject{
			{
				TenantCode: dataType.TenantCode(in.Code),
				ProjectID:  po.DefaultProjectID,
				TargetType: def.TargetUser,
				TargetID:   po.AdminUserID,
				AuthType:   def.AuthAdmin,
			}, {
				TenantCode: dataType.TenantCode(in.Code),
				ProjectID:  po.DefaultProjectID,
				TargetType: def.TargetRole,
				TargetID:   po.AdminRoleID,
				AuthType:   def.AuthAdmin,
			},
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	err = caches.SetTenant(l.ctx, logic.ToTenantInfoCache(po))
	if err != nil {
		l.Error(err)
	}
	err = l.svcCtx.TenantCache.SetData(l.ctx, string(po.Code), logic.ToTenantInfoCache(po))
	if err != nil {
		l.Error(err)
	}
	l.svcCtx.FastEvent.Publish(l.ctx, topics.CoreTenantCreate, po.Code)
	return &sys.WithID{Id: po.ID}, nil
}

func CheckAdmin(ctx context.Context, svcCtx *svc.ServiceContext, userName string, nickName string, password string, accounts ...string) error {
	if userName == "" {
		return errors.Parameter.AddMsg("用户名必填")
	}
	if password == "" {
		return errors.Parameter.AddMsg("密码必填")
	}
	err := logic.CheckPwd(svcCtx, password)
	if err != nil {
		return err
	}
	err = logic.CheckUserName(userName)
	if err != nil {
		return err
	}
	//_, err = relationDB.NewTenantInfoRepo(ctx).FindOneByFilter(ctxs.WithRoot(ctx), relationDB.TenantInfoFilter{
	//	Code: userName,
	//})
	//if err == nil {
	//	return errors.Parameter.AddMsg("用户名已被使用")
	//}
	//if !errors.Cmp(err, errors.NotFind) {
	//	return err
	//}

	po, err := relationDB.NewUserInfoRepo(ctx).FindOneByFilter(ctxs.WithRoot(ctx), relationDB.UserInfoFilter{
		WithTenants: true,
		Accounts:    append([]string{userName}, accounts...),
	})
	if err == nil {
		if po.UserName.Valid && po.UserName.String == userName {
			return errors.Parameter.AddMsg("用户名已被使用")
		}
		if po.Email.Valid && utils.SliceIn(po.Email.String, accounts...) {
			return errors.Parameter.AddMsg("邮箱已被使用")
		}
		if po.Phone.Valid && utils.SliceIn(po.Phone.String, accounts...) {
			return errors.Parameter.AddMsg("手机号已被使用")
		}
		return errors.Parameter.AddMsg("账号已被占用")
	}
	if !errors.Cmp(err, errors.NotFind) {
		return err
	}
	return nil
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
			TenantCode:       po.Code,
			RegisterRoleID:   ris[1].ID,
			DeviceLimit:      svcCtx.Config.TenantConfig.DeviceLimit,
			UserLimit:        svcCtx.Config.TenantConfig.UserLimit,
			ProjectLimit:     svcCtx.Config.TenantConfig.ProjectLimit,
			OperLogKeepDays:  svcCtx.Config.TenantConfig.OperLogKeepDays,
			LoginLogKeepDays: svcCtx.Config.TenantConfig.LoginLogKeepDays,
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

		// 插入用户角色关联
		err = relationDB.NewUserRoleRepo(tx).Insert(ctxs.WithRoot(ctx), &relationDB.SysUserRole{
			TenantCode: po.Code,
			UserID:     ui.UserID,
			RoleID:     po.AdminRoleID,
		})
		if err != nil {
			return err
		}

		err = relationDB.NewUserInfoRepo(tx).Insert(ctxs.WithRoot(ctx), ui)
		if err != nil {
			return err
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
