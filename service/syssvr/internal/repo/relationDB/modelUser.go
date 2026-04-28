package relationDB

import (
	"database/sql"
	"time"

	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
)

// 用户登录信息表
type SysUserInfo struct {
	TenantCode      dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;uniqueIndex:idx_sys_user_info_tc_un;uniqueIndex:idx_sys_user_info_tc_doi;uniqueIndex:idx_sys_user_info_tc_email;uniqueIndex:idx_sys_user_info_huawei_union;uniqueIndex:idx_sys_user_info_huawei_open;uniqueIndex:idx_sys_user_info_tc_phone;uniqueIndex:idx_sys_user_info_tc_wui;uniqueIndex:idx_sys_user_info_tc_woi"` // 租户编码
	UserID          int64               `gorm:"column:user_id;primary_key;AUTO_INCREMENT;type:BIGINT;NOT NULL"`                                                                                                                                                                                                                // 用户id
	UserName        sql.NullString      `gorm:"column:user_name;uniqueIndex:idx_sys_user_info_tc_un;type:VARCHAR(200)"`                                                                                                                                                                                                        // 登录用户名
	NickName        string              `gorm:"column:nick_name;type:VARCHAR(60);NOT NULL"`                                                                                                                                                                                                                                    // 用户的昵称
	Password        string              `gorm:"column:password;type:CHAR(32);NOT NULL"`                                                                                                                                                                                                                                        // 登录密码
	Email           sql.NullString      `gorm:"column:email;uniqueIndex:idx_sys_user_info_tc_email;type:VARCHAR(255)"`                                                                                                                                                                                                         // 邮箱
	Phone           sql.NullString      `gorm:"column:phone;uniqueIndex:idx_sys_user_info_tc_phone;type:VARCHAR(20)"`                                                                                                                                                                                                          // 手机号
	WechatUnionID   sql.NullString      `gorm:"column:wechat_union_id;uniqueIndex:idx_sys_user_info_tc_wui;type:VARCHAR(128)"`                                                                                                                                                                                                 // 微信union id
	WechatOpenID    sql.NullString      `gorm:"column:wechat_open_id;uniqueIndex:idx_sys_user_info_tc_woi;type:VARCHAR(128)"`                                                                                                                                                                                                  // 微信union id
	DingTalkUserID  sql.NullString      `gorm:"column:ding_talk_user_id;uniqueIndex:idx_sys_user_info_tc_doi;type:VARCHAR(128)"`
	DingTalkUnionID sql.NullString      `gorm:"column:ding_talk_union_id;uniqueIndex:idx_sys_user_info_tc_doi;type:VARCHAR(128)"`
	HuaweiUnionID   sql.NullString      `gorm:"column:huawei_union_id;uniqueIndex:idx_sys_user_info_huawei_union;type:VARCHAR(128)"`                                                                                                                                                                                                 // 微信union id
	HuaweiOpenID    sql.NullString      `gorm:"column:huawei_open_id;uniqueIndex:idx_sys_user_info_huawei_open;type:VARCHAR(128)"`                                                                                                                                                                                                  // 微信union id

	LastIP          string              `gorm:"column:last_ip;type:VARCHAR(128);NOT NULL"`                       // 最后登录ip
	LastTokenID     string              `gorm:"column:last_token_id;type:VARCHAR(128);default:''"`               // 最后登录的token ID
	RegIP           string              `gorm:"column:reg_ip;type:VARCHAR(128);NOT NULL"`                        // 注册ip
	Sex             int64               `gorm:"column:sex;type:SMALLINT;default:3;NOT NULL"`                     // 用户的性别，值为1时是男性，值为2时是女性，其他值为未知
	City            string              `gorm:"column:city;type:VARCHAR(50);NOT NULL"`                           // 用户所在城市
	Country         string              `gorm:"column:country;type:VARCHAR(50);NOT NULL"`                        // 用户所在国家
	Province        string              `gorm:"column:province;type:VARCHAR(50);NOT NULL"`                       // 用户所在省份
	Language        string              `gorm:"column:language;type:VARCHAR(50);NOT NULL"`                       // 用户的语言，简体中文为zh_CN
	HeadImg         string              `gorm:"column:head_img;type:VARCHAR(256);NOT NULL"`                      // 用户头像
	Role            int64               `gorm:"column:role;type:BIGINT;NOT NULL"`                                // 用户默认角色（默认使用该角色）
	Tags            map[string]string   `gorm:"column:tags;type:json;serializer:json;NOT NULL"`     // 私有标签,只有管理员可以修改
	PubTags         map[string]string   `gorm:"column:pub_tags;type:json;serializer:json;NOT NULL"` // 公共的标签,用户自己可以修改
	IsAllData       int64               `gorm:"column:is_all_data;type:SMALLINT;default:1;NOT NULL"`             // 是否所有数据权限（1是，2否）
	DeviceCount     int64               `gorm:"column:device_count;default:0"`                                   //用户所拥有的设备数量统计
	Roles           []*SysUserRole      `gorm:"foreignKey:UserID;references:UserID"`
	Tenant          *SysTenantInfo      `gorm:"foreignKey:Code;references:TenantCode"`
	Status          int64               `gorm:"column:status;type:BIGINT;NOT NULL;default:1"` //租戶状态: 1启用 2禁用
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_user_info_tc_un;uniqueIndex:idx_sys_user_info_huawei_union;uniqueIndex:idx_sys_user_info_huawei_open;uniqueIndex:idx_sys_user_info_tc_doi;uniqueIndex:idx_sys_user_info_tc_email;uniqueIndex:idx_sys_user_info_tc_phone;uniqueIndex:idx_sys_user_info_tc_wui;uniqueIndex:idx_sys_user_info_tc_woi"`
}

func (m *SysUserInfo) TableName() string {
	return "sys_user_info"
}

func (m *SysUserInfo) BeforeSave(tx *gorm.DB) error {
	if m == nil {
		return nil
	}
	m.Tags = normalizeJSONMapStringString(m.Tags)
	m.PubTags = normalizeJSONMapStringString(m.PubTags)
	return nil
}

// 应用菜单关联表
type SysUserRole struct {
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                        // id编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL;"`                     // 租户编码
	UserID     int64               `gorm:"column:user_id;uniqueIndex:idx_sys_user_role_ri_mi;NOT NULL;type:BIGINT"` // 用户ID
	RoleID     int64               `gorm:"column:role_id;uniqueIndex:idx_sys_user_role_ri_mi;NOT NULL;type:BIGINT"` // 角色ID
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
	CreatedTime   time.Time           `gorm:"column:created_time;index;autoCreateTime"` // 登录时间
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
	CreatedTime  time.Time           `gorm:"column:created_time;index;autoCreateTime"` // 操作时间
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
