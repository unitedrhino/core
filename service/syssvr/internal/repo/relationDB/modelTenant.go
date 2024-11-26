package relationDB

import (
	"gitee.com/unitedrhino/core/service/syssvr/domain/tenant"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/users"
)

// 租户信息表
type SysTenantInfo struct {
	ID               int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`        // id编号
	Code             stores.TenantCode `gorm:"column:code;uniqueIndex:code;type:VARCHAR(100);NOT NULL"` // 租户编码
	Name             string            `gorm:"column:name;uniqueIndex:name;type:VARCHAR(100);NOT NULL"` // 租户名称
	AdminUserID      int64             `gorm:"column:admin_user_id;type:BIGINT;NOT NULL"`               // 超级管理员id
	AdminRoleID      int64             `gorm:"column:admin_role_id;type:BIGINT;NOT NULL"`               // 超级角色
	Desc             string            `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`                  //应用描述
	DefaultProjectID int64             `gorm:"column:default_project_id;type:BIGINT;NOT NULL"`
	DefaultAreaID    int64             `gorm:"column:default_area_id;type:BIGINT;NOT NULL"`
	UserCount        int64             `gorm:"column:user_count;type:bigint;default:0;"` //租户下用户统计
	SysTenantOem
	Status int64 `gorm:"column:status;type:BIGINT;NOT NULL;default:1"` //租戶状态: 1启用 2禁用
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:code;uniqueIndex:name"`
}

type SysTenantOem struct {
	BackgroundColour string `gorm:"column:background_colour;type:VARCHAR(54);"` //背景颜色
	BackgroundDesc   string `gorm:"column:background_desc;type:VARCHAR(54);"`   //背景描述
	BackgroundImg    string `gorm:"column:background_img;type:VARCHAR(512);"`   //背景图片
	LogoImg          string `gorm:"column:logo_img;type:VARCHAR(512);"`         //租户logo地址
	Title            string `gorm:"column:title;type:VARCHAR(100);"`            //中文标题
	TitleEn          string `gorm:"column:title_en;type:VARCHAR(100);"`         //英文标题
	Footer           string `gorm:"column:footer;type:text;"`                   //页尾
}

func (m *SysTenantInfo) TableName() string {
	return "sys_tenant_info"
}

// 租户开放认证
type SysTenantOpenWebhook struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                        // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"`          // 租户编码
	Code       string            `gorm:"column:code;type:VARCHAR(50);uniqueIndex:tc_ac;NOT NULL"`                 //业务里定义的,推送的内容
	Uri        string            `gorm:"column:uri;type:VARCHAR(100);NOT NULL"`                                   // 参考: /api/v1/system/user/self/captcha?fwefwf=gwgweg&wefaef=gwegwe
	Hosts      []string          `gorm:"column:hosts;type:json;serializer:json;NOT NULL;default:'[]';NOT NULL"`   //访问的地址 host or host:port
	Desc       string            `gorm:"column:desc;type:VARCHAR(500);"`                                          // 备注
	Handler    map[string]string `gorm:"column:handler;type:json;serializer:json;NOT NULL;default:'{}';NOT NULL"` //http头
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:tc_ac"`
}

func (m *SysTenantOpenWebhook) TableName() string {
	return "sys_tenant_open_webhook"
}

//// 租户自定义表
//type TenantOem struct {

//}

