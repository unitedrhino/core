package svc

import (
	"gitee.com/unitedrhino/core/service/timed/internal/repo/event/publish/pubJob"
	"gitee.com/unitedrhino/core/service/timed/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/config"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"os"
)

type ServiceContext struct {
	Config         config.Config
	Store          kv.Store
	Redis          *redis.Redis
	PubJob         *pubJob.PubJob
	AsynqClient    *asynq.Client
	AsynqInspector *asynq.Inspector
	FastEvent      *eventBus.FastEvent
	NodeID         int64
}

func NewServiceContext(c config.Config) *ServiceContext {
	pj, err := pubJob.NewPubJob(c.Event)
	if err != nil {
		logx.Error("初始化消息队列 err", err)
		os.Exit(-1)
	}
	stores.InitConn(c.Database)
	err = relationDB.Migrate(c.Database)
	if err != nil {
		logx.Error("timedjobsvr 数据库初始化失败 err", err)
		os.Exit(-1)
	}
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	serverMsg, err := eventBus.NewFastEvent(c.Event, c.Name, nodeID)
	if err != nil {
		logx.Errorf("NewFastEvent err cfg:%v err:%v", utils.Fmt(c.Event), err)
		os.Exit(-1)
	}
	return &ServiceContext{
		FastEvent:      serverMsg,
		Config:         c,
		PubJob:         pj,
		Redis:          redis.MustNewRedis(c.CacheRedis[0].RedisConf),
		AsynqClient:    clients.NewAsynqClient(c.CacheRedis),
		AsynqInspector: clients.NewAsynqInspector(c.CacheRedis),
		Store:          kv.NewStore(c.CacheRedis),
		NodeID:         nodeID,
	}
}
