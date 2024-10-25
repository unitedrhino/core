package relationDB

import (
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
)

// 租户开放认证
type SysDataOpenAccess struct {
	ID           int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode   stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 租户编码
	UserID       int64             `gorm:"column:user_id;uniqueIndex:tc_ac;type:bigint;NOT NULL"`
	Code         string            `gorm:"column:code;type:VARCHAR(50);uniqueIndex:tc_ac;NOT NULL"` //用来标识用来干嘛的
	AccessSecret string            `gorm:"column:access_secret;type:VARCHAR(256);NOT NULL"`
	Desc         string            `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`     //应用描述
	IpRange      []string          `gorm:"column:ip_range;type:json;serializer:json;"` //ip白名单
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:tc_ac"`
}

func (m *SysDataOpenAccess) TableName() string {
	return "sys_data_open_access"
}

// 用户区域权限表
type SysDataArea struct {
	ID             int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode     stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 租户编码
	TargetType     def.TargetType    `gorm:"column:target_type;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"`
	TargetID       int64             `gorm:"column:target_id;uniqueIndex:ri_mi;type:bigint;NOT NULL"`          // 授权对象的id,角色id,用户id
	ProjectID      stores.ProjectID  `gorm:"column:project_id;uniqueIndex:ri_mi;type:bigint;NOT NULL"`         // 所属项目ID(雪花ID)
	AreaID         int64             `gorm:"column:area_id;uniqueIndex:ri_mi;type:bigint;NOT NULL"`            // 区域ID(雪花ID)
	AreaIDPath     string            `gorm:"column:area_id_path;uniqueIndex:ri_mi;type:VARCHAR(256);NOT NULL"` // 区域ID(雪花ID)
	AuthType       def.AuthType      `gorm:"column:auth_type;type:bigint;NOT NULL"`                            // 授权类型 1 管理员(可以调整本区域及旗下区域的设备区域规划)  2 读写授权(可以对该区域及旗下区域的设备进行管理) 3 只读授权()
	IsAuthChildren int64             `gorm:"column:is_auth_children;type:bigint;default:2;NOT NULL"`           //是否同时授权子节点
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:ri_mi"`
}

func (m *SysDataArea) TableName() string {
	return "sys_data_area"
}

// 用户区域权限授权表
type SysUserAreaApply struct {
	ID         int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 租户编码
	UserID     int64             `gorm:"column:user_id;uniqueIndex:ri_mi;type:bigint;NOT NULL"`          // 用户ID(雪花id)
	ProjectID  stores.ProjectID  `gorm:"column:project_id;uniqueIndex:ri_mi;type:bigint;NOT NULL"`       // 所属项目ID(雪花ID)
	AreaID     stores.AreaID     `gorm:"column:area_id;uniqueIndex:ri_mi;type:bigint;NOT NULL"`          // 区域ID(雪花ID)
	AuthType   def.AuthType      `gorm:"column:auth_type;type:bigint;NOT NULL"`                          // 授权类型 1 管理员(可以调整本区域及旗下区域的设备区域规划)  2 读授权(可以对区域下的设备进行操作,但是不能修改区域) 2 读写授权(可以对区域下的设备进行操作,同时可以对区域进行修改,但是不能新增或删除)
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:ri_mi"`
}

func (m *SysUserAreaApply) TableName() string {
	return "sys_user_area_apply"
}

// 用户项目权限表
type SysDataProject struct {
	ID         int64             `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);default:default"` // 租户编码
	ProjectID  int64             `gorm:"column:project_id;uniqueIndex:ri_mi;type:bigint;NOT NULL"`              // 所属项目ID(雪花ID)
	TargetType def.TargetType    `gorm:"column:target_type;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"`
	TargetID   int64             `gorm:"column:target_id;uniqueIndex:ri_mi;type:bigint;NOT NULL"` // 授权对象的id,角色id,用户id
	AuthType   def.AuthType      `gorm:"column:auth_type;type:bigint;NOT NULL"`                   // 授权类型 1 管理员(可以修改本项目的状态,同时拥有所有区域权限)  2 读授权(可以对项目下的区域进行操作,但是不能修改项目) 2 读写授权(可以对项目下的区域进行操作,同时可以对项目进行修改)
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:ri_mi"`
}

func (m *SysDataProject) TableName() string {
	return "sys_data_project"
}
