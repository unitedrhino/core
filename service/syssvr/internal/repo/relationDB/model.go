package relationDB

import (
	"database/sql"
	"gitee.com/unitedrhino/core/service/syssvr/internal/domain/dept"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"time"
)

// 示例
type SysExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

// 全局配置
type SysConfig struct {
	Sms *SysConfigSms `gorm:"embedded;embeddedPrefix:sms_"` //短信配置,全租户共用
}

type SysConfigSms struct {
	From     string `gorm:"column:from;type:VARCHAR(50);default:'';NOT NULL"`     // 发件人  你自己要发邮件的邮箱
	Host     string `gorm:"column:host;type:VARCHAR(50);default:'';NOT NULL"`     // 服务器地址 例如 smtp.qq.com  请前往QQ或者你要发邮件的邮箱查看其smtp协议
	Secret   string `gorm:"column:secret;type:VARCHAR(50);default:'';NOT NULL"`   // 密钥    用于登录的密钥 最好不要用邮箱密码 去邮箱smtp申请一个用于登录的密钥
	Nickname string `gorm:"column:nickname;type:VARCHAR(50);default:'';NOT NULL"` // 昵称    发件人昵称 通常为自己的邮箱
	Port     int64  `gorm:"column:port;type:int;default:465"`                     // 端口     请前往QQ或者你要发邮件的邮箱查看其smtp协议 大多为 465
	IsSSL    int64  `gorm:"column:is_ssl;type:int;default:2"`                     // 是否SSL   是否开启SSL
}

type SysDictInfo struct {
	ID         int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                  // id编号
	Name       string `gorm:"column:name;uniqueIndex:name;comment:字典名"`                          // 字典名（中）
	Code       string `gorm:"column:code;uniqueIndex:code;type:VARCHAR(50);default:'';NOT NULL"` //编码
	Group      string `gorm:"column:group;type:VARCHAR(50);default:'';NOT NULL"`                 //字典分组
	Desc       string `gorm:"column:desc;comment:描述"`                                            // 描述
	Body       string `gorm:"column:body;type:VARCHAR(1024)"`                                    // 自定义数据
	StructType int64  `gorm:"column:struct_type;type:BIGINT;default:1"`                          //结构类型(不可修改) 1:列表(默认) 2:树型
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:code;uniqueIndex:name"`
	Details     []*SysDictDetail   `gorm:"foreignKey:DictCode;references:Code"`
}

func (SysDictInfo) TableName() string {
	return "sys_dict_info"
}

type SysDictDetail struct {
	ID       int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                        // id编号
	DictCode string `gorm:"column:dict_code;uniqueIndex:value;type:VARCHAR(50);default:'';NOT NULL"` // 关联标记
	Label    string `gorm:"column:label;comment:展示值"`                                                // 展示值
	Value    string `gorm:"column:value;uniqueIndex:value;comment:字典值"`                              // 字典值
	Status   int64  `gorm:"column:status;type:SMALLINT;default:1"`                                   // 状态  1:启用,2:禁用
	Sort     int64  `gorm:"column:sort;comment:排序标记"`                                                // 排序标记
	Desc     string `gorm:"column:desc;comment:描述"`                                                  // 描述
	Body     string `gorm:"column:body;type:VARCHAR(1024)"`                                          // 自定义数据
	IDPath   string `gorm:"column:id_path;type:varchar(100);NOT NULL"`                               // 1-2-3-的格式记录顶级区域到当前id的路径
	ParentID int64  `gorm:"column:parent_id;type:BIGINT"`                                            // id编号
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:value"`
	Children    []*SysDictDetail   `gorm:"foreignKey:parent_id;references:id"`
	Parent      *SysDictDetail     `gorm:"foreignKey:ID;references:ParentID"`
}

func (SysDictDetail) TableName() string {
	return "sys_dict_detail"
}

