package relationDB

import (
	"context"
	"database/sql"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/spf13/cast"
	"gorm.io/gorm/clause"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(context.TODO())
	var needInitColumn bool
	if !db.Migrator().HasTable(&SysUserInfo{}) {
		//需要初始化表
		needInitColumn = true
	}
	err := db.AutoMigrate(
		&SysDeptUser{},
		&SysServiceInfo{},
		&SysDeptInfo{},
		&SysUserMessage{},
		&SysMessageInfo{},
		&SysNotifyConfig{},
		&SysNotifyTemplate{},
		&SysNotifyConfigTemplate{},
		&SysDictInfo{},
		&SysDictDetail{},
		&SysSlotInfo{},
		&SysUserInfo{},
		&SysRoleInfo{},
		&SysRoleMenu{},
		&SysRoleAccess{},
		&SysTenantAgreement{},
		&SysRoleModule{},
		&SysModuleMenu{},
		&SysLoginLog{},
		&SysOperLog{},
		&SysApiInfo{},
		&SysAccessInfo{},
		&SysAreaInfo{},
		&SysProjectInfo{},
		&SysOpsWorkOrder{},
		&SysOpsFeedback{},
		&SysDataArea{},
		&SysDataProject{},
		&SysAppInfo{},
		&SysRoleApp{},
		&SysUserRole{},
		&SysTenantInfo{},
		&SysTenantOpenWebhook{},
		&SysDataOpenAccess{},
		&SysTenantApp{},
		&SysTenantAccess{},
		&SysTenantConfig{},
		&SysModuleInfo{},
		&SysAppModule{},
		&SysTenantAppMenu{},
		&SysTenantAppModule{},
		&SysUserAreaApply{},
		&SysUserProfile{},
		&SysNotifyChannel{},
		&SysProjectProfile{},
		&SysAreaProfile{},
	)
	if err != nil {
		return err
	}
	{
		db := stores.GetCommonConn(context.TODO()).Clauses(clause.OnConflict{DoNothing: true})
		if err := db.CreateInBatches(&MigrateDictInfo, 100).Error; err != nil {
			return err
		}
		if err := db.CreateInBatches(&MigrateDictDetail, 100).Error; err != nil {
			return err
		}

	}

	if needInitColumn {
		return migrateTableColumn()
	}
	return err
}

func migrateTableColumn() error {
	db := stores.GetCommonConn(context.TODO()).Clauses(clause.OnConflict{DoNothing: true})
	if err := db.CreateInBatches(&MigrateUserInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateRoleInfo, 100).Error; err != nil {
		return err
	}

	if err := db.CreateInBatches(&MigrateUserRole, 100).Error; err != nil {
		return err
	}

	if err := db.CreateInBatches(&MigrateTenantInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateProjectInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateTenantConfig, 100).Error; err != nil {
		return err
	}

	if err := db.CreateInBatches(&MigrateDictDetailAdcode, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateDataProject, 100).Error; err != nil {
		return err
	}

	if err := db.Create(&SysDeptInfo{ID: 3, Name: "锚点"}).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateDictInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateDictDetail, 100).Error; err != nil {
		return err
	}

	return nil
}

const (
	adminUserID      = 1740358057038188544
	defaultProjectID = 1786838173980422144
)

// 子应用管理员可以配置自己子应用的角色

var (
	MigrateTenantAppMenu = []SysTenantAppMenu{}
	MigrateTenantConfig  = []SysTenantConfig{
		{TenantCode: def.TenantCodeDefault, RegisterRoleID: 2},
	}
	MigrateProjectInfo = []SysProjectInfo{{TenantCode: def.TenantCodeDefault, AdminUserID: adminUserID, ProjectID: defaultProjectID, ProjectName: "默认项目"}}
	MigrateDataProject = []SysDataProject{{ProjectID: defaultProjectID, TargetType: def.TargetRole, TargetID: 1, AuthType: def.AuthAdmin}}
	MigrateTenantInfo  = []SysTenantInfo{{Code: def.TenantCodeDefault, Name: "默认租户", AdminUserID: adminUserID, DefaultProjectID: defaultProjectID}}
	MigrateUserInfo    = []SysUserInfo{
		{TenantCode: def.TenantCodeDefault, UserID: adminUserID, UserName: sql.NullString{String: "administrator", Valid: true}, Password: "4f0fded4a38abe7a3ea32f898bb82298", Role: 1, NickName: "联犀管理员", IsAllData: def.True},
	}
	MigrateUserRole = []SysUserRole{
		{TenantCode: def.TenantCodeDefault, UserID: adminUserID, RoleID: 1},
		{TenantCode: def.TenantCodeDefault, UserID: adminUserID, RoleID: 2},
		{TenantCode: def.TenantCodeDefault, UserID: adminUserID, RoleID: 3},
	}
	MigrateRoleInfo = []SysRoleInfo{
		{ID: 1, TenantCode: def.TenantCodeDefault, Name: "管理员", Code: def.RoleCodeAdmin},
		{ID: 2, TenantCode: def.TenantCodeDefault, Name: "普通用户", Code: def.RoleCodeClient, Desc: "C端用户"},
		{ID: 3, TenantCode: def.TenantCodeDefault, Name: "超级管理员", Code: def.RoleCodeSupper}}

	MigrateDictInfo = []SysDictInfo{
		{
			Name:  "错误",
			Code:  "error",
			Group: "基础配置",
			Desc:  "系统返回的错误code和对应的描述",
		}, {
			Name:  "区划",
			Code:  "adcode",
			Group: "基础配置",
			Desc:  "中国区划",
		}, {
			Name:  "字典分组",
			Code:  "dictGroup",
			Group: "基础配置",
			Desc:  "字典的分组",
		},
	}
	MigrateDictDetail = []SysDictDetail{
		{DictCode: "dictGroup", Label: "基础配置", Value: def.DictGroupBase},
		{DictCode: "dictGroup", Label: "物联网", Value: def.DictGroupThings},
		{DictCode: "dictGroup", Label: "系统管理", Value: def.DictGroupSystem},
	}
)

func init() {
	for code, msg := range errors.ErrorMap {
		MigrateDictDetail = append(MigrateDictDetail, SysDictDetail{
			DictCode: "error",
			Label:    msg,
			Value:    cast.ToString(code),
			Status:   def.True,
		})
	}
	return
}
