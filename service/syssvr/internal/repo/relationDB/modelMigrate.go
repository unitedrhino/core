package relationDB

import (
	"context"
	"database/sql"

	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/core/share/users"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
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
		&SysDeptSyncJob{},
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
		&SysProjectCrud{},
		&SysAreaProfile{},
	)
	if err != nil {
		return err
	}

	if NeedInitColumn {
		return migrateTableColumn()
	}
	return err
}

func migrateTableColumn() error {
	db := stores.GetCommonConn(context.TODO()).Clauses(clause.OnConflict{DoNothing: true})
	if err := stores.CreateInBatches(db, &MigrateUserInfo, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateRoleInfo, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateUserRole, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateTenantInfo, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateProjectInfo, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateTenantConfig, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateDictDetailAdcode, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateDataProject, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateDictInfo, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateDictDetail, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateAppInfo, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateModuleInfo, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateAppModule, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateTenantApp, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateTenantAppModule, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateNotifyConfig, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateNotifyTemplate, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateNotifyConfigTemplate, 100); err != nil {
		return err
	}
	if err := stores.CreateInBatches(db, &MigrateSlotInfo, 100); err != nil {
		return err
	}
	stores.CreateInBatches(db, &[]SysDeptInfo{{DeletedTime: 666}, {DeletedTime: 777}, {DeletedTime: 888}}, 100)
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
	MigrateTenantInfo  = []SysTenantInfo{{Code: def.TenantCodeDefault, Name: "联犀平台", AdminUserID: adminUserID, AdminRoleID: 3, DefaultProjectID: defaultProjectID}}
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
		//{Code: "client-mini-wx", Name: "c端微信小程序", Type: "mini", SubType: "wx"},
		//{Code: "client-app-android", Name: "客户端安卓", Type: "app", SubType: "android"},
		//{Code: "client-app-ios", Name: "客户端苹果", Type: "app", SubType: "ios"},
	}
	MigrateModuleInfo = []SysModuleInfo{ //后续都通过配置文件去初始化
		//{Code: "systemManage", Type: 1, Order: 2, Name: "系统管理", Path: "system", Url: "", Icon: "icon-menu-xitong", Body: `{}`, HideInMenu: 2, SubType: 3, Tag: 1},
		//{Code: "things", Type: 1, Order: 1, Name: "物联网", Path: "things", Url: "/app/things", Icon: "icon-menu-yingyong2", Body: `{"microAppUrl":"/app/things","microAppName":"物联网","microAppBaseroute":"things"}`, HideInMenu: 2, SubType: 1, Tag: 1},
		//{Code: "myThings", Desc: "企业版", Type: 1, Order: 8, Name: "我的物联", Path: "myThings", Url: "/app/my-things", Icon: "icon-menu-haoyou", Body: `{"microAppUrl":"/app/my-things","microAppName":"我的物联","microAppBaseroute":"myThings"}`, HideInMenu: 2, SubType: 1, Tag: 1},
		//{Code: "lowcode", Desc: "企业版", Type: 1, Order: 8, Name: "低代码", Path: "lowcode", Url: "/app/lowcode", Icon: "icon-menu-yingyong1", Body: `{"microAppUrl":"/app/lowcode","microAppName":"低代码","microAppBaseroute":"lowcode"}`, HideInMenu: 2, SubType: 1, Tag: 1},
	}
	MigrateAppModule = []SysAppModule{ //后续都通过配置文件去初始化
		//{AppCode: "core", ModuleCode: "systemManage"},
		//{AppCode: "core", ModuleCode: "things"},
		//{AppCode: "core", ModuleCode: "myThings"},
		//{AppCode: "core", ModuleCode: "lowcode"},
	}
	MigrateTenantApp = []SysTenantApp{
		{TenantCode: def.TenantCodeDefault, AppCode: "core", LoginTypes: []users.RegType{users.RegPwd}, IsAutoRegister: 1},
		//{TenantCode: def.TenantCodeDefault, AppCode: "client-mini-wx", LoginTypes: []users.RegType{users.RegPwd}, IsAutoRegister: 1},
		//{TenantCode: def.TenantCodeDefault, AppCode: "client-app-android", LoginTypes: []users.RegType{users.RegPwd}, IsAutoRegister: 1},
	}
	MigrateTenantAppModule = []SysTenantAppModule{ //后续都通过配置文件去初始化
		//{TenantCode: def.TenantCodeDefault, SysAppModule: SysAppModule{AppCode: "core", ModuleCode: "systemManage"}},
		//{TenantCode: def.TenantCodeDefault, SysAppModule: SysAppModule{AppCode: "core", ModuleCode: "things"}},
		//{TenantCode: def.TenantCodeDefault, SysAppModule: SysAppModule{AppCode: "core", ModuleCode: "lowcode"}},
		//{TenantCode: def.TenantCodeDefault, SysAppModule: SysAppModule{AppCode: "core", ModuleCode: "myThings"}},
	}

	MigrateNotifyConfig = []SysNotifyConfig{
		{Group: "验证码", Code: "sysUserRegisterCaptcha", Name: "用户注册验证码", SupportTypes: []def.NotifyType{"sms", "email"}, IsRecord: def.False, Params: map[string]string{"code": "验证码", "expr": "过期时间(单位秒,显示分钟)"}},
		{Group: "验证码", Code: "sysUserLoginCaptcha", Name: "用户登录验证码", SupportTypes: []def.NotifyType{"sms", "email"}, IsRecord: def.False, Params: map[string]string{"code": "验证码", "expr": "过期时间(单位秒,显示分钟)"}},
		{Group: "验证码", Code: "sysUserChangePwdCaptcha", Name: "用户修改密码", SupportTypes: []def.NotifyType{"sms", "email"}, IsRecord: def.False, Params: map[string]string{"code": "验证码", "expr": "过期时间(单位秒,显示分钟)"}},
		{Group: "场景联动通知", Code: "ruleScene", Name: "场景联动通知", SupportTypes: []def.NotifyType{"message", "sms", "email", "phoneCall", "dingWebhook", "wxEWebHook", "wxMini", "dingTalk", "dingMini"}, EnableTypes: []def.NotifyType{"message"}, IsRecord: def.True, Params: map[string]string{"body": "内容", "title": "标题"}},
		{Group: "设备", Code: "ruleDeviceAlarm", Name: "设备告警通知", SupportTypes: []def.NotifyType{"sms", "email", "dingWebhook"}, IsRecord: def.True, Params: map[string]string{"productID": "产品ID(若为设备触发)", "deviceName": "触发设备ID(若为设备触发)", "deviceAlias": "设备名称(若为设备触发)", "sceneName": "触发场景名称"}},
		{Group: "系统公告", Code: "sysAnnouncement", Name: "系统公告", SupportTypes: []def.NotifyType{"sms", "email", "wxMini"}, IsRecord: def.True, Params: map[string]string{"body": "内容", "title": "标题"}},
	}
	// 场景联动站内信默认模板（notifyCode=ruleScene, type=message）
	MigrateNotifyTemplate = []SysNotifyTemplate{
		{ID: 1, TenantCode: dataType.TenantCode("common"), Name: "场景联动站内信", NotifyCode: "ruleScene", Type: "message", TemplateCode: "ruleScene_message", Subject: "{{.title}}", Body: "{{.body}}", Desc: "场景联动站内信默认模板"},
	}
	// 默认平台绑定 ruleScene 站内信模板
	MigrateNotifyConfigTemplate = []SysNotifyConfigTemplate{
		{TenantCode: def.TenantCodeDefault, NotifyCode: "ruleScene", Type: "message", TemplateID: 1},
	}

	MigrateSlotInfo = []SysSlotInfo{}

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
		}, {
			Name:  "项目类型",
			Code:  "projectType",
			Group: def.DictGroupBase,
			Desc:  "",
		},
	}
	MigrateDictDetail = []SysDictDetail{
		{DictCode: "dictGroup", Label: "基础配置", Value: def.DictGroupBase},
		{DictCode: "dictGroup", Label: "物联网", Value: def.DictGroupThings},
		{DictCode: "dictGroup", Label: "系统管理", Value: def.DictGroupSystem},
		{DictCode: "projectType", Label: "默认", Value: "default"},
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