type SysDeptInfo struct {
	ID             int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                   // id编号
	TenantCode     dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`                       // 租户编码
	ParentID       int64               `gorm:"column:parent_id;uniqueIndex:name;type:BIGINT"`                      // id编号
	Name           string              `gorm:"column:name;type:VARCHAR(256);uniqueIndex:name;default:'';NOT NULL"` // 部门名称
	AdminUserID    int64               `gorm:"column:admin_user_id;comment:管理员账号;NOT NULL"`
	Status         int64               `gorm:"column:status;type:SMALLINT;default:1"` // 状态  1:启用,2:禁用
	Sort           int64               `gorm:"column:sort;comment:排序标记"`              // 排序标记
	Desc           string              `gorm:"column:desc;comment:描述"`                // 描述
	UserCount      int64               `gorm:"column:user_count;comment:用户统计,包含下级部门的人数"`
	DeviceCount    int64               `gorm:"column:device_count;default:0;comment:部门自己的设备总数"`
	AllDeviceCount int64               `gorm:"column:all_device_count;default:0;comment:部门及其下级的设备总数"`
	ChildrenCount  int64               `gorm:"column:children_count;default:0;comment:部门下级数量"`
	IDPath         string              `gorm:"column:id_path;type:varchar(100);NOT NULL"` // 1-2-3-的格式记录顶级区域到当前id的路径
	DingTalkID     int64               `gorm:"column:ding_talk_id;default:0;"`            //钉钉的部门ID
	Tags           map[string]string   `gorm:"column:tags;type:json;serializer:json"`     //部门标签
	AdminUser      *SysUserInfo        `gorm:"foreignKey:user_id;references:AdminUserID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:name"`
	Children    []*SysDeptInfo     `gorm:"foreignKey:parent_id;references:id"`
	Parent      *SysDeptInfo       `gorm:"foreignKey:ID;references:ParentID"`
}

func (SysDeptInfo) TableName() string {
	return "sys_dept_info"
}

type SysDeptUser struct {
	ID             int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`       // id编号
	TenantCode     dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;"`          // 租户编码
	UserID         int64               `gorm:"column:user_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"`  // 用户ID
	DeptID         int64               `gorm:"column:dept_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"`  // 角色ID
	DeptIDPath     string              `gorm:"column:dept_id_path;type:varchar(100);NOT NULL"`         // 1-2-3-的格式记录顶级区域到当前id的路径
	AuthType       def.AuthType        `gorm:"column:auth_type;type:bigint;NOT NULL"`                  // 授权类型 1 管理员(可以调整本区域及旗下区域的设备区域规划)  2 读写授权(可以对该区域及旗下区域的设备进行管理) 3 只读授权()
	IsAuthChildren int64               `gorm:"column:is_auth_children;type:bigint;default:1;NOT NULL"` //是否同时授权子节点
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:ri_mi"`
	Dept        *SysDeptInfo       `gorm:"foreignKey:ID;references:DeptID"`
	User        *SysUserInfo       `gorm:"foreignKey:UserID;references:UserID"`
}

func (m *SysDeptUser) TableName() string {
	return "sys_dept_user"
}

// 租户下的应用列表
type SysDeptSyncJob struct {
	ID          int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	TenantCode  dataType.TenantCode `gorm:"column:tenant_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 租户编码
	Direction   dept.SyncDirection  `gorm:"column:direction;default:1"`                                     // 同步的方向,1上游同步到联犀(默认),2联犀同步到下游
	ThirdType   def.AppSubType      `gorm:"column:third_type;type:varchar(20)"`                             //同步的类型
	ThirdConfig *SysTenantThird     `gorm:"embedded;embeddedPrefix:third_config"`                           //第三方配置
	FieldMap    map[string]string   `gorm:"column:field_map;type:json;serializer:json"`                     //用户字段映射,左边是联犀的字段,右边是第三方的,不填写就是全量映射
	SyncDeptIDs []int64             `gorm:"column:sync_dept_ids;type:json;serializer:json"`                 //同步的第三方部门id列表,不填为同步全部
	IsAddSync   int64               `gorm:"column:is_add_sync;default:1"`                                   //新增人员自动同步,默认为1
	SyncMode    dept.SyncMode       `gorm:"column:sync_mode;default:1"`                                     //同步模式: 1:手动(默认) 2: 定时同步(半小时) 3: 实时同步

	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:tc_ac"`
}

func (m *SysDeptSyncJob) TableName() string {
	return "sys_dept_sync_job"
}

