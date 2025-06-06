package relationDB

import (
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/stores"
)

// SysProjectInfo 项目信息表,在智能家居中一个项目是一个家庭,一个区域是一个房间
type SysProjectInfo struct {
	TenantCode  dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`                      // 租户编码
	AdminUserID int64               `gorm:"column:admin_user_id;type:BIGINT;NOT NULL"`                         // 超级管理员id,拥有全部权限,默认是创建者
	ProjectID   dataType.ProjectID  `gorm:"column:project_id;type:bigint;primary_key;AUTO_INCREMENT;NOT NULL"` // 项目ID(雪花ID)
	ProjectImg  string              `gorm:"column:project_img;type:varchar(1024);default:''"`
	ProjectName string              `gorm:"column:project_name;type:varchar(100);NOT NULL"` // 项目名称
	//Region      string            `gorm:"column:region;type:varchar(100);NOT NULL"`      // 项目省市区县
	Address      string            `gorm:"column:address;type:varchar(512);"`                // 项目详细地址
	AreaCount    int64             `gorm:"column:area_count;type:bigint;default:0;NOT NULL"` //所属区域的数量统计
	UserCount    int64             `gorm:"column:user_count;type:bigint;default:0;NOT NULL"`
	Position     stores.Point      `gorm:"column:position;"` // 项目地址
	DeviceCount  int64             `gorm:"column:device_count;type:bigint;default:0"`
	Area         float32           `gorm:"column:area;default:0"`
	Ppsm         int64             `gorm:"column:ppsm;type:bigint;default:0"`                           //w.h/m2 每平方米功耗 建筑定额能耗 Power per square meter
	Desc         string            `gorm:"column:desc;type:varchar(100);NOT NULL"`                      // 项目备注
	IsSysCreated int64             `gorm:"column:is_sys_created;type:bigint;default:2;NOT NULL"`        //是否是系统创建的,系统创建的只有管理员可以删除
	Tags         map[string]string `gorm:"column:tags;type:json;serializer:json;NOT NULL;default:'{}'"` // 设备标签
	Areas        []*SysAreaInfo    `gorm:"foreignKey:ProjectID;references:ProjectID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0"`
}

func (m *SysProjectInfo) TableName() string {
	return "sys_project_info"
}

// 区域配置表
type SysProjectProfile struct {
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                // 编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;uniqueIndex:tc_un;"` // 租户编码
	ProjectID  dataType.ProjectID  `gorm:"column:project_id;uniqueIndex:tc_un;type:bigint;NOT NULL"`        // 所属项目ID(雪花ID)
	Code       string              `gorm:"column:code;type:VARCHAR(50);uniqueIndex:tc_un;NOT NULL"`         //配置code
	Params     string              `gorm:"column:params;type:text;NOT NULL"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:tc_un;"`
}

func (m *SysProjectProfile) TableName() string {
	return "sys_project_profile"
}

// 区域信息表
type SysAreaInfo struct {
	TenantCode      dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`                   // 租户编码
	ProjectID       dataType.ProjectID  `gorm:"column:project_id;type:bigint;NOT NULL"`                         // 所属项目ID(雪花ID)
	AreaID          dataType.AreaID     `gorm:"column:area_id;type:bigint;primary_key;AUTO_INCREMENT;NOT NULL"` // 区域ID(雪花ID)
	ParentAreaID    int64               `gorm:"column:parent_area_id;index;type:bigint;default:1"`              // 上级区域ID(雪花ID)
	AreaIDPath      string              `gorm:"column:area_id_path;index;type:varchar(1024);NOT NULL"`          // 1-2-3-的格式记录顶级区域到当前区域的路径
	AreaNamePath    string              `gorm:"column:area_name_path;type:varchar(1024);NOT NULL"`              // 1-2-3-的格式记录顶级区域到当前区域的路径
	AreaName        string              `gorm:"column:area_name;index;type:varchar(100);NOT NULL"`              // 区域名称
	AreaImg         string              `gorm:"column:area_img;type:varchar(1024);NOT NULL"`
	ConfigFile      string              `gorm:"column:config_file;type:varchar(1024);NOT NULL"`
	Position        stores.Point        `gorm:"column:position;"`                                        // 区域定位(默认火星坐标系)
	Desc            string              `gorm:"column:desc;type:varchar(100);NOT NULL"`                  // 区域备注
	LowerLevelCount int64               `gorm:"column:lower_level_count;type:bigint;default:0;NOT NULL"` //下级区域的数量统计
	DeviceCount     int64               `gorm:"column:device_count;type:bigint;default:0;"`
	GroupCount      int64               `gorm:"column:group_count;type:bigint;default:0;"`
	IsLeaf          int64               `gorm:"column:is_leaf;type:bigint;default:1;NOT NULL"`               //是否是叶子节点
	Tags            map[string]string   `gorm:"column:tags;type:json;serializer:json;NOT NULL;default:'{}'"` // 设备标签
	UseBy           string              `gorm:"column:use_by;type:varchar(100);default:''"`                  //用途
	ChildrenAreaIDs []int64             `gorm:"column:children_area_ids;type:json;serializer:json"`          //所有的子区域的id列表
	IsSysCreated    int64               `gorm:"column:is_sys_created;type:bigint;default:2;NOT NULL"`        //是否是系统创建的,系统创建的只有管理员可以删除
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;index"`
	Children    []*SysAreaInfo     `gorm:"foreignKey:ParentAreaID;references:AreaID"`
	Parent      *SysAreaInfo       `gorm:"foreignKey:AreaID;references:ParentAreaID"`
}

func (m *SysAreaInfo) TableName() string {
	return "sys_area_info"
}

// 区域配置表
type SysAreaProfile struct {
	ID         int64               `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                // 编号
	TenantCode dataType.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL;uniqueIndex:tc_un;"` // 租户编码
	ProjectID  dataType.ProjectID  `gorm:"column:project_id;uniqueIndex:tc_un;type:bigint;NOT NULL"`        // 所属项目ID(雪花ID)
	AreaID     dataType.AreaID     `gorm:"column:area_id;uniqueIndex:tc_un;type:bigint;NOT NULL"`           // 区域ID(雪花ID)
	Code       string              `gorm:"column:code;type:VARCHAR(50);uniqueIndex:tc_un;NOT NULL"`         //配置code
	Params     string              `gorm:"column:params;type:text;NOT NULL"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:tc_un;"`
}

func (m *SysAreaProfile) TableName() string {
	return "sys_area_profile"
}
