package svc

import (
	"gitee.com/i-Things/core/service/apisvr/export"
	"gitee.com/i-Things/core/service/apisvr/internal/config"
	"gitee.com/i-Things/core/service/apisvr/internal/middleware"
	"gitee.com/i-Things/core/service/syssvr/client/accessmanage"
	app "gitee.com/i-Things/core/service/syssvr/client/appmanage"
	"gitee.com/i-Things/core/service/syssvr/client/areamanage"
	"gitee.com/i-Things/core/service/syssvr/client/common"
	"gitee.com/i-Things/core/service/syssvr/client/datamanage"
	"gitee.com/i-Things/core/service/syssvr/client/dictmanage"
	"gitee.com/i-Things/core/service/syssvr/client/log"
	module "gitee.com/i-Things/core/service/syssvr/client/modulemanage"
	"gitee.com/i-Things/core/service/syssvr/client/notifymanage"
	"gitee.com/i-Things/core/service/syssvr/client/ops"
	"gitee.com/i-Things/core/service/syssvr/client/projectmanage"
	role "gitee.com/i-Things/core/service/syssvr/client/rolemanage"
	tenant "gitee.com/i-Things/core/service/syssvr/client/tenantmanage"
	user "gitee.com/i-Things/core/service/syssvr/client/usermanage"
	"gitee.com/i-Things/core/service/syssvr/sysExport"
	"gitee.com/i-Things/core/service/syssvr/sysdirect"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/client/timedmanage"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/timedjobdirect"
	"gitee.com/i-Things/core/service/timed/timedschedulersvr/client/timedscheduler"
	"gitee.com/i-Things/core/service/timed/timedschedulersvr/timedschedulerdirect"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/oss"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/share/verify"
	ws "gitee.com/i-Things/share/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"time"
)

type SvrClient struct {
	TenantRpc tenant.TenantManage
	UserRpc   user.UserManage
	RoleRpc   role.RoleManage
	AppRpc    app.AppManage
	ModuleRpc module.ModuleManage
	LogRpc    log.Log
	ProjectM  projectmanage.ProjectManage
	AreaM     areamanage.AreaManage
	DictM     dictmanage.DictManage
	NotifyM   notifymanage.NotifyManage
	Common    common.Common

	AccessRpc      accessmanage.AccessManage
	DataM          datamanage.DataManage
	Timedscheduler timedscheduler.Timedscheduler
	TimedJob       timedmanage.TimedManage
	Ops            ops.Ops
}

type ServiceContext struct {
	SvrClient
	Ws        *ws.Server
	Config    config.Config
	UserCache sysExport.UserCacheT

	CheckTokenWare rest.Middleware
	DataAuthWare   rest.Middleware
	TeardownWare   rest.Middleware
	CheckApiWare   rest.Middleware
	InitCtxsWare   rest.Middleware
	Captcha        *verify.Captcha
	OssClient      *oss.Client
	NodeID         int64
	ServerMsg      *eventBus.FastEvent
}

