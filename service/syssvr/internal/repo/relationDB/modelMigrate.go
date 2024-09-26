package relationDB

import (
	"context"
	"database/sql"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
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
	//{
	//	db := stores.GetCommonConn(context.TODO()).Clauses(clause.OnConflict{DoNothing: true})
	//	if err := db.CreateInBatches(&MigrateDictDetailAdcode, 100).Error; err != nil {
	//		return err
	//	}
	//	if err := db.CreateInBatches(&MigrateDictInfo, 100).Error; err != nil {
	//		return err
	//	}
	//}

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
	MigrateTenantInfo  = []SysTenantInfo{{Code: def.TenantCodeDefault, Name: "默认租户", AdminUserID: adminUserID, DefaultProjectID: defaultProjectID}}
	MigrateUserInfo    = []SysUserInfo{
		{TenantCode: def.TenantCodeDefault, UserID: adminUserID, UserName: sql.NullString{String: "administrator", Valid: true}, Password: "4f0fded4a38abe7a3ea32f898bb82298", Role: 1, NickName: "iThings管理员", IsAllData: def.True},
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
)
