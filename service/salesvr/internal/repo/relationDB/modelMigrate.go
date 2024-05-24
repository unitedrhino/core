package relationDB

import (
	"context"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/stores"
)

func Migrate(c conf.Database) error {
	if c.IsInitTable == false {
		return nil
	}
	db := stores.GetCommonConn(context.TODO())
	var needInitColumn bool
	//if !db.Migrator().HasTable(&SysUserInfo{}) {
	//	//需要初始化表
	//	needInitColumn = true
	//}
	err := db.AutoMigrate()
	if err != nil {
		return err
	}

	if needInitColumn {
		return migrateTableColumn()
	}
	return err
}
func migrateTableColumn() error {
	//db := stores.GetCommonConn(context.TODO()).Clauses(clause.OnConflict{DoNothing: true})

	return nil
}
