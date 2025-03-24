package dataType

import (
	"context"
	"database/sql/driver"
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type AreaID int64

func (t AreaID) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	stmt := db.Statement
	uc := ctxs.GetUserCtxOrNil(ctx)
	expr = clause.Expr{SQL: "?", Vars: []interface{}{int64(t)}}

	authType, areas := ctxs.GetAreaIDs(uc.ProjectID, uc.ProjectAuth)
	if t != def.NotClassified && !(uc.IsAdmin || uc.AllArea || authType <= def.AuthReadWrite || utils.SliceIn(int64(t), areas...)) { //如果没有权限
		stmt.Error = errors.Permissions.WithMsg("区域权限不足")
	}
	return
}
func (t *AreaID) Scan(value interface{}) error {
	ret := utils.ToInt64(value)
	p := AreaID(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t AreaID) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t AreaID) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AreaClause{Field: f, T: t, Opt: stores.Select}}
}
func (t AreaID) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AreaClause{Field: f, T: t, Opt: stores.Update}}
}

func (t AreaID) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AreaClause{Field: f, T: t, Opt: stores.Create}}
}

func (t AreaID) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{AreaClause{Field: f, T: t, Opt: stores.Delete}}
}

type AreaClause struct {
	stores.ClauseInterface
	Field *schema.Field
	T     AreaID
	Opt   stores.Opt
}

func (sd AreaClause) GenAuthKey() string { //查询的时候会调用此接口
	return fmt.Sprintf(stores.AuthModify, "areaID")
}

func (sd AreaClause) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口
	uc := ctxs.GetUserCtxOrNil(stmt.Context)
	if uc == nil {
		return
	}
	authType, areaIDs := ctxs.GetAreaIDs(uc.ProjectID, uc.ProjectAuth)
	_, areaIDPaths := ctxs.GetAreaIDPaths(uc.ProjectID, uc.ProjectAuth)
	if uc.IsAdmin || uc.AllArea || authType <= def.AuthReadWrite {
		return
	}
	switch sd.Opt {
	case stores.Create:
	case stores.Update, stores.Delete, stores.Select:
		if _, ok := stmt.Clauses[sd.GenAuthKey()]; !ok {
			if c, ok := stmt.Clauses["WHERE"]; ok {
				if where, ok := c.Expression.(clause.Where); ok && len(where.Exprs) > 1 {
					for _, expr := range where.Exprs {
						if orCond, ok := expr.(clause.OrConditions); ok && len(orCond.Exprs) == 1 {
							where.Exprs = []clause.Expression{clause.And(where.Exprs...)}
							c.Expression = where
							stmt.Clauses["WHERE"] = c
							break
						}
					}
				}
			}
			var expression []clause.Expression
			func() { //区域授权
				areaIDPathF := stmt.Schema.FieldsByName["AreaIDPath"]
				if areaIDPathF != nil && len(areaIDPaths) > 0 {
					for _, v := range areaIDPaths {
						expression = append(expression, clause.Like{Column: clause.Column{Table: clause.CurrentTable, Name: areaIDPathF.DBName}, Value: v + "%"})
					}
				}
				if len(areaIDs) == 0 { //如果没有权限
					return
				}
				var values = []any{def.NotClassified}
				for _, v := range areaIDs {
					values = append(values, v)
				}
				expression = append(expression, clause.IN{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Values: values})
			}()
			func() { //部门授权
				return
				deptIDPathF := stmt.Schema.FieldsByName["DeptIDPath"]
				if deptIDPathF != nil {
					for _, v := range areaIDPaths {
						expression = append(expression, clause.Like{Column: clause.Column{Table: clause.CurrentTable, Name: deptIDPathF.DBName}, Value: v + "%"})
					}
				}
				deptIDF := stmt.Schema.FieldsByName["DeptID"]
				if deptIDF != nil {
					depts := utils.SetToSlice(uc.Dept)
					if len(depts) == 0 {
						return
					}
					var values = []any{}
					for _, v := range depts {
						values = append(values, v)
					}
					expression = append(expression, clause.IN{Column: clause.Column{Table: clause.CurrentTable, Name: deptIDF.DBName}, Values: values})
				}
			}()
			if len(expression) == 0 {
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{
					clause.IN{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Values: nil},
				}})
			} else {
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{
					clause.OrConditions{Exprs: expression},
				}})
			}
			stmt.Clauses[sd.GenAuthKey()] = clause.Clause{}
		}
	}
}
func GenAreaAuthScope(ctx context.Context, db *gorm.DB) *gorm.DB {
	uc := ctxs.GetUserCtxOrNil(ctx)
	if uc == nil {
		return db
	}
	authType, areas := ctxs.GetAreaIDs(uc.ProjectID, uc.ProjectAuth)
	if uc.IsAdmin || uc.AllArea || authType <= def.AuthReadWrite {
		return db
	}
	if len(areas) == 0 { //如果没有权限
		db.AddError(errors.Permissions.WithMsg("区域权限不足"))
		return db
	}
	var values = []any{def.NotClassified}
	for _, v := range areas {
		values = append(values, v)
	}
	db = db.Where("area_id in ?", values)
	return db
}

func GetAreaAuthIDs(ctx context.Context) ([]int64, error) {
	uc := ctxs.GetUserCtxOrNil(ctx)
	if uc == nil {
		return nil, nil
	}
	authType, areas := ctxs.GetAreaIDs(uc.ProjectID, uc.ProjectAuth)
	if uc.IsAdmin || uc.AllArea || authType <= def.AuthReadWrite {
		return nil, nil
	}
	if len(areas) == 0 { //如果没有权限
		return nil, errors.Permissions.WithMsg("区域权限不足")
	}
	var values = []int64{def.NotClassified}
	for _, v := range areas {
		values = append(values, v)
	}
	return values, nil
}
