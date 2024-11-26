package relationDB

import (
	"gitee.com/unitedrhino/share/stores"
)

// 示例
type DataExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

//// 数据统计
//type DataStatisticsInfo struct {
//	ID         int64             `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"` // 编号
//	TenantCode stores.TenantCode `gorm:"column:tenant_code;type:VARCHAR(50);NOT NULL"`     // 租户编码
//	ProjectID  stores.ProjectID  `gorm:"column:project_id;type:bigint;NOT NULL"`           // 所属项目ID(雪花ID)
//	AreaID     stores.AreaID     `gorm:"column:area_id;type:bigint;NOT NULL"`              // 区域ID(雪花ID)
//	Int1       int64             `gorm:"column:int1;type:BIGINT;"`
//	Int2       int64             `gorm:"column:int2;type:BIGINT;"`
//	Int3       int64             `gorm:"column:int3;type:BIGINT;"`
//	Int4       int64             `gorm:"column:int4;type:BIGINT;"`
//	Int5       int64             `gorm:"column:int5;type:BIGINT;"`
//	Int6       int64             `gorm:"column:int6;type:BIGINT;"`
//	Int7       int64             `gorm:"column:int7;type:BIGINT;"`
//	Int8       int64             `gorm:"column:int8;type:BIGINT;"`
//	Int9       int64             `gorm:"column:int9;type:BIGINT;"`
//	Int10      int64             `gorm:"column:int10;type:BIGINT;"`
//	String1    string            `gorm:"column:string;type:VARCHAR(512);default:''"`
//	String2    string            `gorm:"column:string;type:VARCHAR(512);default:''"`
//	String3    string            `gorm:"column:string;type:VARCHAR(512);default:''"`
//	String4    string            `gorm:"column:string;type:VARCHAR(512);default:''"`
//	String5    string            `gorm:"column:string;type:VARCHAR(512);default:''"`
//	String6    string            `gorm:"column:string;type:VARCHAR(512);default:''"`
//	String7    string            `gorm:"column:string;type:VARCHAR(512);default:''"`
//	String8    string            `gorm:"column:string;type:VARCHAR(512);default:''"`
//	String9    string            `gorm:"column:string;type:VARCHAR(512);default:''"`
//	String10   string            `gorm:"column:string;type:VARCHAR(512);default:''"`
//	stores.OnlyTime
//}
//
//func (m *DataStatisticsInfo) TableName() string {
//	return "sys_module_menu"
//}

type DataStatisticsInfo struct {
	ID              int64                     `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`              // 编号
	IsFilterTenant  int64                     `gorm:"column:is_filter_tenant;type:BIGINT;default:1;"`                //是否要过滤租户
	IsFilterProject int64                     `gorm:"column:is_filter_project;type:BIGINT;default:1;"`               //是否要过滤项目1 是 2 否
	IsFilterArea    int64                     `gorm:"column:is_filter_area;type:BIGINT;default:1;"`                  //是否要过滤区域
	FilterRoles     string                    `gorm:"column:filter_roles;type:VARCHAR(120);default:'';"`             //是否要过滤区域
	IsSoftDelete    int64                     `gorm:"column:is_soft_delete;type:BIGINT;default:1;"`                  //是否是软删除,是的话会默认加过滤
	Code            string                    `gorm:"column:code;type:VARCHAR(120);not null;uniqueIndex:key"`        //查询的code
	Type            string                    `gorm:"column:type;type:VARCHAR(120);not null"`                        //查询的类别: sql:sql模板替换查询   table: 直接查表
	Table           string                    `gorm:"column:table;type:VARCHAR(120);default:''"`                     //table类型查询的表名
	Omits           string                    `gorm:"column:omits;type:VARCHAR(120);default:''"`                     //忽略的字段列表,table类型需要
	IsToHump        int64                     `gorm:"column:is_to_hump;type:BIGINT;default:1;"`                      //是否转换为驼峰,入参转换为下划线
	Sql             string                    `gorm:"column:sql;type:VARCHAR(2000);default:''"`                      //sql类型的sql内容
	OrderBy         string                    `gorm:"column:order_by;type:VARCHAR(120);default:'created_time desc'"` //排序
	Filter          map[string]FilterKeywords `gorm:"column:filter;type:json;serializer:json;NOT NULL;default:'{}'"`
	FilterSlotCode  string                    `gorm:"column:filter_slot_code;type:VARCHAR(120);default:''"` //第三方过滤插槽code
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:key"`
}

func (m *DataStatisticsInfo) TableName() string {
	return "data_statistics_info"
}

type FilterKeywords struct {
	Sql    string `json:"sql"`
	ValNum int64  `json:"valNum"` //问号的数量
	Type   string `json:"type"`   //time:时间类型时间戳 date: 日期类型的字符串
}
