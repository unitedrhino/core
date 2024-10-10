package processTask

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/timed/internal/domain"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/logic/processTask/sqlFunc"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/dop251/goja"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func (t ProcessTask) ScriptExec(ctx context.Context, task *domain.TaskInfo) error {
	var (
		code = errors.OK.Code
		msg  = errors.OK.GetMsg()
	)
	vm := goja.New()
	sf := sqlFunc.NewSqlFunc(ctx, t.svcCtx, task, t.TaskSend, vm)

	func() {
		defer func() {
			if p := recover(); p != nil {
				if e, ok := p.(error); ok {
					logx.WithContext(ctx).Errorf("panic err:%v", e)
					err := errors.Fmt(e)
					sf.ExecuteLog = append(sf.ExecuteLog, &domain.ScriptLog{
						Level:       "error",
						Content:     fmt.Sprintf("catch an panic,err:%v", err.Error()),
						CreatedTime: time.Now().Unix(),
					})
					code = err.GetCode()
					msg = err.GetMsg()
				}
			}
		}()
		var Run func(map[string]string) map[string]any
		err := func() error {
			err := sf.Register()
			if err != nil {
				return err
			}
			_, err = vm.RunString(task.Script.Param.ExecContent)
			if err != nil {
				return err
			}
			err = vm.ExportTo(vm.Get("Main"), &Run)
			if err != nil {
				return err
			}
			return nil
		}()

		e := errors.Fmt(err)
		if e != nil {
			logx.WithContext(ctx).Errorf("ScriptExec.err:%v", e)
			code = e.GetCode()
			msg = e.GetMsg()
		} else if Run != nil {
			ret := Run(task.Script.Param.Param)
			code = cast.ToInt64(ret["code"])
			msg = cast.ToString(ret["msg"])
		}
	}()
	er := relationDB.NewJobLogRepo(ctx).Insert(ctx, &relationDB.TimedTaskLog{
		GroupCode:          task.GroupCode,
		TaskCode:           task.Code,
		Params:             utils.MarshalNoErr(task.Script.Param.Param),
		ResultCode:         code,
		ResultMsg:          msg,
		TimedTaskLogScript: &relationDB.TimedTaskLogScript{ExecLog: sf.ExecuteLog},
		TimedTaskLogSql: &relationDB.TimedTaskLogSql{
			SelectNum: sf.SelectNum,
			ExecNum:   sf.ExecNum,
		},
	})
	if er != nil {
		logx.WithContext(ctx).Errorf("SqlExec.JobLog.Insert err:%v", er)
	}
	return er
}
