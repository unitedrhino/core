package svc

import (
	"gitee.com/unitedrhino/core/service/datasvr/internal/config"
	"gitee.com/unitedrhino/core/service/datasvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/client/common"
	"gitee.com/unitedrhino/core/service/syssvr/client/log"
	role "gitee.com/unitedrhino/core/service/syssvr/client/rolemanage"
	tenant "gitee.com/unitedrhino/core/service/syssvr/client/tenantmanage"
	user "gitee.com/unitedrhino/core/service/syssvr/client/usermanage"
	"gitee.com/unitedrhino/core/service/syssvr/sysExport"
	"gitee.com/unitedrhino/core/service/syssvr/sysdirect"
	"gitee.com/unitedrhino/core/share/middlewares"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	CheckTokenWare rest.Middleware
	Slot           sysExport.SlotCacheT
	InitCtxsWare   rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	stores.InitConn(c.Database)
	logx.Must(relationDB.Migrate(c.Database))
	var (
		TenantRpc tenant.TenantManage
		LogRpc    log.Log
		UserRpc   user.UserManage
		AuthRpc   role.RoleManage
		Common    common.Common
	)
	if c.SysRpc.Mode == conf.ClientModeDirect {
		TenantRpc = sysdirect.NewTenantManage(c.SysRpc.RunProxy)
		LogRpc = sysdirect.NewLog(c.SysRpc.RunProxy)
		UserRpc = sysdirect.NewUser(c.SysRpc.RunProxy)
		AuthRpc = sysdirect.NewRole(c.SysRpc.RunProxy)
		Common = sysdirect.NewCommon(c.SysRpc.RunProxy)
	} else {
		TenantRpc = tenant.NewTenantManage(zrpc.MustNewClient(c.SysRpc.Conf))
		LogRpc = log.NewLog(zrpc.MustNewClient(c.SysRpc.Conf))
		UserRpc = user.NewUserManage(zrpc.MustNewClient(c.SysRpc.Conf))
		AuthRpc = role.NewRoleManage(zrpc.MustNewClient(c.SysRpc.Conf))
		Common = common.NewCommon(zrpc.MustNewClient(c.SysRpc.Conf))
	}
	Slot, err := sysExport.NewSlotCache(Common)
	logx.Must(err)
	return &ServiceContext{
		Config:         c,
		CheckTokenWare: middlewares.NewCheckTokenWareMiddleware(UserRpc, AuthRpc, TenantRpc, LogRpc).Handle,
		InitCtxsWare:   middlewares.InitMiddleware,
		Slot:           Slot,
	}
}
