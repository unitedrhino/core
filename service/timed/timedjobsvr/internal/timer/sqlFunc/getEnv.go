package sqlFunc

import (
	"gitee.com/i-Things/share/errors"
	"github.com/dop251/goja"
)

func (s *SqlFunc) GetEnv() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 1 {
			s.Errorf("timed.SetFunc.GetEnv script use err,need (key string),code:%v,script:%v",
				s.Task.Code, s.Task.Sql.Param.ExecContent)
			panic(errors.Parameter)
		}
		return s.vm.ToValue(s.Task.Env[in.Arguments[0].String()])
	}
}
