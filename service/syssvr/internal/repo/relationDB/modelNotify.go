package relationDB

import "gitee.com/i-Things/share/stores"

/*

 */

// SysNotifyConfig 通知类型配置
type SysNotifyConfig struct {
	ID           int64    `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`        // id编号
	Code         string   `gorm:"column:code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 通知类型编码
	Name         string   //通知的命名
	SupportTypes []string //支持的通知类型
	Desc         string   `gorm:"column:desc;type:varchar(100);NOT NULL"` // 项目备注
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:pn"`
}

func (m *SysNotifyConfig) TableName() string {
	return "sys_project_info"
}

// 通知模版
type SysNotifyTemplate struct {
	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`        // id编号
	Code       string            `gorm:"column:code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 通知类型编码
	TenantCode string            //限定租户
	Name       string            //通知的命名
	ConfigCode string            //使用的模版code
	Type       string            //模版类型
	Body       string            //模版内容
	Params     map[string]string //变量属性
	Desc       string            `gorm:"column:desc;type:varchar(100);NOT NULL"` // 项目备注
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:pn"`
}
