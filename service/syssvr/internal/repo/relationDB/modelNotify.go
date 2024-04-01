package relationDB

import "gitee.com/i-Things/share/stores"

/*

 */

// SysNotifyConfig 通知类型配置
type SysNotifyConfig struct {
	ID           int64    `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                     // id编号
	Code         string   `gorm:"column:code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"`              // 通知类型编码
	Name         string   `gorm:"column:name;type:VARCHAR(50);NOT NULL"`                                //通知的命名
	SupportTypes []string `gorm:"column:support_types;type:json;serializer:json;NOT NULL;default:'[]'"` //支持的通知类型
	Desc         string   `gorm:"column:desc;type:varchar(100);NOT NULL"`                               // 项目备注
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:pn"`
}

func (m *SysNotifyConfig) TableName() string {
	return "sys_notify_config"
}

// 通知模版
type SysNotifyTemplate struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode string            `gorm:"column:tenant_code;type:VARCHAR(50);default:'all'"`              //限定租户,不填是通用的
	Name       string            `gorm:"column:name;type:VARCHAR(50);NOT NULL"`                          //通知的命名
	ConfigCode string            `gorm:"column:config_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` //对应的配置Code
	Type       string            `gorm:"column:type;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"`        //对应的配置类型 sms email
	Code       string            `gorm:"column:code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"`        // 通知类型编码
	SignName   string            `gorm:"column:sign_name;type:VARCHAR(50);default:''"`                   //签名(短信)
	Body       string            `gorm:"column:body;type:VARCHAR(512);default:''"`                       //模版内容(email,钉钉)
	Params     map[string]string `gorm:"column:params;type:json;serializer:json;NOT NULL;default:'{}'"`  //变量属性 key是参数,value是描述
	Desc       string            `gorm:"column:desc;type:varchar(100);NOT NULL"`                         // 项目备注
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
	ConfigCode string             `gorm:"column:config_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` //对应的配置Code
	Type       string             `gorm:"column:type;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"`        //对应的类型
	TemplateID int64              `gorm:"column:template_id;type:BIGINT"`                                 //绑定的模板code
	Template   *SysNotifyTemplate `gorm:"foreignKey:ID;references:TemplateID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:ri_mi"`
}

func (m *SysTenantNotifyTemplate) TableName() string {
	return "sys_tenant_notify_template"
}
