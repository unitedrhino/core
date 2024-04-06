package tenantmanagelogic

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	usermanagelogic "gitee.com/i-Things/core/service/syssvr/internal/logic/usermanage"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"gorm.io/gorm"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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
func (l *TenantInfoCreateLogic) TenantInfoCreate(in *sys.TenantInfoCreateReq) (*sys.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	ctxs.GetUserCtx(l.ctx).AllTenant = true
	defer func() {
		ctxs.GetUserCtx(l.ctx).AllTenant = false
	}()
	userInfo := in.AdminUserInfo
	//首先校验账号格式使用正则表达式，对用户账号做格式校验：只能是大小写字母，数字和下划线，减号
	err := usermanagelogic.CheckUserName(userInfo.UserName)
	if err != nil {
		return nil, err
	}
	//校验密码强度
	err = usermanagelogic.CheckPwd(l.svcCtx, userInfo.Password)
	if err != nil {
		return nil, err
	}
	//1.生成uid
	userID := l.svcCtx.UserID.GetSnowflakeId()
	//2.对密码进行md5加密
	password := utils.MakePwd(userInfo.Password, userID, false)
	ui := relationDB.SysUserInfo{
		TenantCode: stores.TenantCode(in.Info.Code),
		UserID:     userID,
		Phone:      utils.AnyToNullString(userInfo.Phone),
		Email:      utils.AnyToNullString(userInfo.Email),
		UserName:   sql.NullString{String: userInfo.UserName, Valid: true},
		Password:   password,
		NickName:   userInfo.NickName,
		City:       userInfo.City,
		Country:    userInfo.Country,
		Province:   userInfo.Province,
		Language:   userInfo.Language,
		HeadImg:    userInfo.HeadImg,
		Role:       userInfo.Role,
		Sex:        userInfo.Sex,
		IsAllData:  def.True,
	}

	projectPo := relationDB.SysProjectInfo{
		ProjectID:   stores.ProjectID(l.svcCtx.ProjectID.GetSnowflakeId()),
		ProjectName: in.Info.Name,
		//CompanyName: utils.ToEmptyString(in.CompanyName),
		AdminUserID: userID,
		//Region:      utils.ToEmptyString(in.Region),
		//Address:     utils.ToEmptyString(in.Address),
	}

	po := logic.ToTenantInfoPo(in.Info)
	if po.BackgroundImg != "" && in.Info.IsUpdateBackgroundImg {
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessTenantManage, oss.SceneBackgroundImg,
			fmt.Sprintf("%s/%s", po.Code, oss.GetFileNameWithPath(po.BackgroundImg)))
		path, err := l.svcCtx.OssClient.PublicBucket().CopyFromTempBucket(po.BackgroundImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		po.BackgroundImg = path
	}
	if po.LogoImg != "" && in.Info.IsUpdateLogoImg {
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessTenantManage, oss.SceneLogoImg,
			fmt.Sprintf("%s/%s", po.Code, oss.GetFileNameWithPath(po.LogoImg)))
		path, err := l.svcCtx.OssClient.PublicBucket().CopyFromTempBucket(po.LogoImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		po.LogoImg = path
	}
	err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		ris := []*relationDB.SysRoleInfo{{TenantCode: stores.TenantCode(in.Info.Code), Name: "超级管理员"}, {TenantCode: stores.TenantCode(in.Info.Code), Name: "普通用户"}}
		err = relationDB.NewRoleInfoRepo(tx).MultiInsert(l.ctx, ris)
		if err != nil {
			return err
		}
		err := relationDB.NewUserRoleRepo(tx).Insert(l.ctx, &relationDB.SysUserRole{
			TenantCode: stores.TenantCode(in.Info.Code),
			UserID:     ui.UserID,
			RoleID:     ris[0].ID,
		})
		if err != nil {
			return err
		}
		ui.Role = ris[0].ID
		err = relationDB.NewUserInfoRepo(tx).Insert(l.ctx, &ui)
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
		err = relationDB.NewTenantInfoRepo(l.ctx).Insert(l.ctx, po)
		if err != nil {
			return err
		}
		err = relationDB.NewTenantConfigRepo(l.ctx).Insert(l.ctx, &relationDB.SysTenantConfig{
			TenantCode:     stores.TenantCode(in.Info.Code),
			RegisterRoleID: ris[1].ID,
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
	err = l.svcCtx.TenantCache.SetData(l.ctx, po.Code, logic.ToTenantInfoCache(po))
	if err != nil {
		l.Error(err)
	}
	return &sys.WithID{Id: po.ID}, err
}
