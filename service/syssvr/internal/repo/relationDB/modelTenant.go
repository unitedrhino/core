package relationDB

import (
	"gitee.com/i-Things/share/stores"
)

// 租户信息表
type SysTenantInfo struct {
	ID               int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`        // id编号
	Code             string `gorm:"column:code;uniqueIndex:code;type:VARCHAR(100);NOT NULL"` // 租户编码
	Name             string `gorm:"column:name;uniqueIndex:name;type:VARCHAR(100);NOT NULL"` // 租户名称
	AdminUserID      int64  `gorm:"column:admin_user_id;type:BIGINT;NOT NULL"`               // 超级管理员id
	AdminRoleID      int64  `gorm:"column:admin_role_id;type:BIGINT;NOT NULL"`               // 超级角色
	Desc             string `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`                  //应用描述
	DefaultProjectID int64  `gorm:"column:default_project_id;type:BIGINT;NOT NULL"`
	SysTenantOem
	Status int64 `gorm:"column:status;type:BIGINT;NOT NULL;default:1"` //租戶状态: 1启用 2禁用
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:code;uniqueIndex:name"`
}

type SysTenantOem struct {
	BackGroundColour string `gorm:"column:background_colour;type:VARCHAR(54);"` //背景颜色
	BackgroundImg    string `gorm:"column:background_img;type:VARCHAR(512);"`   //背景图片
	LogoImg          string `gorm:"column:logo_img;type:VARCHAR(512);"`         //租户logo地址
	Title            string `gorm:"column:title;type:VARCHAR(100);"`            //中文标题
	TitleEn          string `gorm:"column:title_en;type:VARCHAR(100);"`         //英文标题
}

func (m *SysTenantInfo) TableName() string {
	return "sys_tenant_info"
}

// 租户开放认证
type SysTenantOpenAccess struct {
	ID           int64    `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode   string   `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 租户编码
	UserID       int64    `gorm:"column:user_id;uniqueIndex:tc_ac;type:bigint;NOT NULL"`
	Code         string   `gorm:"column:code;type:VARCHAR(50);uniqueIndex:tc_ac;NOT NULL"` //用来标识用来干嘛的
	AccessSecret string   `gorm:"column:access_secret;type:VARCHAR(256);NOT NULL"`
	Desc         string   `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`     //应用描述
	IpRange      []string `gorm:"column:ip_range;type:json;serializer:json;"` //ip白名单
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:tc_ac"`
}

func (m *SysTenantOpenAccess) TableName() string {
	return "sys_tenant_open_access"
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
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 租户编码
	AppCode    string            `gorm:"column:app_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"`    // 应用编码 这里只关联主应用,主应用授权,子应用也授权了
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
	TempLateID int64             `gorm:"column:template_id;type:BIGINT;NOT NULL"`      // 模板id
	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"` // 租户编码
	AppCode    string            `gorm:"column:app_code;type:VARCHAR(50);NOT NULL"`    // 应用编码 这里只关联主应用,主应用授权,子应用也授权了
	SysModuleMenu
}

func (m *SysTenantAppMenu) TableName() string {
	return "sys_tenant_app_menu"
}

// 租户下的邮箱配置
type SysTenantConfig struct {
	ID                    int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode            stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 租户编码
	Email                 *SysTenantEmail   `gorm:"embedded;embeddedPrefix:email_"`                                 //邮箱配置
	DingTalk              *SysTenantThird   `gorm:"embedded;embeddedPrefix:ding_talk_"`                             //钉钉企业应用接入
	WxMini                *SysTenantThird   `gorm:"embedded;embeddedPrefix:wxmini_"`                                //微信小程序接入
	RegisterRoleID        int64             `gorm:"column:register_role_id;type:BIGINT;NOT NULL"`                   //注册分配的角色id
	RegisterCreateProject int64             `gorm:"column:register_create_project;type:int;default:2"`              //注册自动创建项目
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:ri_mi"`
}

type SysTenantLogin struct {
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
	AppID     string `gorm:"column:app_id;type:VARCHAR(50);default:'';NOT NULL"`
	AppKey    string `gorm:"column:app_key;type:VARCHAR(50);default:'';NOT NULL"`
	AppSecret string `gorm:"column:app_secret;type:VARCHAR(200);default:'';NOT NULL"`
}

func (m *SysTenantConfig) TableName() string {
	return "sys_tenant_config"
}
