package relationDB

import (
	"gitee.com/unitedrhino/core/service/syssvr/internal/domain/module"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
)

// 应用信息
type SysAppInfo struct {
	ID            int64          `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                         // id编号
	Code          string         `gorm:"column:code;uniqueIndex:idx_sys_app_info_code;type:VARCHAR(100);NOT NULL"` // 应用编码
	Name          string         `gorm:"column:name;uniqueIndex:idx_sys_app_info_name;type:VARCHAR(100);NOT NULL"` //应用名称
	Type          def.AppType    `gorm:"column:type;type:VARCHAR(100);default:web;NOT NULL"`                       //应用类型 web:web页面  app:应用  mini:小程序
	SubType       def.AppSubType `gorm:"column:sub_type;type:VARCHAR(100);default:wx;NOT NULL"`                    // 类型  wx:微信小程序  ding:钉钉小程序
	Desc          string         `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`                                   //应用描述
	IsCommon      def.Bool       `gorm:"column:is_common;type:TINYINT;default:2"`                                  //是否是公共应用,公共应用所有租户共用,只有default租户能修改,其他租户只能读
	IsSetDingMini def.Bool       `gorm:"column:is_set_dingMini;type:TINYINT;default:2"`                            //是否设置了钉钉小程序
	IsSetWxMini   def.Bool       `gorm:"column:is_set_wxMini;type:TINYINT;default:2"`                              //是否设置了微信小程序
	IsSetWxOpen   def.Bool       `gorm:"column:is_set_wxOpen;type:TINYINT;default:2"`                              //是否设置了微信开放平台
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_app_info_code;uniqueIndex:idx_sys_app_info_name"`
}

func (m *SysAppInfo) TableName() string {
	return "sys_app_info"
}

// 应用默认绑定的模块
type SysAppModule struct {
	ID         int64  `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                  // id编号
	AppCode    string `gorm:"column:app_code;uniqueIndex:idx_sys_app_module_tc_ac;type:VARCHAR(50);NOT NULL"`    // 应用编码 这里只关联主应用,主应用授权,子应用也授权了
	ModuleCode string `gorm:"column:module_code;uniqueIndex:idx_sys_app_module_tc_ac;type:VARCHAR(50);NOT NULL"` // 模块编码
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_app_module_tc_ac"`
	Module      *SysModuleInfo     `gorm:"foreignKey:Code;references:ModuleCode"`
	App         *SysAppInfo        `gorm:"foreignKey:Code;references:AppCode"`
}

func (m *SysAppModule) TableName() string {
	return "sys_app_module"
}

// 模块管理表 模块是菜单和接口的集合体
type SysModuleInfo struct {
	ID         int64            `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                          // 编号
	Code       string           `gorm:"column:code;uniqueIndex:idx_sys_app_module_code;NOT NULL;type:VARCHAR(50)"` // 编码
	Type       int64            `gorm:"column:type;type:BIGINT;default:1;NOT NULL"`                                // 类型   1:web页面  2:应用  3:小程序
	SubType    int64            `gorm:"column:sub_type;type:BIGINT;default:1;NOT NULL"`                            // 类型   1：微应用   2：iframe内嵌 3: 原生菜单
	Order      int64            `gorm:"column:order;type:BIGINT;default:1;NOT NULL"`                               // 左侧table排序序号
	Name       string           `gorm:"column:name;type:VARCHAR(50);NOT NULL"`                                     // 菜单名称
	Path       string           `gorm:"column:path;type:VARCHAR(64);NOT NULL"`                                     // 系统的path
	Url        string           `gorm:"column:url;type:VARCHAR(200);NOT NULL"`                                     // 页面
	Icon       string           `gorm:"column:icon;type:VARCHAR(64);NOT NULL"`                                     // 图标
	Body       string           `gorm:"column:body;type:VARCHAR(1024)"`                                            // 菜单自定义数据
	HideInMenu int64            `gorm:"column:hide_in_menu;type:BIGINT;default:2;NOT NULL"`                        // 是否隐藏菜单 1-是 2-否
	Desc       string           `gorm:"column:desc;type:VARCHAR(100);NOT NULL"`                                    // 备注
	Tag        int64            `gorm:"column:tag;type:BIGINT;default:1;NOT NULL"`                                 //标签: 1:通用 2:选配
	Purpose    module.Purpose   `gorm:"column:purpose;type:BIGINT;default:1;NOT NULL"`                             // platform(那么只有default租户可以看,然后平台模块http头里不用传租户号) normal project(需要选择项目,默认选择第一个)
	HomeMenuID int64            `gorm:"column:home_menu_id;type:BIGINT;default:1;"`                                // 模块首页全屏菜单,点击应用左上角的logo会跳出这个页面及首次进入该模块的时候,如果为1 则是没有
	Home       *SysModuleMenu   `gorm:"foreignKey:ID;references:HomeMenuID"`                                       // 模块首页全屏菜单,点击应用左上角的logo会跳出这个页面及首次进入该模块的时候,如果为1 则是没有
	Menus      []*SysModuleMenu `gorm:"foreignKey:ModuleCode;references:Code"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_sys_app_module_code"`
}

func (m *SysModuleInfo) TableName() string {
	return "sys_module_info"
}

// 菜单管理表
type SysModuleMenu struct {
	ID          int64            `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`                                         // 编号
	ModuleCode  string           `gorm:"column:module_code;uniqueIndex:idx_sys_app_module_menu_path;type:VARCHAR(50);NOT NULL"`    // 模块编码
	ParentID    int64            `gorm:"column:parent_id;uniqueIndex:idx_sys_app_module_menu_path;type:BIGINT;default:1;NOT NULL"` // 父菜单ID，一级菜单为1
	Type        int64            `gorm:"column:type;type:BIGINT;default:1;NOT NULL"`                                               // 类型   1：菜单或者页面   2：iframe嵌入   3：外链跳转
	Order       int64            `gorm:"column:order;type:BIGINT;default:1;NOT NULL"`                                              // 左侧table排序序号
	Name        string           `gorm:"column:name;type:VARCHAR(50);NOT NULL"`                                                    // 菜单名称
	Path        string           `gorm:"column:path;uniqueIndex:idx_sys_app_module_menu_path;type:VARCHAR(64);"`                   // 系统的path
	Component   string           `gorm:"column:component;type:VARCHAR(1024);"`                                                     // 页面
	Icon        string           `gorm:"column:icon;type:VARCHAR(64);"`                                                            // 图标
	Redirect    string           `gorm:"column:redirect;type:VARCHAR(64)"`                                                         // 路由重定向
	Body        string           `gorm:"column:body;type:VARCHAR(1024)"`                                                           // 菜单自定义数据
	HideInMenu  int64            `gorm:"column:hide_in_menu;type:BIGINT;default:2"`                                                // 是否隐藏菜单 1-是 2-否
	IsCommon    int64            `gorm:"column:is_common;type:BIGINT;default:2;"`                                                  // 是否常用菜单 1-是 2-否
	IsAllTenant int64            `gorm:"column:is_all_tenant;type:BIGINT;default:1;"`                                              // 菜单是否提供给所有租户 1-是 2-否
	Children    []*SysModuleMenu `gorm:"foreignKey:ID;references:ParentID"`
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;uniqueIndex:idx_sys_app_module_menu_path;default:0;index"`
}

func (m *SysModuleMenu) TableName() string {
	return "sys_module_menu"
}
