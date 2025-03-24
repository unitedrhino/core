package dataType

import (
	"context"
	"database/sql/driver"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeptID int64

func (t DeptID) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}

	return
}
func (t *DeptID) Scan(value interface{}) error {
	ret := utils.ToInt64(value)
	p := DeptID(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t DeptID) Value() (driver.Value, error) {
	return int64(t), nil
}
