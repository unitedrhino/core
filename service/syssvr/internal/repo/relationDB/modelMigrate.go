package relationDB

import (
	"context"
	"database/sql"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/domain/slot"
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
		&SysAppPolicy{},
		&SysRoleApp{},
		&SysUserRole{},
		&SysTenantInfo{},
		&SysTenantOpenWebhook{},
		&SysTenantOpenAccess{},
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

	if err := db.CreateInBatches(&MigrateRoleMenu, 100).Error; err != nil {
		return err
	}

	//if err := db.CreateInBatches(&MigrateRoleApi, 100).Error; err != nil {
	//	return err
	//}
	if err := db.CreateInBatches(&MigrateUserRole, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateRoleApp, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateAppInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateTenantInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateProjectInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateTenantApp, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateTenantConfig, 100).Error; err != nil {
		return err
	}

	if err := db.CreateInBatches(&MigrateNotifyInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateNotifyTemplate, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateTenantNotify, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateSlotInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateDictDetailAdcode, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateDictInfo, 100).Error; err != nil {
		return err
	}

	//{
	//	if err := db.CreateInBatches(&MigrateModuleApi, 100).Error; err != nil {
	//		return err
	//	}
	//	for _, v := range MigrateModuleApi {
	//		data := SysTenantAppApi{
	//			TempLateID:   v.ID,
	//			TenantCode:   def.TenantCodeDefault,
	//			AppCode:      def.AppCore,
	//			SysModuleApi: v,
	//		}
	//		data.ID = 0
	//		MigrateTenantAppApi = append(MigrateTenantAppApi, data)
	//	}
	//	if err := db.CreateInBatches(&MigrateTenantAppApi, 100).Error; err != nil {
	//		return err
	//	}
	//}
	{
		if err := db.CreateInBatches(&MigrateModuleMenu, 100).Error; err != nil {
			return err
		}
		for _, v := range MigrateModuleMenu {
			data := SysTenantAppMenu{
				TenantCode:    def.TenantCodeDefault,
				SysModuleMenu: v,
			}
			data.ID = 0
			MigrateTenantAppMenu = append(MigrateTenantAppMenu, data)
		}
		if err := db.CreateInBatches(&MigrateModuleMenu, 100).Error; err != nil {
			return err
		}
	}
	{
		if err := db.CreateInBatches(&MigrateAppModule, 100).Error; err != nil {
			return err
		}
		for _, v := range MigrateAppModule {
			MigrateTenantAppModule = append(MigrateTenantAppModule, SysTenantAppModule{
				TenantCode:   def.TenantCodeDefault,
				SysAppModule: v,
			})
		}
		if err := db.CreateInBatches(&MigrateTenantAppModule, 100).Error; err != nil {
			return err
		}
	}

	{
		if err := db.CreateInBatches(&MigrateNotifyInfo, 100).Error; err != nil {
			return err
		}
		if err := db.CreateInBatches(&MigrateNotifyTemplate, 100).Error; err != nil {
			return err
		}
		if err := db.CreateInBatches(&MigrateTenantNotify, 100).Error; err != nil {
			return err
		}
	}

	return nil
}

func init() {
	for i := int64(1); i <= 100; i++ {
		MigrateRoleMenu = append(MigrateRoleMenu, SysRoleMenu{
			TenantCode: def.TenantCodeDefault,
			RoleID:     1,
			AppCode:    def.AppCore,
			MenuID:     i,
		})
	}
}

const (
	adminUserID      = 1740358057038188544
	defaultProjectID = 1786838173980422144
)

// 子应用管理员可以配置自己子应用的角色