// 租户下的应用列表
type SysTenantApp struct {
	ID             int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode     stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 租户编码
	AppCode        string            `gorm:"column:app_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"`    // 应用编码 这里只关联主应用,主应用授权,子应用也授权了
	DingMini       *SysTenantThird   `gorm:"embedded;embeddedPrefix:ding_mini_"`                             //钉钉企业应用接入
	Android        *SysThirdApp      `gorm:"embedded;embeddedPrefix:android_"`                               //安卓应用
	WxMini         *SysTenantThird   `gorm:"embedded;embeddedPrefix:wx_mini_"`                               //微信小程序接入
	WxOpen         *SysTenantThird   `gorm:"embedded;embeddedPrefix:wx_open_"`                               //微信公众号接入
	LoginTypes     []users.RegType   `gorm:"column:login_types;type:json;serializer:json"`                   //支持的登录类型(不填支持全部登录方式):  	 "email":邮箱 "phone":手机号  "wxMiniP":微信小程序  "wxOfficial": 微信公众号登录   "dingApp":钉钉应用(包含小程序,h5等方式)  "pwd":账号密码注册
	IsAutoRegister int64             `gorm:"column:is_auto_register;type:BIGINT;default:1"`                  //登录未注册是否自动注册
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:tc_ac"`
}

func (m *SysTenantApp) TableName() string {
	return "sys_tenant_app"
}

// 租户下的应用列表
type SysTenantAppModule struct {
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 租户编码
	SysAppModule
}

func (m *SysTenantAppModule) TableName() string {
	return "sys_tenant_app_module"
}

// 菜单管理表
type SysTenantAppMenu struct {
	TempLateID int64             `gorm:"column:template_id;uniqueIndex:template_id;type:BIGINT;NOT NULL"`      // 模板id
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:template_id;type:VARCHAR(50);NOT NULL"` // 租户编码
	AppCode    string            `gorm:"column:app_code;uniqueIndex:template_id;type:VARCHAR(50);NOT NULL"`    // 应用编码 这里只关联主应用,主应用授权,子应用也授权了
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                     // 编号
	ModuleCode string            `gorm:"column:module_code;type:VARCHAR(50);NOT NULL"`                         // 模块编码
	ParentID   int64             `gorm:"column:parent_id;type:BIGINT;default:1;NOT NULL"`                      // 父菜单ID，一级菜单为1
	Type       int64             `gorm:"column:type;type:BIGINT;default:1;NOT NULL"`                           // 类型   1：菜单或者页面   2：iframe嵌入   3：外链跳转
	Order      int64             `gorm:"column:order;type:BIGINT;default:1;NOT NULL"`                          // 左侧table排序序号
	Name       string            `gorm:"column:name;type:VARCHAR(50);NOT NULL"`                                // 菜单名称
	Path       string            `gorm:"column:path;type:VARCHAR(64);NOT NULL"`                                // 系统的path
	Component  string            `gorm:"column:component;type:VARCHAR(1024);NOT NULL"`                         // 页面
	Icon       string            `gorm:"column:icon;type:VARCHAR(64);NOT NULL"`                                // 图标
	Redirect   string            `gorm:"column:redirect;type:VARCHAR(64);NOT NULL"`                            // 路由重定向
	Body       string            `gorm:"column:body;type:VARCHAR(1024)"`                                       // 菜单自定义数据
	HideInMenu int64             `gorm:"column:hide_in_menu;type:BIGINT;default:2;NOT NULL"`                   // 是否隐藏菜单 1-是 2-否
	IsCommon   int64             `gorm:"column:is_common;type:BIGINT;default:2;"`                              // 是否常用菜单 1-是 2-否
	Children   []*SysModuleMenu  `gorm:"foreignKey:ID;references:ParentID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:template_id;default:0;index"`
}

func (m *SysTenantAppMenu) TableName() string {
	return "sys_tenant_app_menu"
}

// 租户下的邮箱配置
type SysTenantConfig struct {
	ID                        int64                               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode                stores.TenantCode                   `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 租户编码
	RegisterRoleID            int64                               `gorm:"column:register_role_id;type:BIGINT;NOT NULL"`                   //注册分配的角色id
	DeviceLimit               int64                               `gorm:"column:device_limit;type:BIGINT;default:0"`                      // 租户下的设备数量限制,0为不限制
	CheckUserDelete           int64                               `gorm:"column:check_user_delete;type:BIGINT;default:2"`                 // 1(禁止项目管理员注销账号) 2(不禁止项目管理员注销账号)
	WeatherKey                string                              `gorm:"column:weather_key;type:VARCHAR(50);default:'';"`                //参考: https://dev.qweather.com/
	RegisterAutoCreateProject []*tenant.RegisterAutoCreateProject `gorm:"column:register_auto_create_project;type:json;serializer:json;default:'[]'"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:ri_mi"`
}

