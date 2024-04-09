package relationDB

import (
	"context"
	"database/sql"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm/clause"
	"net/http"
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
		&SysNotifyInfo{},
		&SysNotifyTemplate{},
		&SysTenantNotify{},
		&SysDictInfo{},
		&SysDictDetail{},
		&SysSlotInfo{},
		&SysUserInfo{},
		&SysRoleInfo{},
		&SysRoleMenu{},
		&SysRoleAccess{},
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
		&SysTenantApp{},
		&SysTenantAccess{},
		&SysTenantConfig{},
		&SysModuleInfo{},
		&SysAppModule{},
		&SysTenantAppMenu{},
		&SysTenantAppModule{},
		&SysUserAreaApply{},
	)
	if err != nil {
		return err
	}
	//{
	//	db := stores.GetCommonConn(context.TODO()).Clauses(clause.OnConflict{DoNothing: true})
	//	if err := db.CreateInBatches(&MigrateNotifyConfig, 100).Error; err != nil {
	//		return err
	//	}
	//	if err := db.CreateInBatches(&MigrateNotifyTemplate, 100).Error; err != nil {
	//		return err
	//	}
	//	if err := db.CreateInBatches(&MigrateTenantNotifyTemplate, 100).Error; err != nil {
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
	if err := db.CreateInBatches(&MigrateAccessInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateApiInfo, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateNotifyConfig, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateNotifyTemplate, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateTenantNotifyTemplate, 100).Error; err != nil {
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
		if err := db.CreateInBatches(&MigrateNotifyConfig, 100).Error; err != nil {
			return err
		}
		if err := db.CreateInBatches(&MigrateNotifyTemplate, 100).Error; err != nil {
			return err
		}
		if err := db.CreateInBatches(&MigrateTenantNotifyTemplate, 100).Error; err != nil {
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
	MigrateNotifyConfig = []SysNotifyInfo{
		{Group: def.NotifyGroupCaptcha, Code: def.NotifyCodeSysUserRegisterCaptcha, Name: "用户注册验证码",
			SupportTypes: []string{def.NotifyTypeSms, def.NotifyTypeEmail}, IsRecord: def.False,
			DefaultSubject: "注册验证码", DefaultBody: "欢迎注册,你的验证码是:{{.code}},有效期为{{.expr}}分钟",
			DefaultTemplateCode: "SMS_288215142", DefaultSignName: "EbelongTool",
			Params: map[string]string{"code": "验证码code"}},
		{Group: def.NotifyGroupCaptcha, Code: def.NotifyCodeSysUserLoginCaptcha, Name: "用户登录验证码",
			DefaultSubject: "登录验证码", DefaultBody: "欢迎登录,你的验证码是:{{.code}},有效期为{{.expr}}分钟",
			DefaultTemplateCode: "SMS_288215142", DefaultSignName: "EbelongTool",
			SupportTypes: []string{def.NotifyTypeSms, def.NotifyTypeEmail}, IsRecord: def.False,
			Params: map[string]string{"code": "验证码code"}},

		{Group: def.NotifyGroupDevice, Code: def.NotifyCodeRuleScene, Name: "场景联动通知",
			DefaultSubject: "场景通知", DefaultBody: "你好,场景联动通知,内容如下:{{.body}}",
			DefaultTemplateCode: "SMS_465414256", DefaultSignName: "EbelongTool",
			SupportTypes: []string{def.NotifyTypeSms, def.NotifyTypeEmail, def.NotifyTypeDingTalk}, IsRecord: def.True,
			Params: map[string]string{"body": "通知的内容"}},
		{Group: def.NotifyGroupDevice, Code: def.NotifyCodeDeviceAlarm, Name: "设备告警通知",
			DefaultSubject: "设备告警通知", DefaultBody: "你好,设备告警,告警级别{{.level}}:{{.body}}",
			DefaultTemplateCode: "SMS_465344291", DefaultSignName: "EbelongTool",
			SupportTypes: []string{def.NotifyTypeSms, def.NotifyTypeEmail, def.NotifyTypeDingTalk}, IsRecord: def.True,
			Params: map[string]string{"body": "通知的内容"}},
	}
	MigrateNotifyTemplate       = []SysNotifyTemplate{}
	MigrateTenantNotifyTemplate = []SysTenantNotify{
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeSysUserRegisterCaptcha, Type: def.NotifyTypeSms, TemplateID: 1},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeSysUserRegisterCaptcha, Type: def.NotifyTypeEmail, TemplateID: 1},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeSysUserLoginCaptcha, Type: def.NotifyTypeSms, TemplateID: 1},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeSysUserLoginCaptcha, Type: def.NotifyTypeEmail, TemplateID: 1},

		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeRuleScene, Type: def.NotifyTypeSms, TemplateID: 1},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeRuleScene, Type: def.NotifyTypeEmail, TemplateID: 1},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeRuleScene, Type: def.NotifyTypeDingTalk, TemplateID: 1},

		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeDeviceAlarm, Type: def.NotifyTypeSms, TemplateID: 1},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeDeviceAlarm, Type: def.NotifyTypeEmail, TemplateID: 1},
		{TenantCode: def.TenantCodeDefault, NotifyCode: def.NotifyCodeDeviceAlarm, Type: def.NotifyTypeDingTalk, TemplateID: 1},
	}

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
		{TenantCode: def.TenantCodeDefault, RegisterRoleID: 2, Email: &SysTenantEmail{
			From:     "godlei6@qq.com",
			Host:     "smtp.qq.com",
			Secret:   "xxx",
			Nickname: "验证码机器人",
			Port:     465,
			IsSSL:    def.True},
		},
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
		{ID: 1, TenantCode: def.TenantCodeDefault, Name: "admin"},
		{ID: 2, TenantCode: def.TenantCodeDefault, Name: "client", Desc: "C端用户"}}
	MigrateRoleMenu []SysRoleMenu
	MigrateRoleApp  = []SysRoleApp{
		{RoleID: 1, TenantCode: def.TenantCodeDefault, AppCode: def.AppCore},
	}
	MigrateAppInfo = []SysAppInfo{
		{Code: def.AppCore, Name: "中台"},
		{Code: def.AppAll, Name: "全部"},
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

	//
	//MigrateModuleMenu = []SysModuleMenu{
	//	{ID: 2, AccessCode: def.AppCore, ParentID: 1, Type: 1, Order: 2, Name: "设备管理", Path: "/deviceManagers", Component: "./deviceManagers/index.tsx", Icon: "icon_data_01", Redirect: "", HideInMenu: def.False},
	//	{ID: 3, AccessCode: def.AppCore, ParentID: 1, Type: 1, Order: 9, Name: "系统管理", Path: "/systemManagers", Component: "./systemManagers/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 4, AccessCode: def.AppCore, ParentID: 1, Type: 1, Order: 4, Name: "运维监控", Path: "/operationsMonitorings", Component: "./operationsMonitorings/index.tsx", Icon: "icon_hvac", Redirect: "", HideInMenu: def.False},
	//	{ID: 6, AccessCode: def.AppCore, ParentID: 2, Type: 1, Order: 1, Name: "产品", Path: "/deviceManagers/product/index", Component: "./deviceManagers/product/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 7, AccessCode: def.AppCore, ParentID: 2, Type: 1, Order: 1, Name: "产品详情", Path: "/deviceManagers/product/detail/:id", Component: "./deviceManagers/product/detail/index", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
	//	{ID: 8, AccessCode: def.AppCore, ParentID: 2, Type: 1, Order: 2, Name: "设备", Path: "/deviceManagers/device/index", Component: "./deviceManagers/device/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 9, AccessCode: def.AppCore, ParentID: 2, Type: 1, Order: 2, Name: "设备详情", Path: "/deviceManagers/device/detail/:id/:name/:type", Component: "./deviceManagers/device/detail/index", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
	//	{ID: 10, AccessCode: def.AppCore, ParentID: 3, Type: 1, Order: 1, Name: "用户管理", Path: "/systemManagers/user/index", Component: "./systemManagers/user/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 11, AccessCode: def.AppCore, ParentID: 3, Type: 1, Order: 2, Name: "角色管理", Path: "/systemManagers/role/index", Component: "./systemManagers/role/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 12, AccessCode: def.AppCore, ParentID: 3, Type: 1, Order: 3, Name: "菜单列表", Path: "/systemManagers/menu/index", Component: "./systemManagers/menu/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 13, AccessCode: def.AppCore, ParentID: 4, Type: 1, Order: 1, Name: "固件升级", Path: "/operationsMonitorings/firmwareUpgrade/index", Component: "./operationsMonitorings/firmwareUpgrade/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 15, AccessCode: def.AppCore, ParentID: 4, Type: 1, Order: 3, Name: "资源管理", Path: "/operationsMonitorings/resourceManagement/index", Component: "./operationsMonitorings/resourceManagement/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 16, AccessCode: def.AppCore, ParentID: 4, Type: 1, Order: 4, Name: "远程配置", Path: "/operationsMonitorings/remoteConfiguration/index", Component: "./operationsMonitorings/remoteConfiguration/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 18, AccessCode: def.AppCore, ParentID: 4, Type: 1, Order: 6, Name: "在线调试", Path: "/operationsMonitorings/onlineDebug/index", Component: "./operationsMonitorings/onlineDebug/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 23, AccessCode: def.AppCore, ParentID: 2, Type: 1, Order: 3, Name: "分组", Path: "/deviceManagers/group/index", Component: "./deviceManagers/group/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 24, AccessCode: def.AppCore, ParentID: 2, Type: 1, Order: 3, Name: "分组详情", Path: "/deviceManagers/group/detail/:id", Component: "./deviceManagers/group/detail/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
	//	{ID: 25, AccessCode: def.AppCore, ParentID: 4, Type: 1, Order: 7, Name: "日志服务", Path: "/operationsMonitorings/logService/index", Component: "./operationsMonitorings/logService/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 35, AccessCode: def.AppCore, ParentID: 1, Type: 1, Order: 1, Name: "首页", Path: "/home", Component: "./home/index.tsx", Icon: "icon_dosing", Redirect: "", HideInMenu: def.False},
	//	{ID: 38, AccessCode: def.AppCore, ParentID: 3, Type: 1, Order: 5, Name: "日志管理", Path: "/systemManagers/log", Component: "./systemManagers/log/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 39, AccessCode: def.AppCore, ParentID: 38, Type: 1, Order: 1, Name: "操作日志", Path: "/systemManagers/log/operationLog/index", Component: "./systemManagers/log/operationLog/index.tsx", Icon: "icon_dosing", Redirect: "", HideInMenu: def.False},
	//	{ID: 41, AccessCode: def.AppCore, ParentID: 38, Type: 1, Order: 2, Name: "登录日志", Path: "/systemManagers/log/loginLog/index", Component: "./systemManagers/log/loginLog/index", Icon: "icon_heat", Redirect: "", HideInMenu: def.False},
	//	{ID: 42, AccessCode: def.AppCore, ParentID: 3, Type: 1, Order: 4, Name: "接口管理", Path: "/systemManagers/api/index", Component: "./systemManagers/api/index", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 43, AccessCode: def.AppCore, ParentID: 1, Type: 1, Order: 5, Name: "告警管理", Path: "/alarmManagers", Component: "./alarmManagers/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.False},
	//	{ID: 44, AccessCode: def.AppCore, ParentID: 43, Type: 1, Order: 1, Name: "告警配置", Path: "/alarmManagers/alarmConfiguration/index", Component: "./alarmManagers/alarmConfiguration/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.False},
	//	{ID: 53, AccessCode: def.AppCore, ParentID: 43, Type: 1, Order: 5, Name: "新增告警配置", Path: "/alarmManagers/alarmConfiguration/save", Component: "./alarmManagers/alarmConfiguration/addAlarmConfig/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.True},
	//	{ID: 54, AccessCode: def.AppCore, ParentID: 43, Type: 1, Order: 5, Name: "告警日志", Path: "/alarmManagers/alarmConfiguration/log/detail/:id/:level", Component: "./alarmManagers/alarmLog/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.True},
	//	{ID: 45, AccessCode: def.AppCore, ParentID: 43, Type: 1, Order: 5, Name: "告警记录", Path: "/alarmManagers/alarmConfiguration/log", Component: "./alarmManagers/alarmRecord/index", Icon: "icon_ap", Redirect: "", HideInMenu: def.False},
	//	{ID: 50, AccessCode: def.AppCore, ParentID: 1, Type: 1, Order: 5, Name: "规则引擎", Path: "/ruleEngine", Component: "./ruleEngine/index.tsx", Icon: "icon_dosing", Redirect: "", HideInMenu: def.False},
	//	{ID: 51, AccessCode: def.AppCore, ParentID: 50, Type: 1, Order: 1, Name: "场景联动", Path: "/ruleEngine/scene/index", Component: "./ruleEngine/scene/index.tsx", Icon: "icon_device", Redirect: "", HideInMenu: def.False},
	//	{ID: 60, AccessCode: def.AppCore, ParentID: 3, Type: 2, Order: 1, Name: "内嵌", Path: "/systemManagers/iframe", Component: "https://www.douyu.com", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 61, AccessCode: def.AppCore, ParentID: 3, Type: 3, Order: 1, Name: "外链", Path: "/systemManagers/links", Component: "https://ant.design", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 70, AccessCode: def.AppCore, ParentID: 3, Type: 1, Order: 1, Name: "任务管理", Path: "/systemManagers/timed", Component: "./systemManagers/timed/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 71, AccessCode: def.AppCore, ParentID: 70, Type: 1, Order: 1, Name: "任务组", Path: "/systemManagers/timed/group", Component: "./systemManagers/timed/group/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 72, AccessCode: def.AppCore, ParentID: 70, Type: 1, Order: 1, Name: "任务组详情", Path: "/systemManagers/timed/group/detail/:id", Component: "./systemManagers/timed/group/detail/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
	//	{ID: 73, AccessCode: def.AppCore, ParentID: 70, Type: 1, Order: 1, Name: "任务", Path: "/systemManagers/timed/task", Component: "./systemManagers/timed/task/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.False},
	//	{ID: 74, AccessCode: def.AppCore, ParentID: 70, Type: 1, Order: 1, Name: "任务详情", Path: "/systemManagers/timed/task/detail/:id", Component: "./systemManagers/timed/task/detail/index.tsx", Icon: "icon_system", Redirect: "", HideInMenu: def.True},
	//	//视频服务菜单项
	//	{ID: 63, AccessCode: def.AppCore, ParentID: 1, Type: 1, Order: 2, Name: "视频服务", Path: "/videoManagers", Component: "./videoManagers", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
	//	{ID: 64, AccessCode: def.AppCore, ParentID: 63, Type: 1, Order: 1, Name: "流服务管理", Path: "/videoManagers/vidsrvmgr/index", Component: "./videoManagers/vidsrvmgr/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
	//	{ID: 65, AccessCode: def.AppCore, ParentID: 63, Type: 1, Order: 3, Name: "视频流广场", Path: "/videoManagers/plaza/index", Component: "./videoManagers/plaza/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
	//	{ID: 66, AccessCode: def.AppCore, ParentID: 63, Type: 1, Order: 2, Name: "视频流管理", Path: "/videoManagers/vidstream/index", Component: "./videoManagers/vidstream/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
	//	{ID: 67, AccessCode: def.AppCore, ParentID: 63, Type: 1, Order: 4, Name: "视频回放", Path: "/videoManagers/playback/index", Component: "./videoManagers/playback/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
	//	{ID: 68, AccessCode: def.AppCore, ParentID: 63, Type: 1, Order: 2, Name: "录像计划", Path: "/videoManagers/recordplan/index", Component: "./videoManagers/recordplan/index.tsx", Icon: "icon_heat", Redirect: "", HideInMenu: 2},
	//	{ID: 69, AccessCode: def.AppCore, ParentID: 63, Type: 1, Order: 1, Name: "流服务详细", Path: "/videoManagers/vidsrvmgr/detail/:id", Component: "./videoManagers/vidsrvmgr/detail/index", Icon: "icon_heat", Redirect: "", HideInMenu: 1},
	//	{ID: 75, AccessCode: def.AppCore, ParentID: 63, Type: 1, Order: 1, Name: "视频流详细", Path: "/videoManagers/vidstream/detail/:id", Component: "./videoManagers/vidstream/detail/index", Icon: "icon_heat", Redirect: "", HideInMenu: 1},
	//}
	MigrateAccessInfo = []SysAccessInfo{
		{Name: "升级任务管理task操作权限", Code: "thingsOtaTaskWrite", Group: "ota升级", IsNeedAuth: 2, Desc: ""},
		{Name: "升级任务管理task读权限", Code: "thingsOtaTaskRead", Group: "ota升级", IsNeedAuth: 2, Desc: ""},
		{Name: "升级包管理firmware操作权限", Code: "thingsOtaFirmwareWrite", Group: "ota升级", IsNeedAuth: 2, Desc: ""},
		{Name: "升级包管理firmware读权限", Code: "thingsOtaFirmwareRead", Group: "ota升级", IsNeedAuth: 2, Desc: ""},
		{Name: "升级包管理操作权限", Code: "thingsOtaOtaFirmwareWrite", Group: "ota升级", IsNeedAuth: 2, Desc: ""},
		{Name: "升级包管理读权限", Code: "thingsOtaOtaFirmwareRead", Group: "ota升级", IsNeedAuth: 2, Desc: ""},
		{Name: "升级批次管理操作权限", Code: "thingsOtaJobWrite", Group: "ota升级", IsNeedAuth: 2, Desc: ""},
		{Name: "升级批次管理读权限", Code: "thingsOtaJobRead", Group: "ota升级", IsNeedAuth: 2, Desc: ""},
		{Name: "产品品类操作权限", Code: "thingsProductCategoryWrite", Group: "产品", IsNeedAuth: 2, Desc: ""},
		{Name: "产品品类读权限", Code: "thingsProductCategoryRead", Group: "产品", IsNeedAuth: 2, Desc: ""},
		{Name: "产品管理操作权限", Code: "thingsProductInfoWrite", Group: "产品", IsNeedAuth: 2, Desc: ""},
		{Name: "产品管理读权限", Code: "thingsProductInfoRead", Group: "产品", IsNeedAuth: 2, Desc: ""},
		{Name: "物模型操作权限", Code: "thingsProductSchemaWrite", Group: "产品", IsNeedAuth: 2, Desc: ""},
		{Name: "物模型读权限", Code: "thingsProductSchemaRead", Group: "产品", IsNeedAuth: 2, Desc: ""},
		{Name: "自定义操作权限", Code: "thingsProductCustomWrite", Group: "产品", IsNeedAuth: 2, Desc: ""},
		{Name: "自定义读权限", Code: "thingsProductCustomRead", Group: "产品", IsNeedAuth: 2, Desc: ""},
		{Name: "远程配置操作权限", Code: "thingsProductRemoteConfigWrite", Group: "产品", IsNeedAuth: 2, Desc: ""},
		{Name: "远程配置读权限", Code: "thingsProductRemoteConfigRead", Group: "产品", IsNeedAuth: 2, Desc: ""},
		{Name: "任务操作权限", Code: "systemJobTaskWrite", Group: "任务管理", IsNeedAuth: 2, Desc: ""},
		{Name: "任务组读权限", Code: "systemJobTaskRead", Group: "任务管理", IsNeedAuth: 2, Desc: ""},
		{Name: "区域管理操作权限", Code: "systemAreaInfoWrite", Group: "区域管理", IsNeedAuth: 2, Desc: ""},
		{Name: "区域管理读权限", Code: "systemAreaInfoRead", Group: "区域管理", IsNeedAuth: 2, Desc: ""},
		{Name: "协议管理操作权限", Code: "thingsProtocolInfoWrite", Group: "协议管理", IsNeedAuth: 2, Desc: ""},
		{Name: "协议管理读权限", Code: "thingsProtocolInfoRead", Group: "协议管理", IsNeedAuth: 2, Desc: ""},
		{Name: "字典信息操作权限", Code: "systemDictInfoWrite", Group: "字典管理", IsNeedAuth: 2, Desc: ""},
		{Name: "字典信息读权限", Code: "systemDictInfoRead", Group: "字典管理", IsNeedAuth: 2, Desc: ""},
		{Name: "字典详情操作权限", Code: "systemDictDetailWrite", Group: "字典管理", IsNeedAuth: 2, Desc: ""},
		{Name: "字典详情读权限", Code: "systemDictDetailRead", Group: "字典管理", IsNeedAuth: 2, Desc: ""},
		{Name: "应用管理操作权限", Code: "systemAppInfoWrite", Group: "应用管理", IsNeedAuth: 2, Desc: ""},
		{Name: "应用管理读权限", Code: "systemAppInfoRead", Group: "应用管理", IsNeedAuth: 2, Desc: ""},
		{Name: "模块操作权限", Code: "systemAppModuleWrite", Group: "应用管理", IsNeedAuth: 2, Desc: ""},
		{Name: "模块读权限", Code: "systemAppModuleRead", Group: "应用管理", IsNeedAuth: 2, Desc: ""},
		{Name: "授权信息操作权限", Code: "systemAccessInfoWrite", Group: "授权管理", IsNeedAuth: 2, Desc: ""},
		{Name: "授权信息读权限", Code: "systemAccessInfoRead", Group: "授权管理", IsNeedAuth: 2, Desc: ""},
		{Name: "接口操作权限", Code: "systemAccessApiWrite", Group: "授权管理", IsNeedAuth: 2, Desc: ""},
		{Name: "接口读权限", Code: "systemAccessApiRead", Group: "授权管理", IsNeedAuth: 2, Desc: ""},
		{Name: "区域操作权限", Code: "systemDataAreaWrite", Group: "数据管理", IsNeedAuth: 2, Desc: ""},
		{Name: "区域用户授权操作权限", Code: "systemUserAreaWrite", Group: "数据管理", IsNeedAuth: 2, Desc: ""},
		{Name: "区域用户授权读权限", Code: "systemUserAreaRead", Group: "数据管理", IsNeedAuth: 2, Desc: ""},
		{Name: "区域读权限", Code: "systemDataAreaRead", Group: "数据管理", IsNeedAuth: 2, Desc: ""},
		{Name: "项目操作权限", Code: "systemUserAuthWrite", Group: "数据管理", IsNeedAuth: 2, Desc: ""},
		{Name: "项目读权限", Code: "systemUserAuthRead", Group: "数据管理", IsNeedAuth: 2, Desc: ""},
		{Name: "日志管理读权限", Code: "systemLogLoginRead", Group: "日志管理", IsNeedAuth: 2, Desc: ""},
		{Name: "日志管理读权限", Code: "systemLogOperRead", Group: "日志管理", IsNeedAuth: 2, Desc: ""},
		{Name: "模块操作权限", Code: "systemModuleInfoWrite", Group: "模块管理", IsNeedAuth: 2, Desc: ""},
		{Name: "模块读权限", Code: "systemModuleInfoRead", Group: "模块管理", IsNeedAuth: 2, Desc: ""},
		{Name: "菜单操作权限", Code: "systemModuleMenuWrite", Group: "模块管理", IsNeedAuth: 2, Desc: ""},
		{Name: "菜单读权限", Code: "systemModuleMenuRead", Group: "模块管理", IsNeedAuth: 2, Desc: ""},
		{Name: "通用物模型操作权限", Code: "thingsSchemaCommonWrite", Group: "物模型管理", IsNeedAuth: 2, Desc: ""},
		{Name: "通用物模型读权限", Code: "thingsSchemaCommonRead", Group: "物模型管理", IsNeedAuth: 2, Desc: ""},
		{Name: "设备分享读权限", Code: "thingsUserDeviceRead", Group: "用户", IsNeedAuth: 2, Desc: ""},
		{Name: "设备收藏操作权限", Code: "thingsUserDeviceWrite", Group: "用户", IsNeedAuth: 2, Desc: ""},
		{Name: "用户管理操作权限", Code: "systemUserRoleWrite", Group: "用户管理", IsNeedAuth: 2, Desc: ""},
		{Name: "用户管理操作权限", Code: "systemUserInfoWrite", Group: "用户管理", IsNeedAuth: 2, Desc: ""},
		{Name: "用户管理读权限", Code: "systemUserRoleRead", Group: "用户管理", IsNeedAuth: 2, Desc: ""},
		{Name: "用户管理读权限", Code: "systemUserInfoRead", Group: "用户管理", IsNeedAuth: 2, Desc: ""},
		{Name: "自己操作权限", Code: "systemUserSelfWrite", Group: "用户管理", IsNeedAuth: 2, Desc: ""},
		{Name: "自己读权限", Code: "systemUserSelfRead", Group: "用户管理", IsNeedAuth: 2, Desc: ""},
		{Name: "应用管理模块管理操作权限", Code: "systemTenantAppWrite", Group: "租户管理", IsNeedAuth: 2, Desc: ""},
		{Name: "应用管理模块管理读权限", Code: "systemTenantAppRead", Group: "租户管理", IsNeedAuth: 2, Desc: ""},
		{Name: "授权管理操作权限", Code: "systemTenantAccessWrite", Group: "租户管理", IsNeedAuth: 2, Desc: ""},
		{Name: "授权管理读权限", Code: "systemTenantAccessRead", Group: "租户管理", IsNeedAuth: 2, Desc: ""},
		{Name: "租户管理操作权限", Code: "systemTenantInfoWrite", Group: "租户管理", IsNeedAuth: 2, Desc: ""},
		{Name: "租户管理读权限", Code: "systemTenantInfoRead", Group: "租户管理", IsNeedAuth: 2, Desc: ""},
		{Name: "告警中心告警日志读权限", Code: "thingsRuleAlarmRead", Group: "规则引擎", IsNeedAuth: 2, Desc: ""},
		{Name: "告警中心场景联动关联操作权限", Code: "thingsRuleAlarmWrite", Group: "规则引擎", IsNeedAuth: 2, Desc: ""},
		{Name: "场景联动操作权限", Code: "thingsRuleSceneWrite", Group: "规则引擎", IsNeedAuth: 2, Desc: ""},
		{Name: "场景联动读权限", Code: "thingsRuleSceneRead", Group: "规则引擎", IsNeedAuth: 2, Desc: ""},
		{Name: "规则编排流操作权限", Code: "thingsRuleFlowWrite", Group: "规则引擎", IsNeedAuth: 2, Desc: ""},
		{Name: "规则编排流读权限", Code: "thingsRuleFlowRead", Group: "规则引擎", IsNeedAuth: 2, Desc: ""},
		{Name: "设备定时操作权限", Code: "thingsRuleDeviceTimerWrite", Group: "规则引擎", IsNeedAuth: 2, Desc: ""},
		{Name: "设备定时读权限", Code: "thingsRuleDeviceTimerRead", Group: "规则引擎", IsNeedAuth: 2, Desc: ""},
		{Name: "国标协议服务操作权限", Code: "thingsVidmgrGbsipWrite", Group: "视频服务", IsNeedAuth: 2, Desc: ""},
		{Name: "流服务交互操作权限", Code: "thingsVidmgrCtrlWrite", Group: "视频服务", IsNeedAuth: 2, Desc: ""},
		{Name: "流服务器管理操作权限", Code: "thingsVidmgrInfoWrite", Group: "视频服务", IsNeedAuth: 2, Desc: ""},
		{Name: "流服务器管理读权限", Code: "thingsVidmgrInfoRead", Group: "视频服务", IsNeedAuth: 2, Desc: ""},
		{Name: "视频流管理操作权限", Code: "thingsVidmgrStreamWrite", Group: "视频服务", IsNeedAuth: 2, Desc: ""},
		{Name: "视频流管理读权限", Code: "thingsVidmgrStreamRead", Group: "视频服务", IsNeedAuth: 2, Desc: ""},
		{Name: "应用操作权限", Code: "systemRoleAppWrite", Group: "角色管理", IsNeedAuth: 2, Desc: ""},
		{Name: "应用读权限", Code: "systemRoleAppRead", Group: "角色管理", IsNeedAuth: 2, Desc: ""},
		{Name: "授权操作权限", Code: "systemRoleAccessWrite", Group: "角色管理", IsNeedAuth: 2, Desc: ""},
		{Name: "授权读权限", Code: "systemRoleAccessRead", Group: "角色管理", IsNeedAuth: 2, Desc: ""},
		{Name: "模块操作权限", Code: "systemRoleModuleWrite", Group: "角色管理", IsNeedAuth: 2, Desc: ""},
		{Name: "模块读权限", Code: "systemRoleModuleRead", Group: "角色管理", IsNeedAuth: 2, Desc: ""},
		{Name: "菜单操作权限", Code: "systemRoleMenuWrite", Group: "角色管理", IsNeedAuth: 2, Desc: ""},
		{Name: "菜单读权限", Code: "systemRoleMenuRead", Group: "角色管理", IsNeedAuth: 2, Desc: ""},
		{Name: "角色管理操作权限", Code: "systemRoleInfoWrite", Group: "角色管理", IsNeedAuth: 2, Desc: ""},
		{Name: "角色管理读权限", Code: "systemRoleInfoRead", Group: "角色管理", IsNeedAuth: 2, Desc: ""},
		{Name: "网关子设备管理操作权限", Code: "thingsDeviceGatewayWrite", Group: "设备", IsNeedAuth: 2, Desc: ""},
		{Name: "网关子设备管理读权限", Code: "thingsDeviceGatewayRead", Group: "设备", IsNeedAuth: 2, Desc: ""},
		{Name: "设备交互操作权限", Code: "thingsDeviceInteractWrite", Group: "设备", IsNeedAuth: 2, Desc: ""},
		{Name: "设备交互读权限", Code: "thingsDeviceInteractRead", Group: "设备", IsNeedAuth: 2, Desc: ""},
		{Name: "设备消息读权限", Code: "thingsDeviceMsgRead", Group: "设备", IsNeedAuth: 2, Desc: ""},
		{Name: "设备管理操作权限", Code: "thingsDeviceInfoWrite", Group: "设备", IsNeedAuth: 2, Desc: ""},
		{Name: "设备管理读权限", Code: "thingsDeviceInfoRead", Group: "设备", IsNeedAuth: 2, Desc: ""},
		{Name: "设备鉴权操作权限", Code: "thingsDeviceAuthWrite", Group: "设备", IsNeedAuth: 2, Desc: ""},
		{Name: "设备分组操作权限", Code: "thingsGroupInfoWrite", Group: "设备分组", IsNeedAuth: 2, Desc: ""},
		{Name: "设备分组操作权限", Code: "thingsGroupDeviceWrite", Group: "设备分组", IsNeedAuth: 2, Desc: ""},
		{Name: "设备分组读权限", Code: "thingsGroupInfoRead", Group: "设备分组", IsNeedAuth: 2, Desc: ""},
		{Name: "设备分组读权限", Code: "thingsGroupDeviceRead", Group: "设备分组", IsNeedAuth: 2, Desc: ""},
		{Name: "工单操作权限", Code: "thingsOpsWorkOrderWrite", Group: "运营维护", IsNeedAuth: 2, Desc: ""},
		{Name: "工单读权限", Code: "thingsOpsWorkOrderRead", Group: "运营维护", IsNeedAuth: 2, Desc: ""},
		{Name: "通用功能操作权限", Code: "systemCommonConfigWrite", Group: "通用功能", IsNeedAuth: 2, Desc: ""},
		{Name: "通用功能操作权限", Code: "systemCommonUploadUrlWrite", Group: "通用功能", IsNeedAuth: 2, Desc: ""},
		{Name: "通用功能操作权限", Code: "systemCommonUploadFileWrite", Group: "通用功能", IsNeedAuth: 2, Desc: ""},
		{Name: "通用功能读权限", Code: "systemCommonWeatherRead", Group: "通用功能", IsNeedAuth: 2, Desc: ""},
		{Name: "项目管理操作权限", Code: "systemProjectInfoWrite", Group: "项目管理", IsNeedAuth: 2, Desc: ""},
		{Name: "项目管理读权限", Code: "systemProjectInfoRead", Group: "项目管理", IsNeedAuth: 2, Desc: ""},
	}
	MigrateApiInfo = []SysApiInfo{
		{AccessCode: "systemAccessApiRead", IsAuthTenant: 1, Route: "/api/v1/system/access/api/index", Method: http.MethodPost, Name: "获取接口列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemAccessApiWrite", IsAuthTenant: 1, Route: "/api/v1/system/access/api/create", Method: http.MethodPost, Name: "添加接口", BusinessType: 1, Desc: ``},
		{AccessCode: "systemAccessApiWrite", IsAuthTenant: 1, Route: "/api/v1/system/access/api/delete", Method: http.MethodPost, Name: "删除接口", BusinessType: 3, Desc: ``},
		{AccessCode: "systemAccessApiWrite", IsAuthTenant: 1, Route: "/api/v1/system/access/api/update", Method: http.MethodPost, Name: "更新接口", BusinessType: 2, Desc: ``},
		{AccessCode: "systemAccessInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/access/info/index", Method: http.MethodPost, Name: "获取授权列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemAccessInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/access/info/create", Method: http.MethodPost, Name: "添加授权", BusinessType: 1, Desc: ``},
		{AccessCode: "systemAccessInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/access/info/delete", Method: http.MethodPost, Name: "删除授权", BusinessType: 3, Desc: ``},
		{AccessCode: "systemAccessInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/access/info/tree", Method: http.MethodPost, Name: "获取授权树", BusinessType: 5, Desc: ``},
		{AccessCode: "systemAccessInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/access/info/update", Method: http.MethodPost, Name: "更新授权", BusinessType: 2, Desc: ``},
		{AccessCode: "systemAppInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/app/info/index", Method: http.MethodPost, Name: "获取应用列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemAppInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/app/info/create", Method: http.MethodPost, Name: "添加应用", BusinessType: 1, Desc: ``},
		{AccessCode: "systemAppInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/app/info/delete", Method: http.MethodPost, Name: "删除应用", BusinessType: 3, Desc: ``},
		{AccessCode: "systemAppInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/app/info/update", Method: http.MethodPost, Name: "更新应用", BusinessType: 2, Desc: ``},
		{AccessCode: "systemAppModuleRead", IsAuthTenant: 1, Route: "/api/v1/system/app/module/index", Method: http.MethodPost, Name: "获取应用绑定的模块列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemAppModuleWrite", IsAuthTenant: 1, Route: "/api/v1/system/app/module/multi-update", Method: http.MethodPost, Name: "批量更新应用绑定的模块", BusinessType: 2, Desc: ``},
		{AccessCode: "systemAreaInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/area/info/index", Method: http.MethodPost, Name: "获取项目区域列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemAreaInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/area/info/read", Method: http.MethodPost, Name: "获取项目区域详情", BusinessType: 4, Desc: ``},
		{AccessCode: "systemAreaInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/area/info/create", Method: http.MethodPost, Name: "新增项目区域", BusinessType: 1, Desc: ``},
		{AccessCode: "systemAreaInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/area/info/delete", Method: http.MethodPost, Name: "删除项目区域", BusinessType: 3, Desc: ``},
		{AccessCode: "systemAreaInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/area/info/update", Method: http.MethodPost, Name: "更新项目区域", BusinessType: 2, Desc: ``},
		{AccessCode: "systemCommonConfigWrite", IsAuthTenant: 1, Route: "/api/v1/system/common/config", Method: http.MethodPost, Name: "获取系统配置", BusinessType: 5, Desc: ``},
		{AccessCode: "systemCommonUploadFileWrite", IsAuthTenant: 1, Route: "/api/v1/system/common/upload-file", Method: http.MethodPost, Name: "文件直传接口", BusinessType: 5, Desc: ``},
		{AccessCode: "systemCommonUploadUrlWrite", IsAuthTenant: 1, Route: "/api/v1/system/common/upload-url/create", Method: http.MethodPost, Name: "获取文件上传地址", BusinessType: 1, Desc: `接口返回signed-url ,前端获取到该url后，往该url put上传文件`},
		{AccessCode: "systemCommonWeatherRead", IsAuthTenant: 1, Route: "/api/v1/system/common/weather/read", Method: http.MethodPost, Name: "获取天气情况", BusinessType: 4, Desc: `参考:
https://dev.qweather.com/docs/api/weather/weather-now/
https://dev.qweather.com/docs/api/air/air-now/`},
		{AccessCode: "systemDataAreaRead", IsAuthTenant: 1, Route: "/api/v1/system/data/area/index", Method: http.MethodPost, Name: "获取区域权限列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemDataAreaWrite", IsAuthTenant: 1, Route: "/api/v1/system/data/area/multi-update", Method: http.MethodPost, Name: "授权区域权限（内部会先全删后重加)", BusinessType: 2, Desc: ``},
		{AccessCode: "systemDictDetailRead", IsAuthTenant: 1, Route: "/api/v1/system/dict/detail/index", Method: http.MethodPost, Name: "获取字典详情列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemDictDetailWrite", IsAuthTenant: 1, Route: "/api/v1/system/dict/detail/create", Method: http.MethodPost, Name: "添加字典详情", BusinessType: 1, Desc: ``},
		{AccessCode: "systemDictDetailWrite", IsAuthTenant: 1, Route: "/api/v1/system/dict/detail/delete", Method: http.MethodPost, Name: "删除字典详情", BusinessType: 3, Desc: ``},
		{AccessCode: "systemDictDetailWrite", IsAuthTenant: 1, Route: "/api/v1/system/dict/detail/update", Method: http.MethodPost, Name: "更新字典详情", BusinessType: 2, Desc: ``},
		{AccessCode: "systemDictInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/dict/info/index", Method: http.MethodPost, Name: "获取字典列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemDictInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/dict/info/read", Method: http.MethodPost, Name: "获取字典", BusinessType: 4, Desc: ``},
		{AccessCode: "systemDictInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/dict/info/create", Method: http.MethodPost, Name: "添加字典", BusinessType: 1, Desc: ``},
		{AccessCode: "systemDictInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/dict/info/delete", Method: http.MethodPost, Name: "删除字典", BusinessType: 3, Desc: ``},
		{AccessCode: "systemDictInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/dict/info/update", Method: http.MethodPost, Name: "更新字典", BusinessType: 2, Desc: ``},
		{AccessCode: "systemJobTaskRead", IsAuthTenant: 1, Route: "/api/v1/system/job/task/group/index", Method: http.MethodPost, Name: "获取任务组列表", BusinessType: 4, Desc: `database这个配置项只有sql和script有`},
		{AccessCode: "systemJobTaskRead", IsAuthTenant: 1, Route: "/api/v1/system/job/task/group/read", Method: http.MethodPost, Name: "获取任务组详情", BusinessType: 4, Desc: ``},
		{AccessCode: "systemJobTaskRead", IsAuthTenant: 1, Route: "/api/v1/system/job/task/info/index", Method: http.MethodPost, Name: "获取任务列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemJobTaskRead", IsAuthTenant: 1, Route: "/api/v1/system/job/task/info/read", Method: http.MethodPost, Name: "获取任务详情", BusinessType: 4, Desc: ``},
		{AccessCode: "systemJobTaskRead", IsAuthTenant: 1, Route: "/api/v1/system/job/task/log/index", Method: http.MethodPost, Name: "获取任务日志列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemJobTaskWrite", IsAuthTenant: 1, Route: "/api/v1/system/job/task/cancel", Method: http.MethodPost, Name: "取消任务", BusinessType: 5, Desc: ``},
		{AccessCode: "systemJobTaskWrite", IsAuthTenant: 1, Route: "/api/v1/system/job/task/group/create", Method: http.MethodPost, Name: "创建任务组", BusinessType: 1, Desc: ``},
		{AccessCode: "systemJobTaskWrite", IsAuthTenant: 1, Route: "/api/v1/system/job/task/group/delete", Method: http.MethodPost, Name: "删除任务组", BusinessType: 3, Desc: ``},
		{AccessCode: "systemJobTaskWrite", IsAuthTenant: 1, Route: "/api/v1/system/job/task/group/update", Method: http.MethodPost, Name: "更新任务组", BusinessType: 2, Desc: ``},
		{AccessCode: "systemJobTaskWrite", IsAuthTenant: 1, Route: "/api/v1/system/job/task/info/create", Method: http.MethodPost, Name: "创建任务", BusinessType: 1, Desc: ""},
		{AccessCode: "systemJobTaskWrite", IsAuthTenant: 1, Route: "/api/v1/system/job/task/info/delete", Method: http.MethodPost, Name: "删除任务", BusinessType: 3, Desc: ``},
		{AccessCode: "systemJobTaskWrite", IsAuthTenant: 1, Route: "/api/v1/system/job/task/info/update", Method: http.MethodPost, Name: "更新任务", BusinessType: 2, Desc: ``},
		{AccessCode: "systemJobTaskWrite", IsAuthTenant: 1, Route: "/api/v1/system/job/task/send", Method: http.MethodPost, Name: "执行任务", BusinessType: 5, Desc: ``},
		{AccessCode: "systemLogLoginRead", IsAuthTenant: 1, Route: "/api/v1/system/log/login/index", Method: http.MethodPost, Name: "获取登录日志列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemLogOperRead", IsAuthTenant: 1, Route: "/api/v1/system/log/oper/index", Method: http.MethodPost, Name: "获取操作日志列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemModuleInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/module/info/index", Method: http.MethodPost, Name: "获取模块列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemModuleInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/module/info/create", Method: http.MethodPost, Name: "添加模块", BusinessType: 1, Desc: ``},
		{AccessCode: "systemModuleInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/module/info/delete", Method: http.MethodPost, Name: "删除模块", BusinessType: 3, Desc: ``},
		{AccessCode: "systemModuleInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/module/info/update", Method: http.MethodPost, Name: "更新模块", BusinessType: 2, Desc: ``},
		{AccessCode: "systemModuleMenuRead", IsAuthTenant: 1, Route: "/api/v1/system/module/menu/index", Method: http.MethodPost, Name: "获取菜单列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemModuleMenuWrite", IsAuthTenant: 1, Route: "/api/v1/system/module/menu/create", Method: http.MethodPost, Name: "添加菜单", BusinessType: 1, Desc: ``},
		{AccessCode: "systemModuleMenuWrite", IsAuthTenant: 1, Route: "/api/v1/system/module/menu/delete", Method: http.MethodPost, Name: "删除菜单", BusinessType: 3, Desc: ``},
		{AccessCode: "systemModuleMenuWrite", IsAuthTenant: 1, Route: "/api/v1/system/module/menu/update", Method: http.MethodPost, Name: "更新菜单", BusinessType: 2, Desc: ``},
		{AccessCode: "systemProjectInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/project/info/index", Method: http.MethodPost, Name: "获取项目列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemProjectInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/project/info/read", Method: http.MethodPost, Name: "获取项目信息", BusinessType: 4, Desc: ``},
		{AccessCode: "systemProjectInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/project/info/create", Method: http.MethodPost, Name: "新增项目", BusinessType: 1, Desc: ``},
		{AccessCode: "systemProjectInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/project/info/delete", Method: http.MethodPost, Name: "删除项目", BusinessType: 3, Desc: ``},
		{AccessCode: "systemProjectInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/project/info/update", Method: http.MethodPost, Name: "更新项目", BusinessType: 2, Desc: ``},
		{AccessCode: "systemRoleAccessRead", IsAuthTenant: 1, Route: "/api/v1/system/role/access/index", Method: http.MethodPost, Name: "获取角色对应授权列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemRoleAccessWrite", IsAuthTenant: 1, Route: "/api/v1/system/role/access/multi-update", Method: http.MethodPost, Name: "更新角色对应授权列表", BusinessType: 2, Desc: ``},
		{AccessCode: "systemRoleAppRead", IsAuthTenant: 1, Route: "/api/v1/system/role/app/index", Method: http.MethodPost, Name: "获取角色对应应用列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemRoleAppWrite", IsAuthTenant: 1, Route: "/api/v1/system/role/app/multi-update", Method: http.MethodPost, Name: "更新角色对应应用列表", BusinessType: 2, Desc: ``},
		{AccessCode: "systemRoleInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/role/info/index", Method: http.MethodPost, Name: "获取角色列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemRoleInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/role/info/create", Method: http.MethodPost, Name: "添加角色", BusinessType: 1, Desc: ``},
		{AccessCode: "systemRoleInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/role/info/delete", Method: http.MethodPost, Name: "删除角色", BusinessType: 3, Desc: ``},
		{AccessCode: "systemRoleInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/role/info/update", Method: http.MethodPost, Name: "更新角色", BusinessType: 2, Desc: ``},
		{AccessCode: "systemRoleMenuRead", IsAuthTenant: 1, Route: "/api/v1/system/role/menu/index", Method: http.MethodPost, Name: "获取角色对应菜单列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemRoleMenuWrite", IsAuthTenant: 1, Route: "/api/v1/system/role/menu/multi-update", Method: http.MethodPost, Name: "更新角色对应菜单列表", BusinessType: 2, Desc: ``},
		{AccessCode: "systemRoleModuleRead", IsAuthTenant: 1, Route: "/api/v1/system/role/module/index", Method: http.MethodPost, Name: "获取角色对应模块列表 ", BusinessType: 4, Desc: ``},
		{AccessCode: "systemRoleModuleWrite", IsAuthTenant: 1, Route: "/api/v1/system/role/module/multi-update", Method: http.MethodPost, Name: "更新角色对应模块列表 ", BusinessType: 2, Desc: ``},
		{AccessCode: "systemTenantAccessRead", IsAuthTenant: 1, Route: "/api/v1/system/tenant/access/info/index", Method: http.MethodPost, Name: "获取租户授权列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemTenantAccessWrite", IsAuthTenant: 1, Route: "/api/v1/system/tenant/access/info/multi-update", Method: http.MethodPost, Name: "批量更新租户授权", BusinessType: 2, Desc: ``},
		{AccessCode: "systemTenantAccessWrite", IsAuthTenant: 1, Route: "/api/v1/system/tenant/access/info/tree", Method: http.MethodPost, Name: "获取授权树", BusinessType: 5, Desc: ``},
		{AccessCode: "systemTenantAppRead", IsAuthTenant: 1, Route: "/api/v1/system/tenant/app/index", Method: http.MethodPost, Name: "获取租户下绑定的应用列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemTenantAppRead", IsAuthTenant: 1, Route: "/api/v1/system/tenant/app/menu/index", Method: http.MethodPost, Name: "获取菜单列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemTenantAppRead", IsAuthTenant: 1, Route: "/api/v1/system/tenant/app/module/index", Method: http.MethodPost, Name: "获取模块绑定列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemTenantAppWrite", IsAuthTenant: 1, Route: "/api/v1/system/tenant/app/create", Method: http.MethodPost, Name: "新增租户下的应用绑定", BusinessType: 1, Desc: ``},
		{AccessCode: "systemTenantAppWrite", IsAuthTenant: 1, Route: "/api/v1/system/tenant/app/delete", Method: http.MethodPost, Name: "删除租户下绑定的应用", BusinessType: 3, Desc: ``},
		{AccessCode: "systemTenantAppWrite", IsAuthTenant: 1, Route: "/api/v1/system/tenant/app/menu/create", Method: http.MethodPost, Name: "添加菜单", BusinessType: 1, Desc: ``},
		{AccessCode: "systemTenantAppWrite", IsAuthTenant: 1, Route: "/api/v1/system/tenant/app/menu/delete", Method: http.MethodPost, Name: "删除菜单", BusinessType: 3, Desc: ``},
		{AccessCode: "systemTenantAppWrite", IsAuthTenant: 1, Route: "/api/v1/system/tenant/app/menu/update", Method: http.MethodPost, Name: "更新菜单", BusinessType: 2, Desc: ``},
		{AccessCode: "systemTenantAppWrite", IsAuthTenant: 1, Route: "/api/v1/system/tenant/app/module/create", Method: http.MethodPost, Name: "添加模块绑定", BusinessType: 1, Desc: ``},
		{AccessCode: "systemTenantAppWrite", IsAuthTenant: 1, Route: "/api/v1/system/tenant/app/module/delete", Method: http.MethodPost, Name: "删除绑定模块", BusinessType: 3, Desc: ``},
		{AccessCode: "systemTenantAppWrite", IsAuthTenant: 1, Route: "/api/v1/system/tenant/app/module/multi-create", Method: http.MethodPost, Name: "批量添加模块绑定", BusinessType: 1, Desc: ``},
		{AccessCode: "systemTenantInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/tenant/info/index", Method: http.MethodPost, Name: "获取租户列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemTenantInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/tenant/info/read", Method: http.MethodPost, Name: "获取租户详情", BusinessType: 4, Desc: ``},
		{AccessCode: "systemTenantInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/tenant/info/create", Method: http.MethodPost, Name: "新增租户", BusinessType: 1, Desc: ``},
		{AccessCode: "systemTenantInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/tenant/info/delete", Method: http.MethodPost, Name: "删除租户", BusinessType: 3, Desc: ``},
		{AccessCode: "systemTenantInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/tenant/info/update", Method: http.MethodPost, Name: "更新租户", BusinessType: 2, Desc: ``},
		{AccessCode: "systemUserAreaRead", IsAuthTenant: 1, Route: "/api/v1/system/user/area/apply/index", Method: http.MethodPost, Name: "获取用户申请区域权限列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemUserAreaWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/area/apply/deal", Method: http.MethodPost, Name: "处理用户申请区域权限", BusinessType: 5, Desc: ``},
		{AccessCode: "systemUserAuthRead", IsAuthTenant: 1, Route: "/api/v1/system/user/auth/project/index", Method: http.MethodPost, Name: "获取用户项目权限列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemUserAuthWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/auth/project/multi-update", Method: http.MethodPost, Name: "授权用户项目权限（内部会先全删后重加）", BusinessType: 2, Desc: ``},
		{AccessCode: "systemUserInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/user/info/index", Method: http.MethodPost, Name: "获取用户信息列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemUserInfoRead", IsAuthTenant: 1, Route: "/api/v1/system/user/info/read", Method: http.MethodPost, Name: "获取用户信息", BusinessType: 4, Desc: ``},
		{AccessCode: "systemUserInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/info/create", Method: http.MethodPost, Name: "创建用户信息", BusinessType: 1, Desc: ``},
		{AccessCode: "systemUserInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/info/delete", Method: http.MethodPost, Name: "删除用户", BusinessType: 3, Desc: ``},
		{AccessCode: "systemUserInfoWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/info/update", Method: http.MethodPost, Name: "更新用户信息", BusinessType: 2, Desc: ``},
		{AccessCode: "systemUserRoleRead", IsAuthTenant: 1, Route: "/api/v1/system/user/role/index", Method: http.MethodPost, Name: "获取用户角色信息列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemUserRoleWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/role/multi-update", Method: http.MethodPost, Name: "用户角色信息批量更新", BusinessType: 2, Desc: ``},
		{AccessCode: "systemUserSelfRead", IsAuthTenant: 1, Route: "/api/v1/system/user/self/app/index", Method: http.MethodPost, Name: "获取用户应用列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemUserSelfRead", IsAuthTenant: 1, Route: "/api/v1/system/user/self/menu/index", Method: http.MethodPost, Name: "获取用户菜单列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemUserSelfRead", IsAuthTenant: 1, Route: "/api/v1/system/user/self/module/index", Method: http.MethodPost, Name: "获取用户模块列表", BusinessType: 4, Desc: ``},
		{AccessCode: "systemUserSelfRead", IsAuthTenant: 1, Route: "/api/v1/system/user/self/read", Method: http.MethodPost, Name: "用户获取自己的用户信息", BusinessType: 4, Desc: ``},
		{AccessCode: "systemUserSelfWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/self/access/tree", Method: http.MethodPost, Name: "获取用户授权树", BusinessType: 5, Desc: ``},
		{AccessCode: "systemUserSelfWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/self/area/apply/create", Method: http.MethodPost, Name: "申请用户区域权限", BusinessType: 1, Desc: ``},
		{AccessCode: "systemUserSelfWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/self/cancel", Method: http.MethodPost, Name: "用户注销账号", BusinessType: 5, Desc: `注册接口`},
		{AccessCode: "systemUserSelfWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/self/captcha", Method: http.MethodPost, Name: "获取验证码", BusinessType: 5, Desc: ``},
		{AccessCode: "systemUserSelfWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/self/change-pwd", Method: http.MethodPost, Name: "用户修改密码", BusinessType: 5, Desc: `注册接口`},
		{AccessCode: "systemUserSelfWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/self/forget-pwd", Method: http.MethodPost, Name: "用户忘记密码", BusinessType: 5, Desc: `注册接口`},
		{AccessCode: "systemUserSelfWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/self/login", Method: http.MethodPost, Name: "登录", BusinessType: 5, Desc: ``},
		{AccessCode: "systemUserSelfWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/self/register", Method: http.MethodPost, Name: "用户注册", BusinessType: 5, Desc: `注册接口`},
		{AccessCode: "systemUserSelfWrite", IsAuthTenant: 1, Route: "/api/v1/system/user/self/update", Method: http.MethodPost, Name: "更新用户基本数据", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsDeviceAuthWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/auth/access", Method: http.MethodPost, Name: "设备操作认证", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsDeviceAuthWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/auth/login", Method: http.MethodPost, Name: "设备登录认证", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsDeviceAuthWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/auth/register", Method: http.MethodPost, Name: "设备动态注册", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsDeviceAuthWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/auth/root-check", Method: http.MethodPost, Name: "鉴定mqtt账号root权限", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsDeviceGatewayRead", IsAuthTenant: 1, Route: "/api/v1/things/device/gateway/index", Method: http.MethodPost, Name: "获取子设备列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsDeviceGatewayWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/gateway/multi-create", Method: http.MethodPost, Name: "批量添加网关子设备", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsDeviceGatewayWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/gateway/multi-delete", Method: http.MethodPost, Name: "批量解绑网关子设备", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsDeviceInfoRead", IsAuthTenant: 1, Route: "/api/v1/things/device/info/index", Method: http.MethodPost, Name: "获取设备列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsDeviceInfoRead", IsAuthTenant: 1, Route: "/api/v1/things/device/info/read", Method: http.MethodPost, Name: "获取设备详情", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsDeviceInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/info/count", Method: http.MethodPost, Name: "设备统计详情", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsDeviceInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/info/create", Method: http.MethodPost, Name: "新增设备", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsDeviceInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/info/delete", Method: http.MethodPost, Name: "删除设备", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsDeviceInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/info/multi-import", Method: http.MethodPost, Name: "批量导入设备", BusinessType: 5, Desc: `#### 前端处理逻辑建议：
- UI text 显示 导入成功 设备数：total - len(errdata)
- UI text 显示 导入失败 设备数：len(errdata)
- UI table 显示 导入失败设备清单明细`},
		{AccessCode: "thingsDeviceInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/info/multi-update", Method: http.MethodPost, Name: "批量更新设备", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsDeviceInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/info/update", Method: http.MethodPost, Name: "更新设备", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsDeviceInteractRead", IsAuthTenant: 1, Route: "/api/v1/things/device/interact/action-read", Method: http.MethodPost, Name: "获取调用设备行为的结果", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsDeviceInteractRead", IsAuthTenant: 1, Route: "/api/v1/things/device/interact/property-control-read", Method: http.MethodPost, Name: "获取调用设备属性的结果", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsDeviceInteractWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/interact/action-send", Method: http.MethodPost, Name: "调用设备行为", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsDeviceInteractWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/interact/property-control-multi-send", Method: http.MethodPost, Name: "批量调用设备属性", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsDeviceInteractWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/interact/property-control-send", Method: http.MethodPost, Name: "调用设备属性", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsDeviceInteractWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/interact/property-get-report-send", Method: http.MethodPost, Name: "请求设备获取设备最新属性", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsDeviceInteractWrite", IsAuthTenant: 1, Route: "/api/v1/things/device/interact/send-msg", Method: http.MethodPost, Name: "发送消息给设备", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsDeviceMsgRead", IsAuthTenant: 1, Route: "/api/v1/things/device/msg/event-log/index", Method: http.MethodPost, Name: "获取物模型事件历史记录", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsDeviceMsgRead", IsAuthTenant: 1, Route: "/api/v1/things/device/msg/hub-log/index", Method: http.MethodPost, Name: "获取云端诊断日志", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsDeviceMsgRead", IsAuthTenant: 1, Route: "/api/v1/things/device/msg/property-latest/index", Method: http.MethodPost, Name: "获取最新属性", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsDeviceMsgRead", IsAuthTenant: 1, Route: "/api/v1/things/device/msg/property-log/index", Method: http.MethodPost, Name: "获取单个id属性历史记录", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsDeviceMsgRead", IsAuthTenant: 1, Route: "/api/v1/things/device/msg/sdk-log/index", Method: http.MethodPost, Name: "获取设备本地日志", BusinessType: 4, Desc: `获取设备主动上传的sdk日志`},
		{AccessCode: "thingsDeviceMsgRead", IsAuthTenant: 1, Route: "/api/v1/things/device/msg/shadow/index", Method: http.MethodPost, Name: "获取设备影子列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsGroupDeviceRead", IsAuthTenant: 1, Route: "/api/v1/things/group/device/index", Method: http.MethodPost, Name: "获取分组设备列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsGroupDeviceWrite", IsAuthTenant: 1, Route: "/api/v1/things/group/device/multi-create", Method: http.MethodPost, Name: "批量更新分组设备", BusinessType: 1, Desc: `会先删除后新增`},
		{AccessCode: "thingsGroupDeviceWrite", IsAuthTenant: 1, Route: "/api/v1/things/group/device/multi-delete", Method: http.MethodPost, Name: "删除分组设备(支持批量)", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsGroupInfoRead", IsAuthTenant: 1, Route: "/api/v1/things/group/info/index", Method: http.MethodPost, Name: "获取分组列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsGroupInfoRead", IsAuthTenant: 1, Route: "/api/v1/things/group/info/read", Method: http.MethodPost, Name: "获取分组详情信息", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsGroupInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/group/info/create", Method: http.MethodPost, Name: "创建分组", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsGroupInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/group/info/delete", Method: http.MethodPost, Name: "删除分组", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsGroupInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/group/info/update", Method: http.MethodPost, Name: "更新分组信息", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsOpsWorkOrderRead", IsAuthTenant: 1, Route: "/api/v1/things/ops/work-order/index", Method: http.MethodPost, Name: "获取维护工单列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsOpsWorkOrderWrite", IsAuthTenant: 1, Route: "/api/v1/things/ops/work-order/create", Method: http.MethodPost, Name: "创建工单", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsOpsWorkOrderWrite", IsAuthTenant: 1, Route: "/api/v1/things/ops/work-order/update", Method: http.MethodPost, Name: "更新维护工单", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsOtaFirmwareRead", IsAuthTenant: 1, Route: "/api/v1/things/ota/firmware/device-info-read", Method: http.MethodPost, Name: "获取升级包可选设备信息,包含可用版本", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsOtaFirmwareRead", IsAuthTenant: 1, Route: "/api/v1/things/ota/firmware/index", Method: http.MethodPost, Name: "获取升级包列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsOtaFirmwareRead", IsAuthTenant: 1, Route: "/api/v1/things/ota/firmware/read", Method: http.MethodPost, Name: "获取升级包详情", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsOtaFirmwareWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/firmware/create", Method: http.MethodPost, Name: "创建升级包版本", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsOtaFirmwareWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/firmware/delete", Method: http.MethodPost, Name: "删除升级包", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsOtaFirmwareWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/firmware/update", Method: http.MethodPost, Name: "更新升级包", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsOtaJobRead", IsAuthTenant: 1, Route: "/api/v1/things/ota/job/read", Method: http.MethodPost, Name: "查询指定升级批次的详情", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsOtaJobWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/job/cancel", Method: http.MethodPost, Name: "取消动态升级策略", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsOtaJobWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/job/deviceIndex", Method: http.MethodPost, Name: "获取设备所在的升级包升级批次列表", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsOtaJobWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/job/dynamicCreate", Method: http.MethodPost, Name: "创建动态升级批次", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsOtaJobWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/job/firmwareIndex", Method: http.MethodPost, Name: "获取升级包下的升级任务批次列表", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsOtaJobWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/job/staticCreate", Method: http.MethodPost, Name: "创建静态升级批次", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsOtaJobWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/job/verify", Method: http.MethodPost, Name: "验证升级包", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsOtaOtaFirmwareRead", IsAuthTenant: 1, Route: "/api/v1/things/ota/otaFirmware/index", Method: http.MethodPost, Name: "升级包列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsOtaOtaFirmwareRead", IsAuthTenant: 1, Route: "/api/v1/things/ota/otaFirmware/read", Method: http.MethodPost, Name: "查询升级包", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsOtaOtaFirmwareWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/otaFirmware/create", Method: http.MethodPost, Name: "添加升级包", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsOtaOtaFirmwareWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/otaFirmware/delete", Method: http.MethodPost, Name: "删除升级包", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsOtaOtaFirmwareWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/otaFirmware/update", Method: http.MethodPost, Name: "更新升级包", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsOtaTaskRead", IsAuthTenant: 1, Route: "/api/v1/things/ota/task/device-index", Method: http.MethodPost, Name: "批次设备列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsOtaTaskRead", IsAuthTenant: 1, Route: "/api/v1/things/ota/task/index", Method: http.MethodPost, Name: "获取升级批次任务列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsOtaTaskRead", IsAuthTenant: 1, Route: "/api/v1/things/ota/task/read", Method: http.MethodPost, Name: "升级任务信息", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsOtaTaskWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/task/analysis", Method: http.MethodPost, Name: "升级状态统计", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsOtaTaskWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/task/cancel", Method: http.MethodPost, Name: "取消所有升级中的任务", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsOtaTaskWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/task/create", Method: http.MethodPost, Name: "创建升级任务", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsOtaTaskWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/task/device-cancel", Method: http.MethodPost, Name: "取消单个设备升级", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsOtaTaskWrite", IsAuthTenant: 1, Route: "/api/v1/things/ota/task/device-retry", Method: http.MethodPost, Name: "重试单个设备升级", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsProductCategoryRead", IsAuthTenant: 1, Route: "/api/v1/things/product/category/index", Method: http.MethodPost, Name: "获取产品品类列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsProductCategoryRead", IsAuthTenant: 1, Route: "/api/v1/things/product/category/read", Method: http.MethodPost, Name: "获取产品品类详情", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsProductCategoryWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/category/create", Method: http.MethodPost, Name: "新增产品品类", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsProductCategoryWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/category/delete", Method: http.MethodPost, Name: "删除产品品类", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsProductCategoryWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/category/update", Method: http.MethodPost, Name: "更新产品品类", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsProductCustomRead", IsAuthTenant: 1, Route: "/api/v1/things/product/custom/read", Method: http.MethodPost, Name: "获取产品自定义信息", BusinessType: 4, Desc: `物联网平台通过定义一种物的描述语言来描述物模型模块和功能，称为TSL（Thing Specification Language）`},
		{AccessCode: "thingsProductCustomWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/custom/update", Method: http.MethodPost, Name: "更新自定义信息", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsProductInfoRead", IsAuthTenant: 1, Route: "/api/v1/things/product/info/index", Method: http.MethodPost, Name: "获取产品列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsProductInfoRead", IsAuthTenant: 1, Route: "/api/v1/things/product/info/read", Method: http.MethodPost, Name: "获取产品详情", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsProductInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/info/create", Method: http.MethodPost, Name: "新增产品", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsProductInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/info/delete", Method: http.MethodPost, Name: "删除产品", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsProductInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/info/update", Method: http.MethodPost, Name: "更新产品", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsProductRemoteConfigRead", IsAuthTenant: 1, Route: "/api/v1/things/product/remote-config/index", Method: http.MethodPost, Name: "获取配置列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsProductRemoteConfigRead", IsAuthTenant: 1, Route: "/api/v1/things/product/remote-config/lastest-read", Method: http.MethodPost, Name: "获取最新配置", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsProductRemoteConfigWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/remote-config/create", Method: http.MethodPost, Name: "创建配置", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsProductRemoteConfigWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/remote-config/push-all", Method: http.MethodPost, Name: "推送配置", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsProductSchemaRead", IsAuthTenant: 1, Route: "/api/v1/things/product/schema/index", Method: http.MethodPost, Name: "获取产品物模型列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsProductSchemaRead", IsAuthTenant: 1, Route: "/api/v1/things/product/schema/tsl-read", Method: http.MethodPost, Name: "获取产品物模型tsl", BusinessType: 4, Desc: `物联网平台通过定义一种物的描述语言来描述物模型模块和功能，称为TSL（Thing Specification Language）`},
		{AccessCode: "thingsProductSchemaWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/schema/create", Method: http.MethodPost, Name: "新增物模型功能", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsProductSchemaWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/schema/delete", Method: http.MethodPost, Name: "删除物模型功能", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsProductSchemaWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/schema/multi-create", Method: http.MethodPost, Name: "批量新增物模型功能", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsProductSchemaWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/schema/tsl-import", Method: http.MethodPost, Name: "导入物模型tsl", BusinessType: 5, Desc: `物联网平台通过定义一种物的描述语言来描述物模型模块和功能，称为TSL（Thing Specification Language）`},
		{AccessCode: "thingsProductSchemaWrite", IsAuthTenant: 1, Route: "/api/v1/things/product/schema/update", Method: http.MethodPost, Name: "更新物模型功能", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsProtocolInfoRead", IsAuthTenant: 1, Route: "/api/v1/things/protocol/info/index", Method: http.MethodPost, Name: "获取协议列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsProtocolInfoRead", IsAuthTenant: 1, Route: "/api/v1/things/protocol/info/read", Method: http.MethodPost, Name: "获取协议详情", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsProtocolInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/protocol/info/create", Method: http.MethodPost, Name: "新增协议", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsProtocolInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/protocol/info/delete", Method: http.MethodPost, Name: "删除协议", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsProtocolInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/protocol/info/update", Method: http.MethodPost, Name: "更新协议", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsRuleAlarmRead", IsAuthTenant: 1, Route: "/api/v1/things/rule/alarm/deal-record/index", Method: http.MethodPost, Name: "获取告警处理记录列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsRuleAlarmRead", IsAuthTenant: 1, Route: "/api/v1/things/rule/alarm/info/index", Method: http.MethodPost, Name: "获取告警信息列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsRuleAlarmRead", IsAuthTenant: 1, Route: "/api/v1/things/rule/alarm/info/read", Method: http.MethodPost, Name: "获取告警详情", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsRuleAlarmRead", IsAuthTenant: 1, Route: "/api/v1/things/rule/alarm/log/index", Method: http.MethodPost, Name: "获取告警流水日志记录列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsRuleAlarmRead", IsAuthTenant: 1, Route: "/api/v1/things/rule/alarm/record/index", Method: http.MethodPost, Name: "获取告警记录列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsRuleAlarmWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/alarm/deal-record/create", Method: http.MethodPost, Name: "新增告警处理记录", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsRuleAlarmWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/alarm/info/create", Method: http.MethodPost, Name: "新增告警", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsRuleAlarmWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/alarm/info/delete", Method: http.MethodPost, Name: "删除告警", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsRuleAlarmWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/alarm/info/update", Method: http.MethodPost, Name: "更新告警", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsRuleAlarmWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/alarm/scene/delete", Method: http.MethodPost, Name: "删除告警和场景的关联", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsRuleAlarmWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/alarm/scene/multi-update", Method: http.MethodPost, Name: "更新告警和场景的关联", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsRuleDeviceTimerRead", IsAuthTenant: 1, Route: "/api/v1/things/rule/device-timer/info/index", Method: http.MethodPost, Name: "获取场景列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsRuleDeviceTimerRead", IsAuthTenant: 1, Route: "/api/v1/things/rule/device-timer/info/read", Method: http.MethodPost, Name: "获取场景信息", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsRuleDeviceTimerWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/device-timer/info/create", Method: http.MethodPost, Name: "创建场景信息", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsRuleDeviceTimerWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/device-timer/info/delete", Method: http.MethodPost, Name: "删除场景信息", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsRuleDeviceTimerWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/device-timer/info/update", Method: http.MethodPost, Name: "更新场景信息", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsRuleFlowRead", IsAuthTenant: 1, Route: "/api/v1/things/rule/flow/info/index", Method: http.MethodPost, Name: "获取流列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsRuleFlowWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/flow/info/create", Method: http.MethodPost, Name: "创建流", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsRuleFlowWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/flow/info/delete", Method: http.MethodPost, Name: "删除流", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsRuleFlowWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/flow/info/update", Method: http.MethodPost, Name: "修改流", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsRuleSceneRead", IsAuthTenant: 1, Route: "/api/v1/things/rule/scene/info/index", Method: http.MethodPost, Name: "获取场景列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsRuleSceneRead", IsAuthTenant: 1, Route: "/api/v1/things/rule/scene/info/read", Method: http.MethodPost, Name: "获取场景信息", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsRuleSceneWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/scene/info/create", Method: http.MethodPost, Name: "创建场景信息", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsRuleSceneWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/scene/info/delete", Method: http.MethodPost, Name: "删除场景信息", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsRuleSceneWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/scene/info/manually-trigger", Method: http.MethodPost, Name: "手动触发场景联动", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsRuleSceneWrite", IsAuthTenant: 1, Route: "/api/v1/things/rule/scene/info/update", Method: http.MethodPost, Name: "更新场景信息", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsSchemaCommonRead", IsAuthTenant: 1, Route: "/api/v1/things/schema/common/index", Method: http.MethodPost, Name: "获取物模型列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsSchemaCommonWrite", IsAuthTenant: 1, Route: "/api/v1/things/schema/common/create", Method: http.MethodPost, Name: "新增物模型功能", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsSchemaCommonWrite", IsAuthTenant: 1, Route: "/api/v1/things/schema/common/delete", Method: http.MethodPost, Name: "删除物模型功能", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsUserDeviceRead", IsAuthTenant: 1, Route: "/api/v1/things/user/device/collect/index", Method: http.MethodPost, Name: "获取用户设备收藏列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsUserDeviceRead", IsAuthTenant: 1, Route: "/api/v1/things/user/device/share/index", Method: http.MethodPost, Name: "获取用户设备分享列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsUserDeviceRead", IsAuthTenant: 1, Route: "/api/v1/things/user/device/share/read", Method: http.MethodPost, Name: "获取用户设备分享详情", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsUserDeviceWrite", IsAuthTenant: 1, Route: "/api/v1/things/user/device/collect/multi-create", Method: http.MethodPost, Name: "批量新增用户收藏的设备", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsUserDeviceWrite", IsAuthTenant: 1, Route: "/api/v1/things/user/device/collect/multi-delete", Method: http.MethodPost, Name: "批量删除用户收藏的设备", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsUserDeviceWrite", IsAuthTenant: 1, Route: "/api/v1/things/user/device/share/create", Method: http.MethodPost, Name: "分享用户设备", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsUserDeviceWrite", IsAuthTenant: 1, Route: "/api/v1/things/user/device/share/delete", Method: http.MethodPost, Name: "删除用户设备分享", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsUserDeviceWrite", IsAuthTenant: 1, Route: "/api/v1/things/user/device/share/update", Method: http.MethodPost, Name: "更新用户设备分享权限", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsVidmgrCtrlWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/ctrl/getsvr", Method: http.MethodPost, Name: "获取流服务状态", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrCtrlWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/ctrl/restart", Method: http.MethodPost, Name: "重启流服务", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrCtrlWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/ctrl/setsvr", Method: http.MethodPost, Name: "修改流服务状态", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrGbsipWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/gbsip/createchn", Method: http.MethodPost, Name: "创建通道", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrGbsipWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/gbsip/createdev", Method: http.MethodPost, Name: "创建设备", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrGbsipWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/gbsip/deletechn", Method: http.MethodPost, Name: "删除通道", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrGbsipWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/gbsip/deletedev", Method: http.MethodPost, Name: "删除设备", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrGbsipWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/gbsip/indexchn", Method: http.MethodPost, Name: "获取通道列表", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrGbsipWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/gbsip/indexdev", Method: http.MethodPost, Name: "获取设备列表", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrGbsipWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/gbsip/playchn", Method: http.MethodPost, Name: "通道播放", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrGbsipWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/gbsip/readchn", Method: http.MethodPost, Name: "获取通道详细", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrGbsipWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/gbsip/readdev", Method: http.MethodPost, Name: "获取设备详细", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrGbsipWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/gbsip/readinfo", Method: http.MethodPost, Name: "获取服务详细", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrGbsipWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/gbsip/stopchn", Method: http.MethodPost, Name: "通道暂停播放", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrGbsipWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/gbsip/updatechn", Method: http.MethodPost, Name: "更新通道信息", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrGbsipWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/gbsip/updatedev", Method: http.MethodPost, Name: "更新设备信息", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrInfoRead", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/info/index", Method: http.MethodPost, Name: "获取流服务器列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsVidmgrInfoRead", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/info/read", Method: http.MethodPost, Name: "获取流服详细", BusinessType: 4, Desc: `{
  "vidmgrID":"1113459"
}`},
		{AccessCode: "thingsVidmgrInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/info/active", Method: http.MethodPost, Name: "激活流服务器", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/info/count", Method: http.MethodPost, Name: "获取设备在线数", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/info/create", Method: http.MethodPost, Name: "新增流服务器", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsVidmgrInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/info/delete", Method: http.MethodPost, Name: "删除流服务器", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsVidmgrInfoWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/info/update", Method: http.MethodPost, Name: "更新流服务器", BusinessType: 2, Desc: ``},
		{AccessCode: "thingsVidmgrStreamRead", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/stream/index", Method: http.MethodPost, Name: "获取流列表", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsVidmgrStreamRead", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/stream/read", Method: http.MethodPost, Name: "查询流详细", BusinessType: 4, Desc: ``},
		{AccessCode: "thingsVidmgrStreamWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/stream/count", Method: http.MethodPost, Name: "统计在线的流", BusinessType: 5, Desc: ``},
		{AccessCode: "thingsVidmgrStreamWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/stream/create", Method: http.MethodPost, Name: "创建流（拉流）", BusinessType: 1, Desc: ``},
		{AccessCode: "thingsVidmgrStreamWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/stream/delete", Method: http.MethodPost, Name: "删除流", BusinessType: 3, Desc: ``},
		{AccessCode: "thingsVidmgrStreamWrite", IsAuthTenant: 1, Route: "/api/v1/things/vidmgr/stream/update", Method: http.MethodPost, Name: "更新流信息", BusinessType: 2, Desc: ``},
	}
)