var (
	MigrateNotifyInfo = []SysNotifyConfig{
		{Group: def.NotifyGroupCaptcha, Code: def.NotifyCodeSysUserRegisterCaptcha, Name: "用户注册验证码",
			SupportTypes: []string{def.NotifyTypeSms, def.NotifyTypeEmail}, IsRecord: def.False,
			Params: map[string]string{"code": "验证码code"}},
		{Group: def.NotifyGroupCaptcha, Code: def.NotifyCodeSysUserLoginCaptcha, Name: "用户登录验证码",
			SupportTypes: []string{def.NotifyTypeSms, def.NotifyTypeEmail}, IsRecord: def.False,
			Params: map[string]string{"code": "验证码code"}},

		{Group: def.NotifyGroupDevice, Code: def.NotifyCodeRuleScene, Name: "场景联动通知",
			SupportTypes: []string{def.NotifyTypeSms, def.NotifyTypeEmail, def.NotifyTypeDingTalk}, IsRecord: def.True,
			Params: map[string]string{"body": "通知的内容"}},
		{Group: def.NotifyGroupDevice, Code: def.NotifyCodeDeviceAlarm, Name: "设备告警通知",
			SupportTypes: []string{def.NotifyTypeSms, def.NotifyTypeEmail, def.NotifyTypeDingTalk}, IsRecord: def.True,
			Params: map[string]string{"body": "通知的内容"}},
	}
	MigrateNotifyTemplate = []SysNotifyTemplate{
		{
			ID:           1,
			TenantCode:   def.TenantCodeDefault,
			Name:         "用户注册验证码",
			NotifyCode:   def.NotifyCodeSysUserRegisterCaptcha,
			Type:         def.NotifyTypeSms,
			TemplateCode: "SMS_288215142",
			SignName:     "EbelongTool",
			Subject:      "注册验证码",
			Body:         "欢迎注册,你的验证码是:{{.code}},有效期为{{.expr}}分钟",
		},
		{
			ID:           2,
			TenantCode:   def.TenantCodeDefault,
			Name:         "登录验证码",
			NotifyCode:   def.NotifyCodeSysUserLoginCaptcha,
			Type:         def.NotifyTypeSms,
			TemplateCode: "SMS_288215142",
			SignName:     "EbelongTool",
			Subject:      "登录验证码",
			Body:         "欢迎登录,你的验证码是:{{.code}},有效期为{{.expr}}分钟",
		},
		{
			ID:           3,
			TenantCode:   def.TenantCodeDefault,
			Name:         "场景通知",
			NotifyCode:   def.NotifyCodeRuleScene,
			Type:         def.NotifyTypeSms,
			TemplateCode: "SMS_465414256",
			SignName:     "EbelongTool",
			Subject:      "场景通知",
			Body:         "你好,场景联动通知,内容如下:{{.body}}",
		},
		{
			ID:           4,
			TenantCode:   def.TenantCodeDefault,
			Name:         "设备告警通知",
			NotifyCode:   def.NotifyCodeDeviceAlarm,
			Type:         def.NotifyTypeSms,
			TemplateCode: "SMS_465344291",
			SignName:     "EbelongTool",
			Subject:      "设备告警通知",
			Body:         "你好,{{.deviceAlias}}设备告警:{{.body}}",
		},
	}

	MigrateTenantNotify = []SysNotifyConfigTemplate{
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeSysUserRegisterCaptcha, Type: def.NotifyTypeSms, TemplateID: 1},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeSysUserRegisterCaptcha, Type: def.NotifyTypeEmail, TemplateID: 0},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeSysUserLoginCaptcha, Type: def.NotifyTypeSms, TemplateID: 2},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeSysUserLoginCaptcha, Type: def.NotifyTypeEmail, TemplateID: 0},

		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeRuleScene, Type: def.NotifyTypeSms, TemplateID: 3},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeRuleScene, Type: def.NotifyTypeEmail, TemplateID: 1},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeRuleScene, Type: def.NotifyTypeDingTalk, TemplateID: 1},

		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeDeviceAlarm, Type: def.NotifyTypeSms, TemplateID: 4},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeDeviceAlarm, Type: def.NotifyTypeEmail, TemplateID: 1},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeDeviceAlarm, Type: def.NotifyTypeDingTalk, TemplateID: 1},
	}
	MigrateTenantNotifyChannel = []SysNotifyChannel{}

	MigrateModuleInfo = []SysModuleInfo{
		{Name: "系统管理", Code: def.ModuleSystemManage},
		{Name: "租户管理", Code: def.ModuleTenantManage},
		{Name: "物联网", Code: def.ModuleThings},
		{Name: "音视频", Code: def.ModuleVideo},
		{Name: "大屏", Code: def.ModuleView},
	}
	MigrateAppModule = []SysAppModule{
		{
			AppCode:    def.AppCore,
			ModuleCode: def.ModuleThings,
		},
		{
			AppCode:    def.AppCore,
			ModuleCode: def.ModuleSystemManage,
		},
		{
			AppCode:    def.AppCore,
			ModuleCode: def.ModuleTenantManage,
		},
		{
			AppCode:    def.AppCore,
			ModuleCode: def.ModuleView,
		},
		{
			AppCode:    def.AppCore,
			ModuleCode: def.ModuleVideo,
		},
	}
	MigrateTenantAppModule = []SysTenantAppModule{}
	MigrateTenantAppMenu   = []SysTenantAppMenu{}
	MigrateTenantConfig    = []SysTenantConfig{
		{TenantCode: def.TenantCodeDefault, RegisterRoleID: 2},
	}
	MigrateProjectInfo = []SysProjectInfo{{TenantCode: def.TenantCodeDefault, AdminUserID: adminUserID, ProjectID: defaultProjectID, ProjectName: "默认项目"}}
	MigrateTenantInfo  = []SysTenantInfo{{Code: def.TenantCodeDefault, Name: "默认租户", AdminUserID: adminUserID, DefaultProjectID: defaultProjectID}}
	MigrateTenantApp   = []SysTenantApp{{TenantCode: def.TenantCodeDefault, AppCode: def.AppCore}}
	MigrateUserInfo    = []SysUserInfo{
		{TenantCode: def.TenantCodeDefault, UserID: adminUserID, UserName: sql.NullString{String: "administrator", Valid: true}, Password: "4f0fded4a38abe7a3ea32f898bb82298", Role: 1, NickName: "iThings管理员", IsAllData: def.True},
	}
	MigrateUserRole = []SysUserRole{
		{TenantCode: def.TenantCodeDefault, UserID: adminUserID, RoleID: 1},
	}
	MigrateRoleInfo = []SysRoleInfo{
		{ID: 1, TenantCode: def.TenantCodeDefault, Name: "admin", Code: "admin"},
		{ID: 2, TenantCode: def.TenantCodeDefault, Name: "client", Code: "client", Desc: "C端用户"}}
	MigrateRoleMenu []SysRoleMenu
	MigrateRoleApp  = []SysRoleApp{
		{RoleID: 1, TenantCode: def.TenantCodeDefault, AppCode: def.AppCore},
	}
	MigrateAppInfo = []SysAppInfo{
		{Code: def.AppCore, Name: "中台"},
		{Code: def.AppAll, Name: "全部"},
	}

	MigrateSlotInfo = []SysSlotInfo{
		{Code: slot.CodeAreaInfo, SubCode: slot.SubCodeCreate, SlotCode: slot.SlotCodeIthings, Method: "POST", Uri: "/api/v1/things/slot/area/create", Hosts: []string{"http://127.0.0.1:7788"}, Body: `{"projectID":"{{.ProjectID}}","areaID":"{{.AreaID}}","parentAreaID":"{{.ParentAreaID}}"}`, Handler: nil, AuthType: "core", Desc: ""},
		{Code: slot.CodeAreaInfo, SubCode: slot.SubCodeDelete, SlotCode: slot.SlotCodeIthings, Method: "POST", Uri: "/api/v1/things/slot/area/delete", Hosts: []string{"http://127.0.0.1:7788"}, Body: `{"projectID":"{{.ProjectID}}","areaID":"{{.AreaID}}","parentAreaID":"{{.ParentAreaID}}"}`, Handler: nil, AuthType: "core", Desc: ""},
		{Code: slot.CodeUserSubscribe, SubCode: def.UserSubscribeDevicePropertyReport, SlotCode: slot.SlotCodeIthings, Method: "POST", Uri: "/api/v1/things/slot/user/subscribe", Hosts: []string{"http://127.0.0.1:7788"}, Body: ``, Handler: nil, AuthType: "core", Desc: ""},
		{Code: slot.CodeUserSubscribe, SubCode: def.UserSubscribeDeviceConn, SlotCode: slot.SlotCodeIthings, Method: "POST", Uri: "/api/v1/things/slot/user/subscribe", Hosts: []string{"http://127.0.0.1:7788"}, Body: ``, Handler: nil, AuthType: "core", Desc: ""},
	}

	MigrateModuleMenu = []SysModuleMenu{
		{ID: 2, ParentID: 1, Type: 1, Order: 2, Name: "设备管理", Path: "/deviceManagers", Component: "./deviceManagers/index.tsx", Icon: "icon_data_01", Redirect: "", HideInMenu: def.False},
		{ID: 6, ParentID: 2, Type: 1, Order: 1, Name: "产品", Path: "/deviceManagers/product/index", Component: "./deviceManagers/product/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 7, ParentID: 2, Type: 1, Order: 1, Name: "产品详情", Path: "/deviceManagers/product/detail/:id", Component: "./deviceManagers/product/detail/index", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
		{ID: 8, ParentID: 2, Type: 1, Order: 2, Name: "设备", Path: "/deviceManagers/device/index", Component: "./deviceManagers/device/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 9, ParentID: 2, Type: 1, Order: 2, Name: "设备详情", Path: "/deviceManagers/device/detail/:id/:name/:type", Component: "./deviceManagers/device/detail/index", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
		{ID: 23, ParentID: 2, Type: 1, Order: 3, Name: "分组", Path: "/deviceManagers/group/index", Component: "./deviceManagers/group/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 24, ParentID: 2, Type: 1, Order: 3, Name: "分组详情", Path: "/deviceManagers/group/detail/:id", Component: "./deviceManagers/group/detail/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.True},

		{ID: 4, ParentID: 1, Type: 1, Order: 4, Name: "运维监控", Path: "/operationsMonitorings", Component: "./operationsMonitorings/index.tsx", Icon: "icon_hvac", Redirect: "", HideInMenu: def.False},
		{ID: 13, ParentID: 4, Type: 1, Order: 1, Name: "固件升级", Path: "/operationsMonitorings/firmwareUpgrade/index", Component: "./operationsMonitorings/firmwareUpgrade/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 15, ParentID: 4, Type: 1, Order: 3, Name: "资源管理", Path: "/operationsMonitorings/resourceManagement/index", Component: "./operationsMonitorings/resourceManagement/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 16, ParentID: 4, Type: 1, Order: 4, Name: "远程配置", Path: "/operationsMonitorings/remoteConfiguration/index", Component: "./operationsMonitorings/remoteConfiguration/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 18, ParentID: 4, Type: 1, Order: 6, Name: "在线调试", Path: "/operationsMonitorings/onlineDebug/index", Component: "./operationsMonitorings/onlineDebug/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},

		{ID: 25, ParentID: 4, Type: 1, Order: 7, Name: "日志服务", Path: "/operationsMonitorings/logService/index", Component: "./operationsMonitorings/logService/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 35, ParentID: 1, Type: 1, Order: 1, Name: "首页", Path: "/home", Component: "./home/index.tsx", Icon: "icon_dosing", Redirect: "", HideInMenu: def.False},

		//{ID: 43, ParentID: 1, Type: 1, Order: 5, Name: "告警管理", Path: "/alarmManagers", Component: "./alarmManagers/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.False},
		//{ID: 44, ParentID: 43, Type: 1, Order: 1, Name: "告警配置", Path: "/alarmManagers/alarmConfiguration/index", Component: "./alarmManagers/alarmConfiguration/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.False},
		//{ID: 53, ParentID: 43, Type: 1, Order: 5, Name: "新增告警配置", Path: "/alarmManagers/alarmConfiguration/save", Component: "./alarmManagers/alarmConfiguration/addAlarmConfig/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.True},
		//{ID: 54, ParentID: 43, Type: 1, Order: 5, Name: "告警日志", Path: "/alarmManagers/alarmConfiguration/log/detail/:id/:level", Component: "./alarmManagers/alarmLog/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.True},
		//{ID: 45, ParentID: 43, Type: 1, Order: 5, Name: "告警记录", Path: "/alarmManagers/alarmConfiguration/log", Component: "./alarmManagers/alarmRecord/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.False},
		//{ID: 50, ParentID: 1, Type: 1, Order: 5, Name: "规则引擎", Path: "/ruleEngine", Component: "./ruleEngine/index.tsx", Icon: "icon_dosing", Redirect: "", HideInMenu: def.False},
		//{ID: 51, ParentID: 50, Type: 1, Order: 1, Name: "场景联动", Path: "/ruleEngine/scene/index", Component: "./ruleEngine/scene/index.tsx", Icon: "icon_device", Redirect: "", HideInMenu: def.False},

		{ID: 60, ParentID: 3, Type: 2, Order: 1, Name: "内嵌", Path: "/systemManagers/iframe", Component: "https://www.douyu.com", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 61, ParentID: 3, Type: 3, Order: 1, Name: "外链", Path: "/systemManagers/links", Component: "https://ant.design", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 70, ParentID: 3, Type: 1, Order: 1, Name: "任务管理", Path: "/systemManagers/timed", Component: "./systemManagers/timed/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 71, ParentID: 70, Type: 1, Order: 1, Name: "任务组", Path: "/systemManagers/timed/group", Component: "./systemManagers/timed/group/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 72, ParentID: 70, Type: 1, Order: 1, Name: "任务组详情", Path: "/systemManagers/timed/group/detail/:id", Component: "./systemManagers/timed/group/detail/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
		{ID: 73, ParentID: 70, Type: 1, Order: 1, Name: "任务", Path: "/systemManagers/timed/task", Component: "./systemManagers/timed/task/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 74, ParentID: 70, Type: 1, Order: 1, Name: "任务详情", Path: "/systemManagers/timed/task/detail/:id", Component: "./systemManagers/timed/task/detail/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
		{ID: 38, ParentID: 3, Type: 1, Order: 5, Name: "日志管理", Path: "/systemManagers/log", Component: "./systemManagers/log/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 39, ParentID: 38, Type: 1, Order: 1, Name: "操作日志", Path: "/systemManagers/log/operationLog/index", Component: "./systemManagers/log/operationLog/index.tsx", Icon: "icon_dosing", Redirect: "", HideInMenu: def.False},
		{ID: 41, ParentID: 38, Type: 1, Order: 2, Name: "登录日志", Path: "/systemManagers/log/loginLog/index", Component: "./systemManagers/log/loginLog/index", Icon: "icon_heat", Redirect: "", HideInMenu: def.False},
		{ID: 42, ParentID: 3, Type: 1, Order: 4, Name: "接口管理", Path: "/systemManagers/api/index", Component: "./systemManagers/api/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 10, ParentID: 3, Type: 1, Order: 1, Name: "用户管理", Path: "/systemManagers/user/index", Component: "./systemManagers/user/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 11, ParentID: 3, Type: 1, Order: 2, Name: "角色管理", Path: "/systemManagers/role/index", Component: "./systemManagers/role/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 12, ParentID: 3, Type: 1, Order: 3, Name: "菜单列表", Path: "/systemManagers/menu/index", Component: "./systemManagers/menu/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
		{ID: 3, ParentID: 1, Type: 1, Order: 9, Name: "系统管理", Path: "/systemManagers", Component: "./systemManagers/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},

		//视频服务菜单项
		{ID: 63, ParentID: 1, Type: 1, Order: 2, Name: "视频服务", Path: "/videoManagers", Component: "./videoManagers", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
		{ID: 64, ParentID: 63, Type: 1, Order: 1, Name: "流服务管理", Path: "/videoManagers/vidsrvmgr/index", Component: "./videoManagers/vidsrvmgr/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
		{ID: 65, ParentID: 63, Type: 1, Order: 3, Name: "视频流广场", Path: "/videoManagers/plaza/index", Component: "./videoManagers/plaza/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
		{ID: 66, ParentID: 63, Type: 1, Order: 2, Name: "视频流管理", Path: "/videoManagers/vidstream/index", Component: "./videoManagers/vidstream/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
		{ID: 67, ParentID: 63, Type: 1, Order: 4, Name: "视频回放", Path: "/videoManagers/playback/index", Component: "./videoManagers/playback/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
		{ID: 68, ParentID: 63, Type: 1, Order: 2, Name: "录像计划", Path: "/videoManagers/recordplan/index", Component: "./videoManagers/recordplan/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
		{ID: 69, ParentID: 63, Type: 1, Order: 1, Name: "流服务详细", Path: "/videoManagers/vidsrvmgr/detail/:id", Component: "./videoManagers/vidsrvmgr/detail/index", Icon: "icon_heat", Redirect: "", HideInMenu: 1},
		{ID: 75, ParentID: 63, Type: 1, Order: 1, Name: "视频流详细", Path: "/videoManagers/vidstream/detail/:id", Component: "./videoManagers/vidstream/detail/index", Icon: "icon_heat", Redirect: "", HideInMenu: 1},
	}

	MigrateDictInfo = []SysDictInfo{
		{
			Name:  "区划",
			Group: "基础配置",
			Code:  "adcode",
		},
	}
)
