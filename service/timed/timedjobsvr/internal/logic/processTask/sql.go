package processTask

import (
	"context"
	"gitee.com/unitedrhino/core/service/timed/internal/domain"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
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
		Params:          task.Sql.Param.Sql,
		ResultCode:      e.GetCode(),
		ResultMsg:       e.GetMsg(),
		TimedTaskLogSql: &relationDB.TimedTaskLogSql{ExecNum: execNum},
	})
	if er != nil {
		logx.WithContext(ctx).Errorf("SqlExec.JobLog.Insert err:%v", er)
	}
	return err
}