type SysSlotInfo struct {
	ID       int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                        // id编号
	Code     string            `gorm:"column:code;uniqueIndex:code_slot;type:VARCHAR(100);NOT NULL"`            // 鉴权的编码
	SubCode  string            `gorm:"column:sub_code;uniqueIndex:code_slot;type:VARCHAR(100);NOT NULL"`        // 鉴权的编码
	SlotCode string            `gorm:"column:slot_code;uniqueIndex:code_slot;type:VARCHAR(100);NOT NULL"`       //slot的编码
	Method   string            `gorm:"column:method;type:VARCHAR(50);default:'POST'"`                           // 请求方式 GET  POST
	Uri      string            `gorm:"column:uri;type:VARCHAR(100);NOT NULL"`                                   // 参考: /api/v1/system/user/self/captcha?fwefwf=gwgweg&wefaef=gwegwe
	Hosts    []string          `gorm:"column:hosts;type:json;serializer:json;NOT NULL;default:'[]';NOT NULL"`   //访问的地址 host or host:port
	Body     string            `gorm:"column:body;type:VARCHAR(100);default:''"`                                // body 参数模板
	Handler  map[string]string `gorm:"column:handler;type:json;serializer:json;NOT NULL;default:'{}';NOT NULL"` //http头
	AuthType string            `gorm:"column:auth_type;type:VARCHAR(100);NOT NULL"`                             //鉴权类型 core
	Desc     string            `gorm:"column:desc;type:VARCHAR(500);"`                                          // 备注
	stores.SoftTime
}

func (m *SysSlotInfo) TableName() string {
	return "sys_slot_info"
}

type SysServiceInfo struct {
	ID      int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"` // id编号
	Code    string `gorm:"column:code;unique;type:VARCHAR(100);NOT NULL"`    // 服务编码
	Name    string `gorm:"column:name;type:VARCHAR(100);NOT NULL"`           // 服务名
	Version string `gorm:"column:version;type:VARCHAR(100);NOT NULL"`        //服务版本
	Desc    string `gorm:"column:desc;type:VARCHAR(500);"`                   // 备注
	stores.NoDelTime
}

func (m *SysServiceInfo) TableName() string {
	return "sys_service_info"
}

// 应用信息
type SysAppInfo struct {
	ID      int64          `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`        // id编号
	Code    string         `gorm:"column:code;uniqueIndex:code;type:VARCHAR(100);NOT NULL"` // 应用编码
	Name    string         `gorm:"column:name;uniqueIndex:name;type:VARCHAR(100);NOT NULL"` //应用名称
	Type    def.AppType    `gorm:"column:type;type:VARCHAR(100);default:web;NOT NULL"`      //应用类型 web:web页面  app:应用  mini:小程序
	SubType def.AppSubType `gorm:"column:sub_type;type:VARCHAR(100);default:wx;NOT NULL"`   // 类型  wx:微信小程序  ding:钉钉小程序
	Desc    string         `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`                  //应用描述
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:code;uniqueIndex:name"`
}

func (m *SysAppInfo) TableName() string {
	return "sys_app_info"
}

// 应用默认绑定的模块
type SysAppModule struct {
	ID         int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`               // id编号
	AppCode    string `gorm:"column:app_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"`    // 应用编码 这里只关联主应用,主应用授权,子应用也授权了
	ModuleCode string `gorm:"column:module_code;uniqueIndex:tc_ac;type:VARCHAR(50);NOT NULL"` // 模块编码
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:tc_ac"`
	Module      *SysModuleInfo     `gorm:"foreignKey:Code;references:ModuleCode"`
	App         *SysAppInfo        `gorm:"foreignKey:Code;references:AppCode"`
}

func (m *SysAppModule) TableName() string {
	return "sys_app_module"
}

