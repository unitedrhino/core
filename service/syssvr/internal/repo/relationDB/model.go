package relationDB

import (
	"gitee.com/unitedrhino/core/service/syssvr/internal/domain/dept"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
)

func normalizeJSONMapStringString(in map[string]string) map[string]string {
	if in == nil {
		return map[string]string{}
	}
	return in
}

func normalizeJSONString(in string) string {
	if in == "" {
		return "{}"
	}
	return in
}

func normalizeStringSlice(in []string) []string {
	if in == nil {
		return []string{}
	}
	return in
}

// 示例
type SysExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

// 全局配置
type SysConfig struct {
	Sms *SysConfigSms `gorm:"embedded;embeddedPrefix:sms_"` //短信配置,全租户共用
}

type Attachment struct {
	ID       int64  `json:"id"`
	FilePath string `json:"filePath"`
	UseBy    string `json:"useBy"`
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
	ID         int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                                 // id编号
	Name       string `gorm:"column:name;uniqueIndex:idx_sys_dict_info_name;type:VARCHAR(100);default:'';NOT NULL;comment:字典名"` // 字典名（中）
	Code       string `gorm:"column:code;uniqueIndex:idx_sys_dict_info_code;type:VARCHAR(50);default:'';NOT NULL"`              //编码
	Group      string `gorm:"column:group;type:VARCHAR(50);default:'';NOT NULL"`                                                //字典分组
	Desc       string `gorm:"column:desc;comment:描述"`                                                                           // 描述
	Body       string `gorm:"column:body;type:VARCHAR(1024)"`                                                                   // 自定义数据
	StructType int64  `gorm:"column:struct_type;type:BIGINT;default:1"`                                                         //结构类型(不可修改) 1:列表(默认) 2:树型
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_dict_info_code;uniqueIndex:idx_sys_dict_info_name"`
	Details     []*SysDictDetail   `gorm:"foreignKey:DictCode;references:Code"`
}

func (SysDictInfo) TableName() string {
	return "sys_dict_info"
}

type SysDictDetail struct {
	ID       int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                                     // id编号
	DictCode string `gorm:"column:dict_code;uniqueIndex:idx_sys_dict_detail_value;type:VARCHAR(50);default:'';NOT NULL"`          // 关联标记
	Label    string `gorm:"column:label;type:VARCHAR(100);default:'';NOT NULL;comment:展示值"`                                       // 展示值
	Value    string `gorm:"column:value;uniqueIndex:idx_sys_dict_detail_value;type:VARCHAR(100);default:'';NOT NULL;comment:字典值"` // 字典值
	Status   int64  `gorm:"column:status;type:SMALLINT;default:1"`                                                                // 状态  1:启用,2:禁用
	Sort     int64  `gorm:"column:sort;comment:排序标记;default:1"`                                                                   // 排序标记
	Desc     string `gorm:"column:desc;comment:描述"`                                                                               // 描述
	Body     string `gorm:"column:body;type:VARCHAR(1024)"`                                                                       // 自定义数据
	IDPath   string `gorm:"column:id_path;type:varchar(100);NOT NULL"`                                                            // 1-2-3-的格式记录顶级区域到当前id的路径
	ParentID int64  `gorm:"column:parent_id;uniqueIndex:idx_sys_dict_detail_value;type:BIGINT"`                                   // id编号
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_dict_detail_value"`
	Children    []*SysDictDetail   `gorm:"foreignKey:parent_id;references:id"`
	Parent      *SysDictDetail     `gorm:"foreignKey:ID;references:ParentID"`
}

func (SysDictDetail) TableName() string {
	return "sys_dict_detail"
}

type SysDeptInfo struct {
	ID             dataType.DeptID     `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                     // id编号
	TenantCode     dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`                                         // 租户编码
	ParentID       int64               `gorm:"column:parent_id;uniqueIndex:idx_sys_dept_info_name;type:BIGINT"`                      // id编号
	Name           string              `gorm:"column:name;type:VARCHAR(256);uniqueIndex:idx_sys_dept_info_name;default:'';NOT NULL"` // 部门名称
	AdminUserID    int64               `gorm:"column:admin_user_id;comment:管理员账号;NOT NULL"`
	Status         int64               `gorm:"column:status;type:SMALLINT;default:1"` // 状态  1:启用,2:禁用
	Sort           int64               `gorm:"column:sort;comment:排序标记"`              // 排序标记
	Desc           string              `gorm:"column:desc;comment:描述"`                // 描述
	UserCount      int64               `gorm:"column:user_count;comment:用户统计,包含下级部门的人数"`
	DeviceCount    int64               `gorm:"column:device_count;default:0;comment:部门自己的设备总数"`
	AllDeviceCount int64               `gorm:"column:all_device_count;default:0;comment:部门及其下级的设备总数"`
	ChildrenCount  int64               `gorm:"column:children_count;default:0;comment:部门下级数量"`
	IDPath         dataType.DeptIDPath `gorm:"column:id_path;type:varchar(100);NOT NULL"` // 1-2-3-的格式记录顶级区域到当前id的路径
	DingTalkID     int64               `gorm:"column:ding_talk_id;default:0;"`            //钉钉的部门ID
	Tags           map[string]string   `gorm:"column:tags;type:json;serializer:json"`     //部门标签
	AdminUser      *SysUserInfo        `gorm:"foreignKey:user_id;references:AdminUserID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_dept_info_name"`
	Children    []*SysDeptInfo     `gorm:"foreignKey:parent_id;references:id"`
	Parent      *SysDeptInfo       `gorm:"foreignKey:ID;references:ParentID"`
}

