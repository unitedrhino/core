package dataType

import (
	"context"
	"database/sql/driver"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeptIDPath string

func (t DeptIDPath) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	expr = clause.Expr{SQL: "?", Vars: []interface{}{string(t)}}
	return
}

func (t *DeptIDPath) Scan(value interface{}) error {
	ret := utils.ToString(value)
	p := DeptIDPath(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.

func (t DeptIDPath) Value() (driver.Value, error) {
	return string(t), nil
}
