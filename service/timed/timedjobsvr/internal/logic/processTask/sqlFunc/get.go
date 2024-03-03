package sqlFunc

import (
	"gitee.com/i-Things/share/errors"
	"github.com/dop251/goja"
)

func (s *SqlFunc) Get() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 1 {
			s.Errorf("timed.SetFunc.Get script use err,need (key string),code:%v,script:%v",
				s.Task.Code, s.Task.Script.Param.ExecContent)
			panic(errors.Parameter)
		}
		ret, err := s.SvcCtx.Store.GetCtx(s.ctx, s.GetStringKey(in.Arguments[0].String()))
		if err != nil {
			s.Errorf("timed.SetFunc.Get script Store.GetCtx err:%v", err)
			panic(errors.Database.AddDetail(err))
		}
		return s.vm.ToValue(ret)
	}
}
