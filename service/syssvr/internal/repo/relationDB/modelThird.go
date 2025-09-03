package relationDB

import (
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
)

type SysThirdInfo struct {
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                     // id编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;uniqueIndex:idx_sys_tenant_config_ri_mi;type:VARCHAR(50);NOT NULL"` // 租户编码
	Name       string              `gorm:"column:name;uniqueIndex:idx_sys_tenant_info_name;type:VARCHAR(100);NOT NULL"`          // 租户名称
	AppType    def.ThirdType       `gorm:"column:app_type;uniqueIndex:idx_sys_user_info_tc_wui;type:varchar(64);NOT NULL"`
	AppID      string              `gorm:"column:app_id;type:VARCHAR(50);default:'';"`
	AppKey     string              `gorm:"column:app_key;type:VARCHAR(50);default:'';"`
	AppSecret  string              `gorm:"column:app_secret;type:VARCHAR(200);default:'';"`
	Desc       string              `gorm:"column:desc;type:VARCHAR(500);"` // 备注
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_tenant_app_tc_ac"`
}

func (m *SysThirdInfo) TableName() string {
	return "sys_third_info"
}