type SysTenantEmail struct {
	From     string `gorm:"column:from;type:VARCHAR(50);default:'';NOT NULL"`     // 发件人  你自己要发邮件的邮箱
	Host     string `gorm:"column:host;type:VARCHAR(50);default:'';NOT NULL"`     // 服务器地址 例如 smtp.qq.com  请前往QQ或者你要发邮件的邮箱查看其smtp协议
	Secret   string `gorm:"column:secret;type:VARCHAR(50);default:'';NOT NULL"`   // 密钥    用于登录的密钥 最好不要用邮箱密码 去邮箱smtp申请一个用于登录的密钥
	Nickname string `gorm:"column:nickname;type:VARCHAR(50);default:'';NOT NULL"` // 昵称    发件人昵称 通常为自己的邮箱
	Port     int64  `gorm:"column:port;type:int;default:465"`                     // 端口     请前往QQ或者你要发邮件的邮箱查看其smtp协议 大多为 465
	IsSSL    int64  `gorm:"column:is_ssl;type:int;default:2"`                     // 是否SSL   是否开启SSL
}

//type SysConfigSms struct {
//	From     string `gorm:"column:from;type:VARCHAR(50);default:'';NOT NULL"`     // 发件人  你自己要发邮件的邮箱
//	Host     string `gorm:"column:host;type:VARCHAR(50);default:'';NOT NULL"`     // 服务器地址 例如 smtp.qq.com  请前往QQ或者你要发邮件的邮箱查看其smtp协议
//	Secret   string `gorm:"column:secret;type:VARCHAR(50);default:'';NOT NULL"`   // 密钥    用于登录的密钥 最好不要用邮箱密码 去邮箱smtp申请一个用于登录的密钥
//	Nickname string `gorm:"column:nickname;type:VARCHAR(50);default:'';NOT NULL"` // 昵称    发件人昵称 通常为自己的邮箱
//	Port     int64  `gorm:"column:port;type:int;default:465"`                     // 端口     请前往QQ或者你要发邮件的邮箱查看其smtp协议 大多为 465
//	IsSSL    int64  `gorm:"column:is_ssl;type:int;default:2"`                     // 是否SSL   是否开启SSL
//}

// 第三方app配置
type SysTenantThird struct {
	AppID string `gorm:"column:app_id;type:VARCHAR(50);default:'';"`
	//MiniAppID string `gorm:"column:mini_app_id;type:VARCHAR(50);default:'';"`
	AppKey    string `gorm:"column:app_key;type:VARCHAR(50);default:'';"`
	AppSecret string `gorm:"column:app_secret;type:VARCHAR(200);default:'';"`
}

// 第三方app配置
type SysThirdApp struct {
	Version     string `gorm:"column:version;type:varchar(64);"`       // 应用版本
	FilePath    string `gorm:"column:file_path;type:varchar(256);"`    // 文件路径,拿来下载文件
	VersionDesc string `gorm:"column:version_desc;type:VARCHAR(100);"` //版本说明
}

func (m *SysTenantConfig) TableName() string {
	return "sys_tenant_config"
}

type SysTenantAgreement struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 租户编码
	Code       string            `gorm:"column:code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"`        // 协议编码
	Name       string            `gorm:"column:name;type:VARCHAR(50);NOT NULL"`                          //协议名称
	Title      string            `gorm:"column:title;type:VARCHAR(50);"`
	Content    string            `gorm:"column:content;type:text;"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:tc_ac"`
}

func (m *SysTenantAgreement) TableName() string {
	return "sys_tenant_agreement"
}
