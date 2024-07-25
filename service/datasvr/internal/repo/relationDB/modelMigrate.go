package relationDB

import (
	"context"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/stores"
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
		{
			IsFilterTenant:  1,
			IsFilterProject: 1,
			IsFilterArea:    1,
			Code:            "dmDevicePower",
			Type:            "table",
			Table:           "data_dm_device_power",
			Omits:           "created_time,updated_time",
			IsToHump:        1,
			Sql:             "",
			OrderBy:         "",
			Filter: map[string]FilterKeywords{
				"startDate": { //开始时间
					Sql:    "?>=date",
					ValNum: 1,
					Type:   "date",
				},
				"endDate": { //开始时间
					Sql:    "?<=date",
					ValNum: 1,
					Type:   "date",
				},
			},
		},
		{
			IsFilterTenant:  1,
			IsFilterProject: 1,
			IsFilterArea:    1,
			Code:            "dmDeviceCount",
			Type:            "table",
			Table:           "dm_device_info",
			Omits:           "created_time,updated_time",
			IsToHump:        1,
			Sql:             "",
			OrderBy:         "",
			Filter: map[string]FilterKeywords{
				"startDate": { //开始时间
					Sql:    "?>=date",
					ValNum: 1,
					Type:   "date",
				},
				"endDate": { //开始时间
					Sql:    "?<=date",
					ValNum: 1,
					Type:   "date",
				},
			},
		},
	}
)
