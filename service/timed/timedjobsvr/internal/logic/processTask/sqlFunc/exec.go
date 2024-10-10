package sqlFunc

import (
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/dop251/goja"
)

func (s *SqlFunc) Exec() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 1 {
			s.Errorf("timed.SetFunc.Exec script use err,"+
				"need (第一个参数是sql 第二个参数是指定的数据库连接code(可选,不填选择默认的连接,需要在config里配置),code:%v,script:%v",
				s.Task.Code, s.Task.Script.Param.ExecContent)
			panic(errors.Parameter)
		}
		sql := in.Arguments[0].String()
		db, close := s.getConn(in, "exec")
		defer close()
		ret, err := db.Exec(sql)
		if err != nil {
			return s.vm.ToValue(ErrRet{Err: stores.ErrFmt(err)})
		}
		RowsAffected, err := ret.RowsAffected()
		if err != nil {
			return s.vm.ToValue(ErrRet{Err: stores.ErrFmt(err)})
		}
		s.ExecNum += RowsAffected
		return s.vm.ToValue(ErrRet{})
	}

}
