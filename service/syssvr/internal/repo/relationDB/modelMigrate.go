package relationDB

import (
	"context"
	"database/sql"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/users"
	"github.com/spf13/cast"
	"gorm.io/gorm/clause"
)

var NeedInitColumn bool

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(context.TODO())
	if !db.Migrator().HasTable(&SysUserInfo{}) {
		//需要初始化表
		NeedInitColumn = true
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

	if NeedInitColumn {
		return migrateTableColumn()
	} else {
		db := stores.GetCommonConn(context.TODO()).Clauses(clause.OnConflict{DoNothing: true})
		if err := db.CreateInBatches(&MigrateDictInfo, 100).Error; err != nil {
			return err
		}
		if err := db.CreateInBatches(&MigrateDictDetail, 100).Error; err != nil {
			return err
		}
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

	if err := db.CreateInBatches(&MigrateDictInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateDictDetail, 100).Error; err != nil {
		return err
	}

	if err := db.CreateInBatches(&MigrateAppInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateModuleInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateAppModule, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateTenantApp, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateTenantAppModule, 100).Error; err != nil {
		return err
	}
	db.Create(&SysDeptInfo{ID: 3, Name: "锚点"})
	db.Delete(&SysDeptInfo{ID: 3, Name: "锚点"})
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
	MigrateAppInfo = []SysAppInfo{
		{Code: "core", Name: "管理后台", Type: "web", SubType: "web"},
		{Code: "client-mini-wx", Name: "c端微信小程序", Type: "mini", SubType: "wx"},
		{Code: "client-mini-wx", Name: "c端微信小程序", Type: "mini", SubType: "wx"},
		{Code: "client-app-android", Name: "客户端安卓", Type: "app", SubType: "android"},
		{Code: "client-app-ios", Name: "客户端苹果", Type: "app", SubType: "ios"},
	}
	MigrateModuleInfo = []SysModuleInfo{
		{Code: "systemManage", Type: 1, Order: 2, Name: "系统管理", Path: "system", Url: "", Icon: "icon-menu-xitong", Body: `{}`, HideInMenu: 2, SubType: 3, Tag: 1},
		{Code: "things", Type: 1, Order: 1, Name: "物联网", Path: "things", Url: "/app/things", Icon: "icon-menu-yingyong2", Body: `"{""microAppUrl"":""/app/things"",""microAppName"":""物联网"",""microAppBaseroute"":""things""}"`, HideInMenu: 2, SubType: 1, Tag: 1},
		{Code: "myThings", Type: 1, Order: 8, Name: "我的物联", Path: "myThings", Url: "/app/my-things", Icon: "icon-menu-haoyou", Body: `"{""microAppUrl"":""/app/my-things"",""microAppName"":""我的物联"",""microAppBaseroute"":""myThings""}"`, HideInMenu: 2, SubType: 1, Tag: 1},
	}
	MigrateAppModule = []SysAppModule{
		{AppCode: "core", ModuleCode: "systemManage"},
		{AppCode: "core", ModuleCode: "things"},
		{AppCode: "core", ModuleCode: "myThings"},
	}
	MigrateTenantApp = []SysTenantApp{
		{TenantCode: def.TenantCodeDefault, AppCode: "core", LoginTypes: []users.RegType{users.RegPwd}, IsAutoRegister: 1},
		{TenantCode: def.TenantCodeDefault, AppCode: "client-mini-wx", LoginTypes: []users.RegType{users.RegPwd}, IsAutoRegister: 1},
		{TenantCode: def.TenantCodeDefault, AppCode: "client-app-android", LoginTypes: []users.RegType{users.RegPwd}, IsAutoRegister: 1},
	}
	MigrateTenantAppModule = []SysTenantAppModule{
		{TenantCode: def.TenantCodeDefault, SysAppModule: SysAppModule{AppCode: "core", ModuleCode: "systemManage"}},
		{TenantCode: def.TenantCodeDefault, SysAppModule: SysAppModule{AppCode: "core", ModuleCode: "things"}},
		{TenantCode: def.TenantCodeDefault, SysAppModule: SysAppModule{AppCode: "core", ModuleCode: "myThings"}},
	}

	MigrateRoleInfo = []SysRoleInfo{
		{ID: 1, TenantCode: def.TenantCodeDefault, Name: "管理员", Code: def.RoleCodeAdmin},
		{ID: 2, TenantCode: def.TenantCodeDefault, Name: "普通用户", Code: def.RoleCodeClient, Desc: "C端用户"},
		{ID: 3, TenantCode: def.TenantCodeDefault, Name: "超级管理员", Code: def.RoleCodeSupper}}

	MigrateDictInfo = []SysDictInfo{
		{
			Name:  "错误",
			Code:  "error",
			Group: def.DictGroupBase,
			Desc:  "系统返回的错误code和对应的描述",
		}, {
			Name:       "区划",
			Code:       "adcode",
			Group:      def.DictGroupThings,
			Desc:       "中国区划",
			StructType: 2,
		}, {
			Name:  "字典分组",
			Code:  "dictGroup",
			Group: def.DictGroupBase,
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