func NewServiceContext(c config.Config) *ServiceContext {
	var (
		appRpc        app.AppManage
		projectM      projectmanage.ProjectManage
		areaM         areamanage.AreaManage
		sysCommon     common.Common
		timedSchedule timedscheduler.Timedscheduler
		timedJob      timedmanage.TimedManage
		tenantM       tenant.TenantManage
		DataM         datamanage.DataManage
		DictM         dictmanage.DictManage
		accessM       accessmanage.AccessManage
		Ops           ops.Ops
		NotifyM       notifymanage.NotifyManage
	)
	var ur user.UserManage
	var ro role.RoleManage
	var me module.ModuleManage
	var lo log.Log

	caches.InitStore(c.CacheRedis)
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	serverMsg, err := eventBus.NewFastEvent(c.Event, c.Name, nodeID)
	logx.Must(err)
	ws.StartWsDp(false, nodeID, serverMsg, c.CacheRedis)
	if c.SysRpc.Enable {
		if c.SysRpc.Mode == conf.ClientModeGrpc {
			projectM = projectmanage.NewProjectManage(zrpc.MustNewClient(c.SysRpc.Conf))
			areaM = areamanage.NewAreaManage(zrpc.MustNewClient(c.SysRpc.Conf))
			ur = user.NewUserManage(zrpc.MustNewClient(c.SysRpc.Conf))
			ro = role.NewRoleManage(zrpc.MustNewClient(c.SysRpc.Conf))
			me = module.NewModuleManage(zrpc.MustNewClient(c.SysRpc.Conf))
			lo = log.NewLog(zrpc.MustNewClient(c.SysRpc.Conf))
			sysCommon = common.NewCommon(zrpc.MustNewClient(c.SysRpc.Conf))
			appRpc = app.NewAppManage(zrpc.MustNewClient(c.SysRpc.Conf))
			tenantM = tenant.NewTenantManage(zrpc.MustNewClient(c.SysRpc.Conf))
			DataM = datamanage.NewDataManage(zrpc.MustNewClient(c.SysRpc.Conf))
			accessM = accessmanage.NewAccessManage(zrpc.MustNewClient(c.SysRpc.Conf))
			DictM = dictmanage.NewDictManage(zrpc.MustNewClient(c.SysRpc.Conf))
			Ops = ops.NewOps(zrpc.MustNewClient(c.SysRpc.Conf))
			NotifyM = notifymanage.NewNotifyManage(zrpc.MustNewClient(c.SysRpc.Conf))
		} else {
			projectM = sysdirect.NewProjectManage(c.SysRpc.RunProxy)
			areaM = sysdirect.NewAreaManage(c.SysRpc.RunProxy)
			ur = sysdirect.NewUser(c.SysRpc.RunProxy)
			ro = sysdirect.NewRole(c.SysRpc.RunProxy)
			me = sysdirect.NewModule(c.SysRpc.RunProxy)
			lo = sysdirect.NewLog(c.SysRpc.RunProxy)
			sysCommon = sysdirect.NewCommon(c.SysRpc.RunProxy)
			appRpc = sysdirect.NewApp(c.SysRpc.RunProxy)
			tenantM = sysdirect.NewTenantManage(c.SysRpc.RunProxy)
			DataM = sysdirect.NewData(c.SysRpc.RunProxy)
			accessM = sysdirect.NewAccess(c.SysRpc.RunProxy)
			DictM = sysdirect.NewDict(c.SysRpc.RunProxy)
			Ops = sysdirect.NewOps(c.SysRpc.RunProxy)
			NotifyM = sysdirect.NewNotify(c.SysRpc.RunProxy)
		}
	}

	if c.TimedSchedulerRpc.Enable {
		if c.TimedSchedulerRpc.Mode == conf.ClientModeGrpc {
			timedSchedule = timedscheduler.NewTimedscheduler(zrpc.MustNewClient(c.TimedSchedulerRpc.Conf))
		} else {
			timedSchedule = timedschedulerdirect.NewScheduler(c.TimedSchedulerRpc.RunProxy)
		}
	}
	if c.TimedJobRpc.Enable {
		if c.TimedJobRpc.Mode == conf.ClientModeGrpc {
			timedJob = timedmanage.NewTimedManage(zrpc.MustNewClient(c.TimedJobRpc.Conf))
		} else {
			timedJob = timedjobdirect.NewTimedJob(c.TimedJobRpc.RunProxy)
		}
	}

	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}
	userCache, err := sysExport.NewUserInfoCache(ur, serverMsg)
	logx.Must(err)
	captcha := verify.NewCaptcha(c.Captcha.ImgHeight, c.Captcha.ImgWidth,
		c.Captcha.KeyLong, c.CacheRedis, time.Duration(c.Captcha.KeepTime)*time.Second)
	return &ServiceContext{
		Config:         c,
		CheckTokenWare: export.NewCheckTokenWareMiddleware(ur, ro, tenantM).Handle,
		DataAuthWare:   middleware.NewDataAuthWareMiddleware(c).Handle,
		TeardownWare:   middleware.NewTeardownWareMiddleware(c, lo).Handle,
		CheckApiWare:   middleware.NewCheckApiWareMiddleware().Handle,
		InitCtxsWare:   ctxs.InitMiddleware,
		UserCache:      userCache,
		Captcha:        captcha,
		OssClient:      ossClient,
		Ws:             ws.MustNewServer(c.RestConf),
		NodeID:         nodeID,
		ServerMsg:      serverMsg,
		SvrClient: SvrClient{
			TenantRpc:      tenantM,
			AppRpc:         appRpc,
			UserRpc:        ur,
			RoleRpc:        ro,
			ModuleRpc:      me,
			LogRpc:         lo,
			AccessRpc:      accessM,
			NotifyM:        NotifyM,
			Timedscheduler: timedSchedule,
			TimedJob:       timedJob,
			ProjectM:       projectM,
			AreaM:          areaM,
			DataM:          DataM,
			DictM:          DictM,
			Common:         sysCommon,
			Ops:            Ops,
		},
		//OSS:        ossClient,
	}
}
