package dataType

import (
	"context"
	"database/sql/driver"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AreaIDPath string

func (t AreaIDPath) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口

	stmt := db.Statement
	uc := ctxs.GetUserCtxOrNil(ctx)
	expr = clause.Expr{SQL: "?", Vars: []interface{}{string(t)}}
	authType, areas := ctxs.GetAreaIDPaths(uc.ProjectID, uc.ProjectAuth)
	if t != "" && !(uc.IsAdmin || uc.AllArea || authType <= def.AuthReadWrite || utils.SliceIn(string(t), areas...)) { //如果没有权限
		stmt.Error = errors.Permissions.WithMsg("区域权限不足")
	}
	return
}

func (t *AreaIDPath) Scan(value interface{}) error {
	if v, ok := value.([]byte); ok {
		if len(v) == 1 && v[0] == 0 {
			*t = ""
			return nil
		}
	}
	ret := utils.ToString(value)
	p := AreaIDPath(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.

func (t AreaIDPath) Value() (driver.Value, error) {
	return string(t), nil
}

//	func (t AreaIDPath) QueryClauses(f *schema.Field) []clause.Interface {
//		return []clause.Interface{AreaIDPathClause{Field: f, T: t, Opt: Select}}
//	}
//
//	func (t AreaIDPath) UpdateClauses(f *schema.Field) []clause.Interface {
//		return []clause.Interface{AreaIDPathClause{Field: f, T: t, Opt: Update}}
//	}
//
//	func (t AreaIDPath) CreateClauses(f *schema.Field) []clause.Interface {
//		return []clause.Interface{AreaIDPathClause{Field: f, T: t, Opt: Create}}
//	}
//
//	func (t AreaIDPath) DeleteClauses(f *schema.Field) []clause.Interface {
//		return []clause.Interface{AreaIDPathClause{Field: f, T: t, Opt: Delete}}
//	}
//
//	type AreaIDPathClause struct {
//		clauseInterface
//		Field *schema.Field
//		T     AreaIDPath
//		Opt   Opt
//	}
//
// func (sd AreaIDPathClause) GenAuthKey() string { //查询的时候会调用此接口
//
//		return fmt.Sprintf(AuthModify, "AreaIDPath")
//	}
//
// func (sd AreaIDPathClause) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口
//
//		uc := ctxs.GetUserCtxOrNil(stmt.Context)
//		if uc == nil {
//			return
//		}
//		authType, areas := ctxs.GetAreaIDPaths(uc.ProjectID, uc.ProjectAuth)
//		if uc.IsAdmin || uc.AllArea || authType <= def.AuthReadWrite {
//			return
//		}
//		switch sd.Opt {
//		case Create:
//		case Update, Delete, Select:
//			if _, ok := stmt.Clauses[sd.GenAuthKey()]; !ok {
//				if c, ok := stmt.Clauses["WHERE"]; ok {
//					if where, ok := c.Expression.(clause.Where); ok && len(where.Exprs) > 1 {
//						for _, expr := range where.Exprs {
//							if orCond, ok := expr.(clause.OrConditions); ok && len(orCond.Exprs) == 1 {
//								where.Exprs = []clause.Expression{clause.And(where.Exprs...)}
//								c.Expression = where
//								stmt.Clauses["WHERE"] = c
//								break
//							}
//						}
//					}
//				}
//				if len(areas) == 0 { //如果没有权限
//					//stmt.Error = errors.Permissions.WithMsg("区域权限不足")
//					stmt.AddClause(clause.Where{Exprs: []clause.Expression{
//						clause.IN{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Values: nil},
//					}})
//					stmt.Clauses[sd.GenAuthKey()] = clause.Clause{}
//					return
//				}
//				for _, v := range areas {
//					stmt.AddClause(clause.Where{Exprs: []clause.Expression{
//						clause.Like{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Value: v + "%"},
//					}})
//				}
//
//				stmt.Clauses[sd.GenAuthKey()] = clause.Clause{}
//			}
//		}
//	}
func GetAreaAuthIDPaths(ctx context.Context) ([]string, error) {
	uc := ctxs.GetUserCtxOrNil(ctx)
	if uc == nil {
		return nil, nil
	}
	authType, areas := ctxs.GetAreaIDPaths(uc.ProjectID, uc.ProjectAuth)
	if uc.IsAdmin || uc.AllArea || authType <= def.AuthReadWrite {
		return nil, nil
	}
	if len(areas) == 0 { //如果没有权限
		return nil, errors.Permissions.WithMsg("区域权限不足")
	}
	var values = []string{}
	for _, v := range areas {
		values = append(values, v)
	}
	return values, nil
}
