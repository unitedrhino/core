package svc

import (
	"gitee.com/unitedrhino/core/service/apisvr/exportMiddleware"
	"gitee.com/unitedrhino/core/service/datasvr/internal/config"
	"gitee.com/unitedrhino/core/service/datasvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/client/common"
	"gitee.com/unitedrhino/core/service/syssvr/client/log"
	role "gitee.com/unitedrhino/core/service/syssvr/client/rolemanage"
	tenant "gitee.com/unitedrhino/core/service/syssvr/client/tenantmanage"
	user "gitee.com/unitedrhino/core/service/syssvr/client/usermanage"
	"gitee.com/unitedrhino/core/service/syssvr/sysExport"
	"gitee.com/unitedrhino/share/ctxs"
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
	var ur user.UserManage
	var ro role.RoleManage
	ur = user.NewUserManage(zrpc.MustNewClient(c.SysRpc.Conf))
	ro = role.NewRoleManage(zrpc.MustNewClient(c.SysRpc.Conf))
	tm := tenant.NewTenantManage(zrpc.MustNewClient(c.SysRpc.Conf))
	lo := log.NewLog(zrpc.MustNewClient(c.SysRpc.Conf))
	Slot, err := sysExport.NewSlotCache(common.NewCommon(zrpc.MustNewClient(c.SysRpc.Conf)))
	logx.Must(err)
	return &ServiceContext{
		Config:         c,
		CheckTokenWare: exportMiddleware.NewCheckTokenWareMiddleware(ur, ro, tm, lo).Handle,
		InitCtxsWare:   ctxs.InitMiddleware,
		Slot:           Slot,
	}
}
