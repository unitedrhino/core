package dataType

import (
	"context"
	"database/sql/driver"
	"fmt"
	"reflect"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type TenantCodeWithCommonR string //非root可以看,不可写

func (t TenantCodeWithCommonR) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) { //更新的时候会调用此接口
	stmt := db.Statement
	uc := ctxs.GetUserCtx(ctx)
	if uc == nil { //系统初始化的时候会掉用这里
		expr = clause.Expr{SQL: "?", Vars: []interface{}{string(t)}}
		return
	}
	if uc.TenantCode == "" {
		stmt.Error = errors.Parameter.AddDetail("tenantCode is empty")
		return
	}
	if t != "" {
		if uc.TenantCode == string(t) || uc.IsRoot() {
			expr = clause.Expr{SQL: "?", Vars: []interface{}{string(t)}}
			return
		}
		stmt.Error = errors.Parameter.AddDetailf("tenantCode not eq uc:%v t:%v", uc.TenantCode, string(t))
		return
	}
	expr = clause.Expr{SQL: "?", Vars: []interface{}{uc.TenantCode}}
	return
}
func (t *TenantCodeWithCommonR) Scan(value interface{}) error {
	ret := cast.ToString(value)
	p := TenantCodeWithCommonR(ret)
	*t = p
	return nil
}

// Value implements the driver Valuer interface.
func (t TenantCodeWithCommonR) Value() (driver.Value, error) {
	return string(t), nil
}

func (t TenantCodeWithCommonR) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{TenantCodeWitCommonClause{Field: f, T: t, Opt: stores.Select}}
}

func (t TenantCodeWithCommonR) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{TenantCodeWitCommonClause{Field: f, T: t, Opt: stores.Update}}
}

func (t TenantCodeWithCommonR) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{TenantCodeWitCommonClause{Field: f, T: t, Opt: stores.Create}}
}

func (t TenantCodeWithCommonR) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{TenantCodeWitCommonClause{Field: f, T: t, Opt: stores.Delete}}
}

func (t TenantCodeWithCommonR) GetAuthIDs(f *schema.Field) stores.GetValues {
	return func(stmt *gorm.Statement) (authIDs []any, isRoot bool, allData bool, err error) {
		uc := ctxs.GetUserCtx(stmt.Context)
		if uc == nil {
			return nil, false, false, nil
		}
		if uc.TenantCode == def.TenantCodeDefault { //只有core租户的可以修改其他租户的租户号
			isRoot = true
		}
		return []any{TenantCodeWithCommonR(uc.TenantCode)}, isRoot, uc.AllTenant, nil
	}
}

type TenantCodeWitCommonClause struct {
	stores.ClauseInterface
	Field *schema.Field
	T     TenantCodeWithCommonR
	Opt   stores.Opt
}

func (sd TenantCodeWitCommonClause) GenAuthKey() string { //查询的时候会调用此接口
	return fmt.Sprintf(stores.AuthModify, "tenantCode")
}

func (sd TenantCodeWitCommonClause) ModifyStatement(stmt *gorm.Statement) { //查询的时候会调用此接口

	uc := ctxs.GetUserCtxNoNil(stmt.Context)

	switch sd.Opt {
	case stores.Create:
		destV := reflect.ValueOf(stmt.Dest)
		if destV.Kind() == reflect.Array || destV.Kind() == reflect.Slice {
			for i := 0; i < destV.Len(); i++ {
				dest := destV.Index(i)
				if dest.Kind() == reflect.Pointer || dest.Kind() == reflect.Interface {
					dest = dest.Elem()
				}
				field := dest.FieldByName(sd.Field.Name)
				if field.IsZero() {
					var v TenantCodeWithCommonR
					v = TenantCodeWithCommonR(uc.TenantCode)
					field.Set(reflect.ValueOf(v))
					continue
				}
				vv := field.Interface().(TenantCodeWithCommonR)
				if string(vv) == uc.TenantCode {
					continue
				}
				if vv == def.TenantCodeDefault || vv == def.TenantCodeCommon {
					if uc.IsRoot() {
						continue
					}
				}
				stmt.Error = errors.Parameter.AddDetail("tenantCode not eq uc")
				return
			}
			return
		}
		field := destV.Elem().FieldByName(sd.Field.Name)
		if field.IsZero() {
			var v TenantCodeWithCommonR
			v = TenantCodeWithCommonR(uc.TenantCode)
			field.Set(reflect.ValueOf(v))
			return
		}
		vv := field.Interface().(TenantCodeWithCommonR)
		if string(vv) == uc.TenantCode {
			return
		}
		if vv == def.TenantCodeDefault || vv == def.TenantCodeCommon {
			if uc.IsRoot() {
				return
			}
		}
		stmt.Error = errors.Parameter.AddDetail("tenantCode not eq uc")

	case stores.Select:
		if uc.IsRoot() && uc.AllTenant { //只有超管能修改其他租户
			return
		}
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
			values := []any{uc.TenantCode, def.TenantCodeDefault}
			stmt.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.IN{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Values: values},
			}})
			stmt.Clauses[sd.GenAuthKey()] = clause.Clause{}
		}
	case stores.Update, stores.Delete:
		if uc.IsRoot() && uc.AllTenant { //只有超管能修改其他租户
			return
		}
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
			values := []any{uc.TenantCode}
			if uc.IsRoot() {
				values = append(values, def.TenantCodeDefault)
			}
			stmt.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.IN{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Values: values},
			}})
			stmt.Clauses[sd.GenAuthKey()] = clause.Clause{}
		}

	}
}
