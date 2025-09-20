package relationDB

import (
	"database/sql"
	"time"

	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
)

// 用户登录信息表
type SysUserInfo struct {
	UserID      int64            `gorm:"column:user_id;primary_key;AUTO_INCREMENT;type:BIGINT;NOT NULL"`        // 用户id
	UserName    sql.NullString   `gorm:"column:user_name;uniqueIndex:idx_sys_user_info_tc_un;type:VARCHAR(20)"` // 登录用户名
	NickName    string           `gorm:"column:nick_name;type:VARCHAR(60);NOT NULL"`                            // 用户的昵称
	Password    string           `gorm:"column:password;type:CHAR(32);NOT NULL"`                                // 登录密码
	Email       sql.NullString   `gorm:"column:email;uniqueIndex:idx_sys_user_info_tc_email;type:VARCHAR(255)"` // 邮箱
	Phone       sql.NullString   `gorm:"column:phone;uniqueIndex:idx_sys_user_info_tc_phone;type:VARCHAR(20)"`  // 手机号
	LastIP      string           `gorm:"column:last_ip;type:VARCHAR(128);NOT NULL"`                             // 最后登录ip
	LastTokenID string           `gorm:"column:last_token_id;type:VARCHAR(128);default:''"`                     // 最后登录的token ID
	RegIP       string           `gorm:"column:reg_ip;type:VARCHAR(128);NOT NULL"`                              // 注册ip
	Sex         int64            `gorm:"column:sex;type:SMALLINT;default:3;NOT NULL"`                           // 用户的性别，值为1时是男性，值为2时是女性，其他值为未知
	City        string           `gorm:"column:city;type:VARCHAR(50);NOT NULL"`                                 // 用户所在城市
	Country     string           `gorm:"column:country;type:VARCHAR(50);NOT NULL"`                              // 用户所在国家
	Province    string           `gorm:"column:province;type:VARCHAR(50);NOT NULL"`                             // 用户所在省份
	Language    string           `gorm:"column:language;type:VARCHAR(50);NOT NULL"`                             // 用户的语言，简体中文为zh_CN
	HeadImg     string           `gorm:"column:head_img;type:VARCHAR(256);NOT NULL"`                            // 用户头像
	Tenants     []*SysUserTenant `gorm:"foreignKey:UserID;references:UserID"`
	Thirds      []*SysUserThird  `gorm:"foreignKey:UserID;references:UserID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_user_info_tc_un;uniqueIndex:idx_sys_user_info_tc_doi;uniqueIndex:idx_sys_user_info_tc_email;uniqueIndex:idx_sys_user_info_tc_phone;uniqueIndex:idx_sys_user_info_tc_wui;"`
}

func (m *SysUserInfo) TableName() string {
	return "sys_user_info"
}

type SysUserTenant struct {
	TenantCode   dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);uniqueIndex:sys_user_tenant_user_tenant;NOT NULL;"` // 租户编码
	UserID       int64               `gorm:"column:user_id;uniqueIndex:sys_user_tenant_user_tenant;type:BIGINT;NOT NULL"`           // 用户id
	DeviceCount  int64               `gorm:"column:device_count;default:0"`                                                         //用户所拥有的设备数量统计
	Roles        []*SysUserRole      `gorm:"foreignKey:UserID;references:UserID"`
	Status       int64               `gorm:"column:status;type:BIGINT;NOT NULL;default:1"`                //用户状态: 1启用 2禁用
	Tags         map[string]string   `gorm:"column:tags;type:json;serializer:json;NOT NULL;default:'{}'"` // 产品标签
	User         *SysUserInfo        `gorm:"foreignKey:UserID;references:UserID"`
	TenantInfo   *SysTenantInfo      `gorm:"foreignKey:Code;references:TenantCode"`
	TenantConfig *SysTenantConfig    `gorm:"foreignKey:TenantCode;references:TenantCode"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:sys_user_tenant_user_tenant;"`
}

func (m *SysUserTenant) TableName() string {
	return "sys_user_tenant"
}

type SysUserThird struct {
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;"` // 租户编码,如果是公共应用登录的,这里填写__common__
	AppType    def.ThirdType       `gorm:"column:app_type;uniqueIndex:idx_sys_user_info_tc_wui;type:varchar(64);NOT NULL"`
	AppID      string              `gorm:"column:app_id;uniqueIndex:idx_sys_user_info_tc_wui;type:VARCHAR(128);NOT NULL"`
	UserID     int64               `gorm:"column:user_id;type:BIGINT;NOT NULL"`                                                    // 用户id
	UnionID    string              `gorm:"column:union_id;index;type:VARCHAR(128);default:''"`                                     // 微信union id
	OpenID     string              `gorm:"column:open_id;uniqueIndex:idx_sys_user_info_tc_wui;index;type:VARCHAR(128);default:''"` // 钉钉里是UserID
	User       *SysUserInfo        `gorm:"foreignKey:UserID;references:UserID"`
	//UserTenants []*SysUserTenant             `gorm:"foreignKey:UserID;references:UserID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_user_info_tc_wui;"`
}

