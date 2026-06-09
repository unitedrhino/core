package relationDB

import (
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
)

// SysUserPushClient uni-push 2.0 客户端 cid 绑定
type SysUserPushClient struct {
	ID           int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`
	TenantCode   dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;uniqueIndex:idx_sys_user_push_client_tc_un"`
	UserID       int64               `gorm:"column:user_id;type:BIGINT;NOT NULL;uniqueIndex:idx_sys_user_push_client_tc_un"`
	PushClientID string              `gorm:"column:push_client_id;type:VARCHAR(64);NOT NULL;uniqueIndex:idx_sys_user_push_client_tc_un"`
	Platform     string              `gorm:"column:platform;type:VARCHAR(16);NOT NULL"`
	AppID        string              `gorm:"column:app_id;type:VARCHAR(32);NOT NULL"`
	AppVersion   string              `gorm:"column:app_version;type:VARCHAR(32)"`
	IsActive     int64               `gorm:"column:is_active;type:BIGINT;NOT NULL;default:1"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_user_push_client_tc_un"`
}

func (m *SysUserPushClient) TableName() string {
	return "sys_user_push_client"
}

func (m *SysUserPushClient) BeforeCreate() {
	if m.IsActive == 0 {
		m.IsActive = def.True
	}
}
