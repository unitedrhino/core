package svc

import (
	"gitee.com/i-Things/core/service/apisvr/export"
	"gitee.com/i-Things/core/service/datasvr/internal/config"
	"gitee.com/i-Things/core/service/datasvr/internal/repo/relationDB"
	role "gitee.com/i-Things/core/service/syssvr/client/rolemanage"
	tenant "gitee.com/i-Things/core/service/syssvr/client/tenantmanage"
	user "gitee.com/i-Things/core/service/syssvr/client/usermanage"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	CheckTokenWare rest.Middleware
	InitCtxsWare   rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	stores.InitConn(c.Database)
	logx.Must(relationDB.Migrate(c.Database))
	var ur user.UserManage
	var ro role.RoleManage
	ur = user.NewUserManage(zrpc.MustNewClient(c.SysRpc.Conf))
	ro = role.NewRoleManage(zrpc.MustNewClient(c.SysRpc.Conf))
	tm := tenant.NewTenantManage(zrpc.MustNewClient(c.SysRpc.Conf))
	return &ServiceContext{
		Config:         c,
		CheckTokenWare: export.NewCheckTokenWareMiddleware(ur, ro, tm).Handle,
		InitCtxsWare:   ctxs.InitMiddleware,
	}
}