func (SysDeptInfo) TableName() string {
	return "sys_dept_info"
}

type Company struct {
	LegalName   string `gorm:"column:legal_name;type:varchar(255);comment:法人名字;default:'';"`
	Code        string `gorm:"column:code;type:varchar(255);comment:统一社会信用代码;default:'';"`
	Address     string `gorm:"column:address;type:varchar(255);comment:地址;default:''"`
	BankName    string `gorm:"column:bank_name;type:varchar(255);comment:开户银行;default:''"`
	BankAccount string `gorm:"column:bank_account;type:varchar(255);comment:银行账户;default:''"`
}

type SysDeptUser struct {
	ID             int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                        // id编号
	TenantCode     dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;"`                           // 租户编码
	UserID         int64               `gorm:"column:user_id;uniqueIndex:idx_sys_dept_user_ri_mi;NOT NULL;type:BIGINT"` // 用户ID
	DeptID         int64               `gorm:"column:dept_id;uniqueIndex:idx_sys_dept_user_ri_mi;NOT NULL;type:BIGINT"` // 角色ID
	DeptIDPath     string              `gorm:"column:dept_id_path;type:varchar(100);NOT NULL"`                          // 1-2-3-的格式记录顶级区域到当前id的路径
	AuthType       def.AuthType        `gorm:"column:auth_type;type:bigint;NOT NULL"`                                   // 授权类型 1 管理员(可以调整本区域及旗下区域的设备区域规划)  2 读写授权(可以对该区域及旗下区域的设备进行管理) 3 只读授权()
	IsAuthChildren int64               `gorm:"column:is_auth_children;type:bigint;default:1;NOT NULL"`                  //是否同时授权子节点
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_dept_user_ri_mi"`
	Dept        *SysDeptInfo       `gorm:"foreignKey:ID;references:DeptID"`
	User        *SysUserInfo       `gorm:"foreignKey:UserID;references:UserID"`
}

func (m *SysDeptUser) TableName() string {
	return "sys_dept_user"
}

// 租户下的应用列表
type SysDeptSyncJob struct {
	ID          int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                     // id编号
	TenantCode  dataType.TenantCode `gorm:"column:tenant_code;uniqueIndex:idx_sys_dept_sync_job_tc_ac;type:VARCHAR(50);NOT NULL"` // 租户编码
	Direction   dept.SyncDirection  `gorm:"column:direction;default:1"`                                                           // 同步的方向,1上游同步到联犀(默认),2联犀同步到下游
	ThirdType   def.AppSubType      `gorm:"column:third_type;type:varchar(20)"`                                                   //同步的类型
	ThirdConfig *SysTenantThird     `gorm:"embedded;embeddedPrefix:third_config"`                                                 //第三方配置
	FieldMap    map[string]string   `gorm:"column:field_map;type:json;serializer:json"`                                           //用户字段映射,左边是联犀的字段,右边是第三方的,不填写就是全量映射
	SyncDeptIDs []int64             `gorm:"column:sync_dept_ids;type:json;serializer:json"`                                       //同步的第三方部门id列表,不填为同步全部
	IsAddSync   int64               `gorm:"column:is_add_sync;default:1"`                                                         //新增人员自动同步,默认为1
	SyncMode    dept.SyncMode       `gorm:"column:sync_mode;default:1"`                                                           //同步模式: 1:手动(默认) 2: 定时同步(半小时) 3: 实时同步

	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_dept_sync_job_tc_ac"`
}

func (m *SysDeptSyncJob) TableName() string {
	return "sys_dept_sync_job"
}

type SysSlotInfo struct {
	ID       int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                    // id编号
	Code     string            `gorm:"column:code;uniqueIndex:idx_sys_slot_info_code_slot;type:VARCHAR(100);NOT NULL"`      // 鉴权的编码
	SubCode  string            `gorm:"column:sub_code;uniqueIndex:idx_sys_slot_info_code_slot;type:VARCHAR(100);NOT NULL"`  // 鉴权的编码
	SlotCode string            `gorm:"column:slot_code;uniqueIndex:idx_sys_slot_info_code_slot;type:VARCHAR(100);NOT NULL"` //slot的编码
	Method   string            `gorm:"column:method;type:VARCHAR(50);default:'POST'"`                                       // 请求方式 GET  POST
	Uri      string            `gorm:"column:uri;type:VARCHAR(100);NOT NULL"`                                               // 参考: /api/v1/system/user/self/captcha?fwefwf=gwgweg&wefaef=gwegwe
	Hosts    []string          `gorm:"column:hosts;type:json;serializer:json;NOT NULL"`                                     //访问的地址 host or host:port
	Body     string            `gorm:"column:body;type:VARCHAR(100);default:''"`                                            // body 参数模板
	Handler  map[string]string `gorm:"column:handler;type:json;serializer:json;NOT NULL"`                                   //http头
	AuthType string            `gorm:"column:auth_type;type:VARCHAR(100);NOT NULL"`                                         //鉴权类型 core
	Desc     string            `gorm:"column:desc;type:VARCHAR(500);"`                                                      // 备注
	stores.SoftTime
}

func (m *SysSlotInfo) TableName() string {
	return "sys_slot_info"
}

func (m *SysSlotInfo) BeforeSave(tx *gorm.DB) error {
	if m == nil {
		return nil
	}
	m.Hosts = normalizeStringSlice(m.Hosts)
	m.Handler = normalizeJSONMapStringString(m.Handler)
	return nil
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