func (m *SysUserThird) TableName() string {
	return "sys_user_third"
}

// 应用菜单关联表
type SysUserRole struct {
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                        // id编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;uniqueIndex:idx_sys_user_role_ri_mi;index;type:VARCHAR(50);NOT NULL;"` // 租户编码
	UserID     int64               `gorm:"column:user_id;uniqueIndex:idx_sys_user_role_ri_mi;NOT NULL;type:BIGINT"`                 // 用户ID
	RoleID     int64               `gorm:"column:role_id;uniqueIndex:idx_sys_user_role_ri_mi;NOT NULL;type:BIGINT"`                 // 角色ID
	Role       *SysRoleInfo        `gorm:"foreignKey:ID;references:RoleID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_user_role_ri_mi"`
}

func (m *SysUserRole) TableName() string {
	return "sys_user_role"
}

// 登录日志管理
type SysLoginLog struct {
	ID            int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`             // 编号
	TenantCode    dataType.TenantCode `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`           // 租户编码
	AppCode       string              `gorm:"column:app_code;NOT NULL;type:VARCHAR(50)"`                    // 应用ID
	UserID        int64               `gorm:"column:user_id;type:BIGINT;NOT NULL"`                          // 用户id
	UserName      string              `gorm:"column:user_name;type:VARCHAR(50)"`                            // 登录账号
	IpAddr        string              `gorm:"column:ip_addr;type:VARCHAR(50)"`                              // 登录IP地址
	LoginLocation string              `gorm:"column:login_location;type:VARCHAR(100)"`                      // 登录地点
	Browser       string              `gorm:"column:browser;type:VARCHAR(50)"`                              // 浏览器类型
	Os            string              `gorm:"column:os;type:VARCHAR(50)"`                                   // 操作系统
	Code          int64               `gorm:"column:code;type:BIGINT;default:200;NOT NULL"`                 // 登录状态（200成功 其它失败）
	Msg           string              `gorm:"column:msg;type:VARCHAR(255)"`                                 // 提示消息
	CreatedTime   time.Time           `gorm:"column:created_time;index;default:CURRENT_TIMESTAMP;NOT NULL"` // 登录时间
}

func (m *SysLoginLog) TableName() string {
	return "sys_login_log"
}

// 操作日志管理
type SysOperLog struct {
	ID           int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`             // 编号
	TenantCode   dataType.TenantCode `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`           // 租户编码
	AppCode      string              `gorm:"column:app_code;NOT NULL;type:VARCHAR(50)"`                    // 应用ID
	OperUserID   int64               `gorm:"column:oper_user_id;type:BIGINT;NOT NULL"`                     // 用户id
	OperUserName string              `gorm:"column:oper_user_name;type:VARCHAR(50)"`                       // 操作人员名称
	OperName     string              `gorm:"column:oper_name;type:VARCHAR(50)"`                            // 操作名称
	BusinessType int64               `gorm:"column:business_type;type:BIGINT;NOT NULL"`                    // 业务类型（1新增 2修改 3删除 4查询 5其它）
	Uri          string              `gorm:"column:uri;type:VARCHAR(100)"`                                 // 请求地址
	OperIpAddr   string              `gorm:"column:oper_ip_addr;type:VARCHAR(50)"`                         // 主机地址
	OperLocation string              `gorm:"column:oper_location;type:VARCHAR(255)"`                       // 操作地点
	Req          sql.NullString      `gorm:"column:req;type:TEXT"`                                         // 请求参数
	Resp         sql.NullString      `gorm:"column:resp;type:TEXT"`                                        // 返回参数
	Code         int64               `gorm:"column:code;type:BIGINT;default:200;NOT NULL"`                 // 返回状态（200成功 其它失败）
	Msg          string              `gorm:"column:msg;type:VARCHAR(255)"`                                 // 提示消息
	CreatedTime  time.Time           `gorm:"column:created_time;index;default:CURRENT_TIMESTAMP;NOT NULL"` // 操作时间
}

func (m *SysOperLog) TableName() string {
	return "sys_oper_log"
}

// 用户配置表
type SysUserProfile struct {
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                     // 编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;uniqueIndex:idx_sys_user_profile_tc_un;"` // 租户编码
	UserID     int64               `gorm:"column:user_id;uniqueIndex:idx_sys_user_profile_tc_un;type:BIGINT;NOT NULL"`           // 用户id
	Code       string              `gorm:"column:code;type:VARCHAR(50);uniqueIndex:idx_sys_user_profile_tc_un;NOT NULL"`         //配置code
	Params     string              `gorm:"column:params;type:text;NOT NULL"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_user_profile_tc_un;"`
}

func (m *SysUserProfile) TableName() string {
	return "sys_user_profile"
}
