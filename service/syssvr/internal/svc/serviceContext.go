package svc

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/config"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/cache"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/caches"
	cas "gitee.com/i-Things/share/casbin"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/domain/tenant"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"os"
)

type ServiceContext struct {
	Config      config.Config
	ProjectID   *utils.SnowFlake
	AreaID      *utils.SnowFlake
	UserID      *utils.SnowFlake
	Casbin      *casbin.Enforcer
	Slot        *cache.Slot
	OssClient   *oss.Client
	Store       kv.Store
	PwdCheck    *cache.PwdCheck
	Captcha     *cache.Captcha
	Cm          *ClientsManage
	ServerMsg   *eventBus.FastEvent
	TenantCache *caches.Cache[tenant.Info]
	Sms         *clients.Sms
}

func NewServiceContext(c config.Config) *ServiceContext {
	//conn := sqlx.NewMysql(c.Database.DSN)
	stores.InitConn(c.Database)
	err := relationDB.Migrate(c.Database)
	if err != nil {
		logx.Error("syssvr 数据库初始化失败 err", err)
		os.Exit(-1)
	}
	// 自动迁移数据库
	db := stores.GetCommonConn(context.Background())
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	ProjectID := utils.NewSnowFlake(nodeID)
	AreaID := utils.NewSnowFlake(nodeID)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	UserID := utils.NewSnowFlake(nodeId)
	dbRaw, err := db.DB()
	if err != nil {
		logx.Error("sys failed to  database err: %v", err)
	}
	ca := cas.NewCasbinWithRedisWatcher(dbRaw, c.Database.DBType, c.CacheRedis[0].RedisConf)
	store := kv.NewStore(c.CacheRedis)
	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}
	serverMsg, err := eventBus.NewFastEvent(c.Event, c.Name, nodeID)
	logx.Must(err)
	sms, err := clients.NewSms(c.Sms)
	logx.Must(err)
	//sms.SendSms(clients.SendSmsParam{
	//	PhoneNumbers: []string{"17052709767"},
	//	SignName:     "萤科物联小程序",
	//	TemplateCode: "1842188",
	//	TemplateParam: map[string]any{
	//		"1": "123",
	//		"2": "333",
	//	},
	//})
	return &ServiceContext{
		ServerMsg: serverMsg,
		Captcha:   cache.NewCaptcha(store),
		PwdCheck:  cache.NewPwdCheck(store),
		Slot:      cache.NewSlot(),
		Cm:        NewClients(c),
		Config:    c,
		ProjectID: ProjectID,
		OssClient: ossClient,
		AreaID:    AreaID,
		UserID:    UserID,
		Casbin:    ca,
		Store:     store,
		Sms:       sms,
	}
}