// 模块管理表 模块是菜单和接口的集合体
type SysModuleInfo struct {
	ID         int64            `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`       // 编号
	Code       string           `gorm:"column:code;uniqueIndex:code;NOT NULL;type:VARCHAR(50)"` // 编码
	Type       int64            `gorm:"column:type;type:BIGINT;default:1;NOT NULL"`             // 类型   1:web页面  2:应用  3:小程序
	SubType    int64            `gorm:"column:sub_type;type:BIGINT;default:1;NOT NULL"`         // 类型   1：微应用   2：iframe内嵌 3: 原生菜单
	Order      int64            `gorm:"column:order;type:BIGINT;default:1;NOT NULL"`            // 左侧table排序序号
	Name       string           `gorm:"column:name;type:VARCHAR(50);NOT NULL"`                  // 菜单名称
	Path       string           `gorm:"column:path;type:VARCHAR(64);NOT NULL"`                  // 系统的path
	Url        string           `gorm:"column:url;type:VARCHAR(200);NOT NULL"`                  // 页面
	Icon       string           `gorm:"column:icon;type:VARCHAR(64);NOT NULL"`                  // 图标
	Body       string           `gorm:"column:body;type:VARCHAR(1024)"`                         // 菜单自定义数据
	HideInMenu int64            `gorm:"column:hide_in_menu;type:BIGINT;default:2;NOT NULL"`     // 是否隐藏菜单 1-是 2-否
	Desc       string           `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`                 // 备注
	Tag        int64            `gorm:"column:tag;type:BIGINT;default:1;NOT NULL"`              //标签: 1:通用 2:选配
	Menus      []*SysModuleMenu `gorm:"foreignKey:ModuleCode;references:Code"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:code"`
}

func (m *SysModuleInfo) TableName() string {
	return "sys_module_info"
}

