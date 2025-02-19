package svc

import (
	"gitee.com/unitedrhino/core/service/syssvr/internal/config"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/cache"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/timedjobdirect"
	"gitee.com/unitedrhino/core/share/domain/tenant"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/clients/dingClient"
	"gitee.com/unitedrhino/share/clients/smsClient"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/tools"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"sync"
)

type CaptchaLimit struct {
	PhoneIp      *tools.Limit
	PhoneAccount *tools.Limit
	EmailIp      *tools.Limit
	EmailAccount *tools.Limit
}

type LoginLimit struct {
	PwdIp      *tools.Limit
	PwdAccount *tools.Limit
}

type ServiceContext struct {
	Config             config.Config
	ProjectID          *utils.SnowFlake
	AreaID             *utils.SnowFlake
	UserID             *utils.SnowFlake
	Slot               *cache.Slot
	OssClient          *oss.Client
	Store              kv.Store
	Captcha            *cache.Captcha
	CaptchaLimit       CaptchaLimit
	LoginLimit         LoginLimit
	Cm                 *ClientsManage
	FastEvent          *eventBus.FastEvent
	UsersCache         *cache.UserCache
	TenantCache        *caches.Cache[tenant.Info, string]
	TenantConfigCache  *caches.Cache[sys.TenantConfig, string]
	ProjectCache       *caches.Cache[sys.ProjectInfo, int64]
	UserCache          *caches.Cache[sys.UserInfo, int64]
	AreaCache          *caches.Cache[sys.AreaInfo, int64]
	ApiCache           *caches.Cache[relationDB.SysApiInfo, string]
	RoleAccessCache    *caches.Cache[map[int64]struct{}, string]
	Sms                *smsClient.Sms
	DingStreamMap      map[string]*dingClient.StreamClient //key是租户号,value是需要同步的stream
	DingStreamMapMutex sync.RWMutex
	TimedM             timedmanage.TimedManage
	NodeID             int64
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		timedJob timedmanage.TimedManage
	)
	stores.InitConn(c.Database)
	err := relationDB.Migrate(c.Database)
	if err != nil {
		logx.Error("syssvr 数据库初始化失败 err", err)
		os.Exit(-1)
	}
	// 自动迁移数据库
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	ProjectID := utils.NewSnowFlake(nodeID)
	AreaID := utils.NewSnowFlake(nodeID)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	UserID := utils.NewSnowFlake(nodeId)
	store := kv.NewStore(c.CacheRedis)
	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}
	serverMsg, err := eventBus.NewFastEvent(c.Event, c.Name, nodeID)
	logx.Must(err)
	sms, err := smsClient.NewSms(c.Sms)
	logx.Must(err)

	cl := CaptchaLimit{
		PhoneIp:      tools.NewLimit(c.CaptchaPhoneIpLimit, "captcha", "phone:ip", config.DefaultIpLimit),
		PhoneAccount: tools.NewLimit(c.CaptchaPhoneIpLimit, "captcha", "phone:account", config.DefaultAccountLimit),
		EmailIp:      tools.NewLimit(c.CaptchaPhoneIpLimit, "captcha", "email:ip", config.DefaultIpLimit),
		EmailAccount: tools.NewLimit(c.CaptchaPhoneIpLimit, "captcha", "email:account", config.DefaultAccountLimit),
	}
	ll := LoginLimit{
		PwdIp:      tools.NewLimit(c.LoginPwdIpLimit, "login", "pwd:ip", config.DefaultIpLimit),
		PwdAccount: tools.NewLimit(c.LoginPwdAccountLimit, "login", "pwd:account", config.DefaultAccountLimit),
	}
	if c.TimedJobRpc.Enable {
		if c.TimedJobRpc.Mode == conf.ClientModeGrpc {
			timedJob = timedmanage.NewTimedManage(zrpc.MustNewClient(c.TimedJobRpc.Conf))
		} else {
			timedJob = timedjobdirect.NewTimedJob(c.TimedJobRpc.RunProxy)
		}
	}
	return &ServiceContext{
		FastEvent:     serverMsg,
		Captcha:       cache.NewCaptcha(store),
		Slot:          cache.NewSlot(),
		Cm:            NewClients(c),
		Config:        c,
		CaptchaLimit:  cl,
		LoginLimit:    ll,
		ProjectID:     ProjectID,
		OssClient:     ossClient,
		AreaID:        AreaID,
		UserID:        UserID,
		Store:         store,
		Sms:           sms,
		NodeID:        nodeID,
		TimedM:        timedJob,
		DingStreamMap: make(map[string]*dingClient.StreamClient),
	}
}
