package sqlFunc

import (
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	sq "gitee.com/unitedrhino/squirrel"
	"github.com/dop251/goja"
)

type CreateOneRet struct {
	Err          error `json:"err"`
	LastInsertId int64 `json:"lastInsertId"` //最后更新的id
	RowsAffected int64 `json:"rowsAffected"` //受影响的行数
}

func (s *SqlFunc) CreateOne() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 2 {
			s.Errorf("timed.SetFunc.Exec script use err,"+
				"need (第一个参数是sql 第二个参数是指定的数据库连接code(可选,不填选择默认的连接,需要在config里配置),code:%v,script:%v",
				s.Task.Code, s.Task.Script.Param.ExecContent)
			panic(errors.Parameter)
		}
		table := in.Arguments[0].String()
		data := in.Arguments[1].Export()
		db, Close := s.getConn(in, "exec")
		defer Close()
		var columns []string
		var values []any
		datas := utils.ToStringMap(data)
		for k, v := range datas {
			columns = append(columns, k)
			values = append(values, v)
		}
		sql, vals, err := sq.Insert(table).Columns(columns...).Values(values...).ToSql()
		if err != nil {
			return ToJsStu(s.vm, CreateOneRet{Err: errors.Fmt(err)})
		}
		ret, err := db.Exec(sql, vals...)
		if err != nil {
			return ToJsStu(s.vm, CreateOneRet{Err: errors.Fmt(err)})
		}
		//ret := conn.Table(table).CreateOne(data)
		//	err := ret.Error
		RowsAffected, _ := ret.RowsAffected()
		LastInsertId, _ := ret.LastInsertId()
		s.ExecNum += RowsAffected
		if err != nil {
			return ToJsStu(s.vm, CreateOneRet{Err: errors.Fmt(err)})
		}
		return ToJsStu(s.vm, CreateOneRet{RowsAffected: RowsAffected, LastInsertId: LastInsertId})
	}
}
