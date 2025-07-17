package relationDB

import (
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/stores"
)

// 角色管理表
type SysRoleInfo struct {
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                 // id编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;uniqueIndex:idx_sys_role_info_tc_ac;type:VARCHAR(50);NOT NULL"` // 角色所属租户编码
	Name       string              `gorm:"column:name;uniqueIndex:idx_sys_role_info_tc_ac;type:VARCHAR(100);NOT NULL"`       // 角色名称
	Desc       string              `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`                                           //描述
	Code       string              `gorm:"column:code;uniqueIndex:idx_sys_role_info_tc_ac;type:VARCHAR(50);NOT NULL"`        // 编码
	Status     int64               `gorm:"column:status;type:SMALLINT;default:1"`                                            // 状态  1:启用,2:禁用
	Apps       []*SysRoleApp       `gorm:"foreignKey:RoleID;references:ID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_role_info_tc_ac"`
}

func (m *SysRoleInfo) TableName() string {
	return "sys_role_info"
}

// 应用菜单关联表
type SysRoleApp struct {
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                // id编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;uniqueIndex:idx_sys_role_app_tc_ac;type:VARCHAR(50);NOT NULL"` // 角色所属租户编码
	RoleID     int64               `gorm:"column:role_id;uniqueIndex:idx_sys_role_app_tc_ac;NOT NULL;type:BIGINT"`          // 角色ID
	AppCode    string              `gorm:"column:app_code;uniqueIndex:idx_sys_role_app_tc_ac;type:VARCHAR(50);NOT NULL"`    // 应用编码
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_role_app_tc_ac"`
}

func (m *SysRoleApp) TableName() string {
	return "sys_role_app"
}

// 应用菜单关联表
type SysRoleModule struct {
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                   // id编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;uniqueIndex:idx_sys_role_module_tc_ac;type:VARCHAR(50);NOT NULL"` // 角色所属租户编码
	RoleID     int64               `gorm:"column:role_id;uniqueIndex:idx_sys_role_module_tc_ac;NOT NULL;type:BIGINT"`          // 角色ID
	AppCode    string              `gorm:"column:app_code;uniqueIndex:idx_sys_role_module_tc_ac;type:VARCHAR(50);NOT NULL"`    // 应用编码
	ModuleCode string              `gorm:"column:module_code;uniqueIndex:idx_sys_role_module_tc_ac;type:VARCHAR(50);NOT NULL"` // 模块编码
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_role_module_tc_ac"`
}

func (m *SysRoleModule) TableName() string {
	return "sys_role_module"
}

// 应用菜单关联表
type SysRoleMenu struct {
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                 // id编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;uniqueIndex:idx_sys_role_menu_ri_mi;type:VARCHAR(50);NOT NULL"` // 角色所属租户编码
	RoleID     int64               `gorm:"column:role_id;uniqueIndex:idx_sys_role_menu_ri_mi;NOT NULL;type:BIGINT"`          // 角色ID
	AppCode    string              `gorm:"column:app_code;uniqueIndex:idx_sys_role_menu_ri_mi;type:VARCHAR(50);NOT NULL"`    // 应用编码
	ModuleCode string              `gorm:"column:module_code;uniqueIndex:idx_sys_role_menu_ri_mi;type:VARCHAR(50);NOT NULL"` // 模块编码
	MenuID     int64               `gorm:"column:menu_id;uniqueIndex:idx_sys_role_menu_ri_mi;NOT NULL;type:BIGINT"`          // 菜单ID
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_role_menu_ri_mi"`
}

func (m *SysRoleMenu) TableName() string {
	return "sys_role_menu"
}

// api权限管理
type SysRoleAccess struct {
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                   // id编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;uniqueIndex:idx_sys_role_access_ri_mi;type:VARCHAR(50);NOT NULL"` // 角色所属租户编码
	RoleID     int64               `gorm:"column:role_id;uniqueIndex:idx_sys_role_access_ri_mi;NOT NULL;type:BIGINT"`          // 角色ID
	AccessCode string              `gorm:"column:access_code;uniqueIndex:idx_sys_role_access_ri_mi;type:VARCHAR(50);NOT NULL"` // 范围编码
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_role_access_ri_mi"`
	Api         *SysAccessInfo     `gorm:"foreignKey:AccessCode;references:Code"`
}

func (m *SysRoleAccess) TableName() string {
	return "sys_role_access"
}
