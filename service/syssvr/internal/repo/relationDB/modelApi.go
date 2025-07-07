package relationDB

import (
	"gitee.com/unitedrhino/core/service/syssvr/internal/domain/access"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/stores"
)

// 功能权限范围
type SysAccessInfo struct {
	ID         int64           `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                 // 编号
	Name       string          `gorm:"column:name;type:VARCHAR(100);NOT NULL"`                           // 请求名称
	Module     string          `gorm:"column:module;type:VARCHAR(100);default:'系统管理'"`               //所属模块
	Code       string          `gorm:"column:code;type:VARCHAR(100);uniqueIndex:idx_app_route;NOT NULL"` // 请求名称
	Group      string          `gorm:"column:group;type:VARCHAR(100);NOT NULL"`                          // 接口组
	IsNeedAuth int64           `gorm:"column:is_need_auth;type:BIGINT;default:1;NOT NULL"`               // 是否需要认证（1是 2否）
	AuthType   access.AuthType `gorm:"column:is_auth_tenant;type:BIGINT;default:1;NOT NULL"`             // 1(all) 全部人可以操作 2(admin) 默认授予租户管理员权限 3(superAdmin,supper) default租户才可以操作(超管是跨租户的)
	Desc       string          `gorm:"column:desc;type:VARCHAR(500);NOT NULL"`                           // 备注
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_app_route"`
	Apis        []*SysApiInfo      `gorm:"foreignKey:AccessCode;references:Code"`
}

func (m *SysAccessInfo) TableName() string {
	return "sys_access_info"
}

// 接口管理
type SysApiInfo struct {
	ID            int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`              // 编号
	AccessCode    string `gorm:"column:access_code;type:VARCHAR(50);NOT NULL"`                  // 范围编码
	Method        string `gorm:"column:method;uniqueIndex:idx_route;type:VARCHAR(50);NOT NULL"` // 请求方式（1 GET 2 POST 3 HEAD 4 OPTIONS 5 PUT 6 DELETE 7 TRACE 8 CONNECT 9 其它）
	Route         string `gorm:"column:route;uniqueIndex:idx_route;type:VARCHAR(100);NOT NULL"` // 路由
	Name          string `gorm:"column:name;type:VARCHAR(100);NOT NULL"`                        // 请求名称
	BusinessType  int64  `gorm:"column:business_type;type:BIGINT;NOT NULL"`                     // 业务类型（1(add)新增 2修改(modify) 3删除(delete) 4查询(find) 5其它(other)
	RecordLogMode int64  `gorm:"column:record_log_mode;type:BIGINT;default:1;"`                 //1为自动模式(读取类型忽略,其他类型记录日志) 2全部记录 3不记录
	Desc          string `gorm:"column:desc;type:VARCHAR(500);NOT NULL"`                        // 备注
	//AuthType     int64  `gorm:"column:is_auth_tenant;type:BIGINT;default:1;NOT NULL"`      // 1(all) 全部人可以操作 2(admin) 默认授予租户管理员权限 3(superAdmin,supper) default租户才可以操作(超管是跨租户的)
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_route"`
	Access      *SysAccessInfo     `gorm:"foreignKey:Code;references:AccessCode"`
}

func (m *SysApiInfo) TableName() string {
	return "sys_api_info"
}

// 应用菜单关联表
type SysTenantAccess struct {
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                           // id编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;uniqueIndex:idx_tenant_scope;type:VARCHAR(50);NOT NULL;"` // 租户编码
	AccessCode string              `gorm:"column:access_code;uniqueIndex:idx_tenant_scope;type:VARCHAR(50);NOT NULL"`  // 范围编码
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_tenant_scope"`
}

func (m *SysTenantAccess) TableName() string {
	return "sys_tenant_access"
}
