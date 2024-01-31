package svc

import (
	"gitee.com/i-Things/core/service/timed/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/timed/timedschedulersvr/internal/config"
	"gitee.com/i-Things/core/shared/clients"
	"gitee.com/i-Things/core/shared/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"os"
)

type ServiceContext struct {
	Config       config.Config
	Scheduler    *clients.TimedScheduler
	Store        kv.Store
	SchedulerRun bool //只启动单例
}

func NewServiceContext(c config.Config) *ServiceContext {
	stores.InitConn(c.Database)
	err := relationDB.Migrate(c.Database)
	if err != nil {
		logx.Error("初始化数据库错误 err", err)
		os.Exit(-1)
	}
	Scheduler := clients.NewTimedScheduler(c.CacheRedis)
	return &ServiceContext{
		Scheduler: Scheduler,
		Config:    c,
		Store:     kv.NewStore(c.CacheRedis),
	}
}
