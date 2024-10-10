package svc

import (
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/timedjobdirect"
	"gitee.com/unitedrhino/core/service/timed/timedschedulersvr/internal/config"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
)

type ServiceContext struct {
	Config       config.Config
	Scheduler    *clients.TimedScheduler
	Store        kv.Store
	SchedulerRun bool //只启动单例
	TimedM       timedmanage.TimedManage
	Queue        *clients.NatsClient
	NodeID       int64
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		TimedM timedmanage.TimedManage
	)

	stores.InitConn(c.Database)
	err := relationDB.Migrate(c.Database)
	if err != nil {
		logx.Error("初始化数据库错误 err", err)
		os.Exit(-1)
	}
	Scheduler := clients.NewTimedScheduler(c.CacheRedis)
	if c.TimedJobRpc.Enable {
		if c.TimedJobRpc.Mode == conf.ClientModeGrpc {
			TimedM = timedmanage.NewTimedManage(zrpc.MustNewClient(c.TimedJobRpc.Conf))
		} else {
			TimedM = timedjobdirect.NewTimedJob(c.TimedJobRpc.RunProxy)
		}
	}
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	queue, err := clients.NewNatsClient2(c.Event.Mode, c.Name, c.Event.Nats, nodeID)
	logx.Must(err)
	return &ServiceContext{
		Scheduler: Scheduler,
		Config:    c,
		Store:     kv.NewStore(c.CacheRedis),
		TimedM:    TimedM,
		Queue:     queue,
		NodeID:    nodeID,
	}
}