// 菜单管理表
type SysModuleMenu struct {
	ID         int64            `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                      // 编号
	ModuleCode string           `gorm:"column:module_code;uniqueIndex:menu_path;type:VARCHAR(50);NOT NULL"`    // 模块编码
	ParentID   int64            `gorm:"column:parent_id;uniqueIndex:menu_path;type:BIGINT;default:1;NOT NULL"` // 父菜单ID，一级菜单为1
	Type       int64            `gorm:"column:type;type:BIGINT;default:1;NOT NULL"`                            // 类型   1：菜单或者页面   2：iframe嵌入   3：外链跳转
	Order      int64            `gorm:"column:order;type:BIGINT;default:1;NOT NULL"`                           // 左侧table排序序号
	Name       string           `gorm:"column:name;type:VARCHAR(50);NOT NULL"`                                 // 菜单名称
	Path       string           `gorm:"column:path;uniqueIndex:menu_path;type:VARCHAR(64);NOT NULL"`           // 系统的path
	Component  string           `gorm:"column:component;type:VARCHAR(1024);NOT NULL"`                          // 页面
	Icon       string           `gorm:"column:icon;type:VARCHAR(64);NOT NULL"`                                 // 图标
	Redirect   string           `gorm:"column:redirect;type:VARCHAR(64);NOT NULL"`                             // 路由重定向
	Body       string           `gorm:"column:body;type:VARCHAR(1024)"`                                        // 菜单自定义数据
	HideInMenu int64            `gorm:"column:hide_in_menu;type:BIGINT;default:2;NOT NULL"`                    // 是否隐藏菜单 1-是 2-否
	IsCommon   int64            `gorm:"column:is_common;type:BIGINT;default:2;"`                               // 是否常用菜单 1-是 2-否
	Children   []*SysModuleMenu `gorm:"foreignKey:ID;references:ParentID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:menu_path;default:0;index"`
}

func (m *SysModuleMenu) TableName() string {
	return "sys_module_menu"
}

// 用户登录信息表
type SysUserInfo struct {
	TenantCode      dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;uniqueIndex:tc_un;uniqueIndex:tc_doi;uniqueIndex:tc_email;uniqueIndex:tc_phone;uniqueIndex:tc_wui;uniqueIndex:tc_woi"` // 租户编码
	UserID          int64               `gorm:"column:user_id;primary_key;AUTO_INCREMENT;type:BIGINT;NOT NULL"`                                                                                                    // 用户id
	UserName        sql.NullString      `gorm:"column:user_name;uniqueIndex:tc_un;type:VARCHAR(20)"`                                                                                                               // 登录用户名
	NickName        string              `gorm:"column:nick_name;type:VARCHAR(60);NOT NULL"`                                                                                                                        // 用户的昵称
	Password        string              `gorm:"column:password;type:CHAR(32);NOT NULL"`                                                                                                                            // 登录密码
	Email           sql.NullString      `gorm:"column:email;uniqueIndex:tc_email;type:VARCHAR(255)"`                                                                                                               // 邮箱
	Phone           sql.NullString      `gorm:"column:phone;uniqueIndex:tc_phone;type:VARCHAR(20)"`                                                                                                                // 手机号
	WechatUnionID   sql.NullString      `gorm:"column:wechat_union_id;uniqueIndex:tc_wui;type:VARCHAR(128)"`                                                                                                       // 微信union id
	WechatOpenID    sql.NullString      `gorm:"column:wechat_open_id;uniqueIndex:tc_woi;type:VARCHAR(128)"`                                                                                                        // 微信union id
	DingTalkUserID  sql.NullString      `gorm:"column:ding_talk_user_id;uniqueIndex:tc_doi;type:VARCHAR(128)"`
	DingTalkUnionID sql.NullString      `gorm:"column:ding_talk_union_id;uniqueIndex:tc_doi;type:VARCHAR(128)"`
	LastIP          string              `gorm:"column:last_ip;type:VARCHAR(128);NOT NULL"`                   // 最后登录ip
	LastTokenID     string              `gorm:"column:last_token_id;type:VARCHAR(128);default:''"`           // 最后登录的token ID
	RegIP           string              `gorm:"column:reg_ip;type:VARCHAR(128);NOT NULL"`                    // 注册ip
	Sex             int64               `gorm:"column:sex;type:SMALLINT;default:3;NOT NULL"`                 // 用户的性别，值为1时是男性，值为2时是女性，其他值为未知
	City            string              `gorm:"column:city;type:VARCHAR(50);NOT NULL"`                       // 用户所在城市
	Country         string              `gorm:"column:country;type:VARCHAR(50);NOT NULL"`                    // 用户所在国家
	Province        string              `gorm:"column:province;type:VARCHAR(50);NOT NULL"`                   // 用户所在省份
	Language        string              `gorm:"column:language;type:VARCHAR(50);NOT NULL"`                   // 用户的语言，简体中文为zh_CN
	HeadImg         string              `gorm:"column:head_img;type:VARCHAR(256);NOT NULL"`                  // 用户头像
	Role            int64               `gorm:"column:role;type:BIGINT;NOT NULL"`                            // 用户默认角色（默认使用该角色）
	Tags            map[string]string   `gorm:"column:tags;type:json;serializer:json;NOT NULL;default:'{}'"` // 产品标签
	IsAllData       int64               `gorm:"column:is_all_data;type:SMALLINT;default:1;NOT NULL"`         // 是否所有数据权限（1是，2否）
	DeviceCount     int64               `gorm:"column:device_count;default:0"`                               //用户所拥有的设备数量统计
	Roles           []*SysUserRole      `gorm:"foreignKey:UserID;references:UserID"`
	Tenant          *SysTenantInfo      `gorm:"foreignKey:Code;references:TenantCode"`
	Status          int64               `gorm:"column:status;type:BIGINT;NOT NULL;default:1"` //租戶状态: 1启用 2禁用
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:tc_un;uniqueIndex:tc_doi;uniqueIndex:tc_email;uniqueIndex:tc_phone;uniqueIndex:tc_wui;uniqueIndex:tc_woi"`
}

func (m *SysUserInfo) TableName() string {
	return "sys_user_info"
}

// 应用菜单关联表
type SysUserRole struct {
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`      // id编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL;"`   // 租户编码
	UserID     int64               `gorm:"column:user_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"` // 用户ID
	RoleID     int64               `gorm:"column:role_id;uniqueIndex:ri_mi;NOT NULL;type:BIGINT"` // 角色ID
	Role       *SysRoleInfo        `gorm:"foreignKey:ID;references:RoleID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:ri_mi"`
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
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                // 编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;uniqueIndex:tc_un;"` // 租户编码
	UserID     int64               `gorm:"column:user_id;uniqueIndex:tc_un;type:BIGINT;NOT NULL"`           // 用户id
	Code       string              `gorm:"column:code;type:VARCHAR(50);uniqueIndex:tc_un;NOT NULL"`         //配置code
	Params     string              `gorm:"column:params;type:text;NOT NULL"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:tc_un;"`
}

func (m *SysUserProfile) TableName() string {
	return "sys_user_profile"
}
