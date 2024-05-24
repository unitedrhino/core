package relationDB

import (
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
)

/*

 */

// SysNotifyInfo 通知类型配置
type SysNotifyInfo struct {
	ID                  int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                     // id编号
	Group               string            `gorm:"column:group;type:VARCHAR(50);NOT NULL"`                               //分组
	Code                string            `gorm:"column:code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"`              // 通知类型编码
	Name                string            `gorm:"column:name;type:VARCHAR(50);NOT NULL"`                                //通知的命名
	SupportTypes        []def.NotifyType  `gorm:"column:support_types;type:json;serializer:json;NOT NULL;default:'[]'"` //支持的通知类型
	Desc                string            `gorm:"column:desc;type:varchar(100);NOT NULL"`                               // 项目备注
	IsRecord            int64             `gorm:"column:is_record;type:BIGINT"`                                         //是否记录该消息,是的情况下会将消息存一份到消息中心
	DefaultSubject      string            `gorm:"column:default_subject;type:VARCHAR(256);NOT NULL"`                    //默认消息主题
	DefaultBody         string            `gorm:"column:default_body;type:VARCHAR(512);default:''"`                     //默认模版内容
	DefaultTemplateCode string            `gorm:"column:default_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"`      // 通知类型编码
	DefaultSignName     string            `gorm:"column:default_sign_name;type:VARCHAR(50);default:''"`                 //默认签名(短信)
	Params              map[string]string `gorm:"column:params;type:json;serializer:json;NOT NULL;default:'{}'"`        //变量属性 key是参数,value是描述
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:ri_mi;"`
}

func (m *SysNotifyInfo) TableName() string {
	return "sys_notify_info"
}

// 通知配置
type SysNotifyTemplate struct {
	ID           int64                   `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode   stores.TenantCode       `gorm:"column:tenant_code;type:VARCHAR(50);default:'common'"`           //限定租户,不填是通用的
	Name         string                  `gorm:"column:name;type:VARCHAR(50);NOT NULL"`                          //通知的命名
	NotifyCode   string                  `gorm:"column:notify_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` //对应的配置Code
	Type         def.NotifyType          `gorm:"column:type;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"`        //对应的配置类型 sms email
	TemplateCode string                  `gorm:"column:code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"`        // 短信通知模版编码
	SignName     string                  `gorm:"column:sign_name;type:VARCHAR(50);default:''"`                   //签名(短信)
	Subject      string                  `gorm:"column:subject;type:VARCHAR(256);NOT NULL"`                      //默认消息主题
	Body         string                  `gorm:"column:body;type:VARCHAR(512);default:''"`                       //默认模版内容
	Desc         string                  `gorm:"column:desc;type:varchar(100)"`                                  // 备注
	ChannelID    int64                   `gorm:"column:channel_id;type:BIGINT;"`
	Channel      *SysTenantNotifyChannel `gorm:"foreignKey:ID;references:ChannelID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:ri_mi"`
}

func (m *SysNotifyTemplate) TableName() string {
	return "sys_notify_template"
}

// 租户下的通知配置
type SysTenantNotifyTemplate struct {
	ID         int64              `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode stores.TenantCode  `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 租户编码
	NotifyCode string             `gorm:"column:notify_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` //对应的配置Code
	Type       string             `gorm:"column:type;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"`        //对应的类型
	TemplateID int64              `gorm:"column:template_id;type:BIGINT;default:1"`                       //绑定的模板id,1为默认
	Template   *SysNotifyTemplate `gorm:"foreignKey:ID;references:TemplateID"`
	Config     *SysNotifyInfo     `gorm:"foreignKey:Code;references:NotifyCode"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:ri_mi"`
}

func (m *SysTenantNotifyTemplate) TableName() string {
	return "sys_tenant_notify_template"
}

// 租户下的通道配置
type SysTenantNotifyChannel struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`        // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`            // 租户编码,为common是公共的
	Type       def.NotifyType    `gorm:"column:type;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` //对应的配置类型 sms email
	Email      *SysTenantEmail   `gorm:"embedded;embeddedPrefix:email_"`                          //邮箱配置
	//AppCode    string            `gorm:"column:app_code;type:VARCHAR(100)"`                       //绑定的应用
	WebHook string `gorm:"column:webhook;type:VARCHAR(256)"` //钉钉webhook模式及企业微信webhook方式
	Name    string `gorm:"column:name;uniqueIndex:ri_mi;type:VARCHAR(100);NOT NULL"`
	Desc    string `gorm:"column:desc;type:VARCHAR(100);NOT NULL"` //应用描述
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:ri_mi"`
}

func (m *SysTenantNotifyChannel) TableName() string {
	return "sys_tenant_channel"
}

type SysMessageInfo struct {
	ID             int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`         // id编号
	TenantCode     stores.TenantCode `gorm:"column:tenant_code;index:ri_mi;type:VARCHAR(50);NOT NULL"` // 租户编码
	Group          string            `gorm:"column:group;type:VARCHAR(50);NOT NULL"`                   //消息分类
	NotifyCode     string            `gorm:"column:notify_code;type:VARCHAR(50);NOT NULL"`             //对应的配置Code
	Subject        string            `gorm:"column:subject;type:VARCHAR(256);NOT NULL"`                //消息主题
	Body           string            `gorm:"column:body;type:text;NOT NULL"`                           //消息内容
	Str1           string            `gorm:"column:str1;index:ri_mi;type:VARCHAR(50);NOT NULL"`        //自定义字段(用来添加搜索索引),如产品id
	Str2           string            `gorm:"column:str2;index:ri_mi;type:VARCHAR(50);NOT NULL"`        //自定义字段(用来添加搜索索引),如设备id
	Str3           string            `gorm:"column:str3;index:ri_mi;type:VARCHAR(50);NOT NULL"`
	IsGlobal       int64             `gorm:"column:is_global;index;type:bigint;default:2"`        //是否是全局消息,是的话所有用户都能看到
	IsDirectNotify int64             `gorm:"column:is_direct_notify;index;type:bigint;default:2"` //是否是发送通知消息创建
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;"`
}

func (m *SysMessageInfo) TableName() string {
	return "sys_message_info"
}

type SysUserMessage struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 租户编码
	UserID     int64             `gorm:"column:user_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"`          // 用户ID
	Group      string            `gorm:"column:group;type:VARCHAR(50);NOT NULL"`                         //消息分类
	MessageID  int64             `gorm:"column:message_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"`       //消息id
	IsRead     int64             `gorm:"column:is_read;NOT NULL;type:BIGINT;default:2"`                  //是否已读
	Message    *SysMessageInfo   `gorm:"foreignKey:ID;references:MessageID"`
	stores.Time
}

func (m *SysUserMessage) TableName() string {
	return "sys_user_message"
}
