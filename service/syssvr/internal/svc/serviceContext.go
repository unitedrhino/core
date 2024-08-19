package svc

import (
	"gitee.com/i-Things/core/service/syssvr/internal/config"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/cache"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/clients/smsClient"
	"gitee.com/i-Things/share/domain/tenant"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"os"
)

type ServiceContext struct {
	Config            config.Config
	ProjectID         *utils.SnowFlake
	AreaID            *utils.SnowFlake
	UserID            *utils.SnowFlake
	Slot              *cache.Slot
	OssClient         *oss.Client
	Store             kv.Store
	PwdCheck          *cache.PwdCheck
	Captcha           *cache.Captcha
	Cm                *ClientsManage
	FastEvent         *eventBus.FastEvent
	UserTokenInfo     *cache.UserToken
	TenantCache       *caches.Cache[tenant.Info, string]
	TenantConfigCache *caches.Cache[sys.TenantConfig, string]
	ProjectCache      *caches.Cache[sys.ProjectInfo, int64]
	UserCache         *caches.Cache[sys.UserInfo, int64]
	ApiCache          *caches.Cache[relationDB.SysApiInfo, string]
	RoleAccessCache   *caches.Cache[map[int64]struct{}, string]
	Sms               *smsClient.Sms
}

func NewServiceContext(c config.Config) *ServiceContext {
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
	userTokenInfo, err := cache.NewUserToken(serverMsg)
	logx.Must(err)
	return &ServiceContext{
		FastEvent:     serverMsg,
		Captcha:       cache.NewCaptcha(store),
		PwdCheck:      cache.NewPwdCheck(store),
		Slot:          cache.NewSlot(),
		Cm:            NewClients(c),
		Config:        c,
		ProjectID:     ProjectID,
		OssClient:     ossClient,
		AreaID:        AreaID,
		UserID:        UserID,
		Store:         store,
		Sms:           sms,
		UserTokenInfo: userTokenInfo,
	}
}
