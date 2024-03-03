package sqlFunc

import (
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"github.com/dop251/goja"
)

func (s *SqlFunc) Select() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 1 {
			s.Errorf("timed.SetFunc.Select script use err,"+
				"need (第一个参数是sql 第二个参数是指定的数据库连接code(可选,不填选择默认的连接,需要在config里配置),code:%v,script:%v",
				s.Task.Code, s.Task.Script.Param.ExecContent)
			panic(errors.Parameter)
		}
		sql := in.Arguments[0].String()
		db, Close := s.getConn(in, "select")
		defer Close()
		var ret []map[string]any
		err := stores.QueryContext(s.ctx, db, sql, &ret)
		if err != nil {
			panic(errors.Database.AddDetail(err))
		}
		s.SelectNum += int64(len(ret))

		return s.vm.ToValue(ret)
	}

}
