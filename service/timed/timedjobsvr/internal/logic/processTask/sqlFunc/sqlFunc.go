package sqlFunc

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/unitedrhino/core/service/timed/internal/domain"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/domain/task"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/dop251/goja"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

type TaskSendFunc func(context.Context, *timedjob.TaskSendReq) error
type SqlFunc struct {
	SvcCtx     *svc.ServiceContext
	ctx        context.Context
	Task       *domain.TaskInfo
	vm         *goja.Runtime
	ExecuteLog []*domain.ScriptLog
	TaskSend   TaskSendFunc
	SelectNum  int64 //查询的数量
	ExecNum    int64 //执行的数量
	kvKeyPre   string
	logx.Logger
}

func NewSqlFunc(ctx context.Context, svcCtx *svc.ServiceContext, task *domain.TaskInfo, TaskSend TaskSendFunc, vm *goja.Runtime) *SqlFunc {
	kvKeyPre := fmt.Sprintf("timed:sql:%s:", task.GroupCode)
	if code := task.Env["code"]; code != "" {
		kvKeyPre = fmt.Sprintf("timed:sql:%s:", task.GroupCode)
	}
	return &SqlFunc{SvcCtx: svcCtx, ctx: ctx, Logger: logx.WithContext(ctx), Task: task, TaskSend: TaskSend, vm: vm, kvKeyPre: kvKeyPre}
}

func (s *SqlFunc) Register() error {
	var funcList = []struct {
		Name string
		f    func(in goja.FunctionCall) goja.Value
	}{
		{"Set", s.Set()},
		{"Get", s.Get()},
		{"Select", s.Select()},
		{"Exec", s.Exec()},
		{"LogError", s.LogError()},
		{"LogInfo", s.LogInfo()},
		{"GetEnv", s.GetEnv()},
		{"Hexists", s.Hexists()},
		{"Hdel", s.Hdel()},
		{"Hget", s.Hget()},
		{"Hset", s.Hset()},
		{"HgetAll", s.HGetAll()},
		{"CreateOne", s.CreateOne()},
		{"TaskSendSqlJs", s.TaskSendSqlJs()},
	}
	for _, f := range funcList {
		err := s.vm.Set(f.Name, f.f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SqlFunc) getConn(in goja.FunctionCall, tp string) (*sql.DB, func() error) {
	dsn := s.Task.Env[task.SqlEnvDsn]
	dbType := s.Task.Env[task.SqlEnvDBType]
	driver := s.Task.Env[task.SqlEnvDriver]
	if len(in.Arguments) > 1 {
		dbName := in.Arguments[1].String()
		c, ok := s.Task.Script.Config.Database[dbName]
		if ok {
			dsn = c.DSN
			dbType = c.DBType
			driver = c.Driver
		}
	}
	if dsn == "" { //判断系统配置
		c, ok := s.Task.Script.Config.Database[tp]
		if ok {
			dsn = c.DSN
			dbType = c.DBType
			driver = c.Driver
		} else {
			db, _ := stores.GetCommonConn(s.ctx).DB()
			return db, func() error {
				return nil
			}
		}
	}
	fmt.Println(driver)
	db, err := func() (*sql.DB, error) {
		switch dbType {
		case conf.Tdengine:
			td, err := clients.NewTDengine(conf.TSDB{
				DBType: dbType,
				Driver: driver,
				DSN:    dsn,
			})
			if err != nil {
				return nil, err
			}
			return td.DB, nil
		default:
			conn, err := stores.GetConn(conf.Database{
				DBType: dbType,
				DSN:    dsn,
			})
			if err != nil {
				return nil, err
			}
			return conn.DB()
		}
	}()
	if err != nil {
		panic(errors.Database.AddMsgf("getConn.GetConn failure dsn:%v dbType:%v err:%v", dsn, dbType, err))
	}
	return db, db.Close
}
func (s *SqlFunc) GetStringKey(key string) string {
	return s.kvKeyPre + "string:" + key
}
func (s *SqlFunc) GetHashKey(key string) string {
	return s.kvKeyPre + "hash:" + key
}

func (s *SqlFunc) GetHashField(field string) string {
	date := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s:%s", date, field)
}
func (s *SqlFunc) GetHashFieldWithDay(field string, day int) string {
	date := time.Now().Add(time.Hour * 24 * time.Duration(day)).Format("2006-01-02")
	return fmt.Sprintf("%s:%s", date, field)
}
func (s *SqlFunc) ToRealHashField(field string) string {
	_, ret, find := strings.Cut(field, ":")
	if !find {
		return field
	}
	return ret
}
