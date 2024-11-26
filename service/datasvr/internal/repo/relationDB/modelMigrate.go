package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm/clause"
	"sync"
)

var once sync.Once

func Migrate(c conf.Database) (err error) {
	//if c.IsInitTable == false {
	//	return
	//}
	once.Do(func() {
		db := stores.GetCommonConn(context.TODO())
		var needInitColumn bool
		if !db.Migrator().HasTable(&DataStatisticsInfo{}) {
			//需要初始化表
			needInitColumn = true
		}
		err = db.AutoMigrate(
			&DataStatisticsInfo{},
		)
		if err != nil {
			return
		}
		if needInitColumn {
			err = migrateTableColumn()
		}
	})
	return
}
func migrateTableColumn() error {
	db := stores.GetCommonConn(context.TODO()).Clauses(clause.OnConflict{DoNothing: true})
	if err := db.CreateInBatches(&MigrateStatisticsInfo, 100).Error; err != nil {
		return err
	}
	return nil
}

var (
	MigrateStatisticsInfo = []DataStatisticsInfo{
		{IsFilterTenant: 1, IsFilterProject: 1, IsFilterArea: 1, Code: "dmDevicePower", Type: "table", Table: "data_dm_device_power", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{
			"startDate": {Sql: "?>=date", ValNum: 1, Type: "date"}, //开始时间
			"endDate":   {Sql: "?<=date", ValNum: 1, Type: "date"}, //开始时间
		}},
		{IsFilterTenant: 1, IsFilterProject: 1, IsFilterArea: 2, Code: "dmDeviceCount", Type: "table", Table: "dm_device_info", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{
			"areas": {Sql: "area_id in ?", ValNum: 1, Type: "array"},
		}},
		{IsFilterTenant: 2, IsFilterProject: 2, IsFilterArea: 2, Code: "sysOpsWorkOrder", Type: "table", Table: "sys_ops_work_order", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{}},
		{IsFilterTenant: 2, IsFilterProject: 2, IsFilterArea: 2, Code: "sysUserAreaApply", Type: "table", Table: "sys_user_area_apply", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{}},
		{IsFilterTenant: 1, IsFilterProject: 2, IsFilterArea: 2, Code: "dmDeviceCountDistributor", Type: "table", Table: "dm_device_info", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{
			"areas": {Sql: "area_id in ?", ValNum: 1, Type: "array"},
		}},
		{IsFilterTenant: 2, IsFilterProject: 2, IsFilterArea: 2, Code: "dmProductCount", Type: "table", Table: "dm_product_info", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{}},
		{IsFilterTenant: 1, IsFilterProject: 2, IsFilterArea: 2, Code: "saleDistributorCount", Type: "table", Table: "sale_distributor_info", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{}},
		{IsFilterTenant: 1, IsFilterProject: 2, IsFilterArea: 2, Code: "saleDistributorApplyCount", Type: "table", Table: "sale_distributor_apply", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{}},
		{IsFilterTenant: 1, IsFilterProject: 2, IsFilterArea: 2, Code: "saleDistributorWaterCount", Type: "table", Table: "sale_distributor_water", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{}},
		{IsFilterTenant: 1, IsFilterProject: 2, IsFilterArea: 2, Code: "sysUserInfo", Type: "table", Table: "sys_user_info", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{}},
		{IsFilterTenant: 1, IsFilterProject: 2, IsFilterArea: 2, Code: "saleOrderInfoCount", Type: "table", Table: "sale_order_info", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{}},
		{IsFilterTenant: 1, IsFilterProject: 2, IsFilterArea: 2, Code: "sysOperLog", Type: "table", Table: "sys_oper_log", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{}},
		{IsFilterTenant: 1, IsFilterProject: 2, IsFilterArea: 2, Code: "sysLoginLog", Type: "table", Table: "sys_login_log", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{}},
		{IsFilterTenant: 2, IsFilterProject: 2, IsFilterArea: 2, Code: "dmDeviceMsgCount", Type: "table", Table: "dm_device_msg_count", Omits: "created_time,updated_time", OrderBy: "created_time desc", Filter: map[string]FilterKeywords{}},
	}
)
