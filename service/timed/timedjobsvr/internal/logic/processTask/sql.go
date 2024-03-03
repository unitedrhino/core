package processTask

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/timed/internal/domain"
	"gitee.com/i-Things/core/service/timed/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/internal/logic/processTask/sqlFunc"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"github.com/dop251/goja"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func (t ProcessTask) SqlExec(ctx context.Context, task *domain.TaskInfo) error {
	var execNum int64
	err := func() error {
		dsn := cast.ToString(task.Env[domain.SqlEnvDsn])
		dbType := cast.ToString(task.Env[domain.SqlEnvDBType])
		if dsn == "" { //走默认值
			err := stores.GetCommonConn(ctx).Exec(task.Sql.Param.Sql).Error
			return stores.ErrFmt(err)
		}
		driver := cast.ToString(task.Env[domain.SqlEnvDriver])
		db, err := stores.GetConnDB(conf.Database{
			Driver: driver, //只支持这种模式
			DSN:    dsn,
			DBType: dbType,
		})
		if err != nil {
			return err
		}
		defer db.Close()
		ret, err := db.Exec(task.Sql.Param.Sql)
		if err != nil {
			return stores.ErrFmt(err)
		}
		execNum, _ = ret.RowsAffected()
		return nil
	}()
	e := errors.Fmt(err)
	er := relationDB.NewJobLogRepo(ctx).Insert(ctx, &relationDB.TimedTaskLog{
		Params:          task.Params,
		ResultCode:      e.GetCode(),
		ResultMsg:       e.GetMsg(),
		TimedTaskLogSql: &relationDB.TimedTaskLogSql{ExecNum: execNum},
	})
	if er != nil {
		logx.WithContext(ctx).Errorf("SqlExec.JobLog.Insert err:%v", er)
	}
	return err
}

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
		Params:             task.Params,
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
