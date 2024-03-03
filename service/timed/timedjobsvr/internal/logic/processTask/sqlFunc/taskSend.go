package sqlFunc

import (
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/i-Things/share/errors"
	"github.com/dop251/goja"
	"github.com/gogf/gf/v2/util/gconv"
)

type TaskSend struct {
	Code        string            `json:"code"` //任务编码
	ExecContent string            `json:"execContent"`
	Param       map[string]string `json:"param"`
	Sql         string            `json:"sql"`
}

func (s *SqlFunc) TaskSendSqlJs() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		taskMap, ok := in.Arguments[0].Export().(map[string]any)
		if !ok {
			s.Errorf("timed.SetFunc.TaskSend script use err,"+
				"need an object,code:%v,script:%v",
				s.Task.Code, s.Task.Script.Param.ExecContent)
			panic(errors.Parameter.AddMsg("TaskSend param not rigth"))
		}
		var task TaskSend
		err := gconv.Struct(taskMap, &task)
		if err != nil {
			s.Errorf("timed.SetFunc.TaskSend gconv.Struct err:%v",
				err)
			panic(errors.Parameter.AddMsg("TaskSend param not rigth"))
		}
		err = s.TaskSend(s.ctx, &timedjob.TaskSendReq{
			GroupCode: s.Task.GroupCode,
			Code:      task.Code,
			ParamScript: func() *timedjob.TaskParamScript {
				if task.Param == nil {
					return nil
				}
				return &timedjob.TaskParamScript{Param: task.Param, ExecContent: task.ExecContent}
			}(),
			ParamSql: func() *timedjob.TaskParamSql {
				if task.Param == nil {
					return nil
				}
				return &timedjob.TaskParamSql{Sql: task.Sql}
			}(),
		})
		if err != nil {
			return s.vm.ToValue(ErrRet{Err: err})
		}
		return nil
	}

}
