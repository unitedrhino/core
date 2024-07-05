package sysdirect

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/internal/config"
	accessmanageServer "gitee.com/i-Things/core/service/syssvr/internal/server/accessmanage"
	appmanageServer "gitee.com/i-Things/core/service/syssvr/internal/server/appmanage"
	areamanageServer "gitee.com/i-Things/core/service/syssvr/internal/server/areamanage"
	commonServer "gitee.com/i-Things/core/service/syssvr/internal/server/common"
	datamanageServer "gitee.com/i-Things/core/service/syssvr/internal/server/datamanage"
	logServer "gitee.com/i-Things/core/service/syssvr/internal/server/log"
	modulemanageServer "gitee.com/i-Things/core/service/syssvr/internal/server/modulemanage"
	opsServer "gitee.com/i-Things/core/service/syssvr/internal/server/ops"
	projectmanageServer "gitee.com/i-Things/core/service/syssvr/internal/server/projectmanage"
	rolemanageServer "gitee.com/i-Things/core/service/syssvr/internal/server/rolemanage"
	tenantmanageServer "gitee.com/i-Things/core/service/syssvr/internal/server/tenantmanage"
	usermanageServer "gitee.com/i-Things/core/service/syssvr/internal/server/usermanage"

	"gitee.com/i-Things/core/service/syssvr/internal/startup"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/interceptors"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"sync"
)

type Config = config.Config

var (
	ctxSvc     *svc.ServiceContext
	svcOnce    sync.Once
	runSvrOnce sync.Once
	c          config.Config
)

func GetSvcCtx() *svc.ServiceContext {
	svcOnce.Do(func() {
		conf.MustLoad("etc/sys.yaml", &c)
		ctxSvc = svc.NewServiceContext(c)
		startup.Init(ctxSvc)
		logx.Infof("enabled syssvr")
	})
	return ctxSvc
}

// RunServer 如果是直连模式,同时提供Grpc的能力
func RunServer(svcCtx *svc.ServiceContext) {
	runSvrOnce.Do(func() {
		utils.Go(context.Background(), func() {
			Run(svcCtx)
		})
	})
}

func Run(svcCtx *svc.ServiceContext) {
	c := svcCtx.Config
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		sys.RegisterUserManageServer(grpcServer, usermanageServer.NewUserManageServer(svcCtx))
		sys.RegisterAccessManageServer(grpcServer, accessmanageServer.NewAccessManageServer(svcCtx))
		sys.RegisterRoleManageServer(grpcServer, rolemanageServer.NewRoleManageServer(svcCtx))
		sys.RegisterAppManageServer(grpcServer, appmanageServer.NewAppManageServer(svcCtx))
		sys.RegisterModuleManageServer(grpcServer, modulemanageServer.NewModuleManageServer(svcCtx))
		sys.RegisterCommonServer(grpcServer, commonServer.NewCommonServer(svcCtx))
		sys.RegisterLogServer(grpcServer, logServer.NewLogServer(svcCtx))
		sys.RegisterProjectManageServer(grpcServer, projectmanageServer.NewProjectManageServer(svcCtx))
		sys.RegisterAreaManageServer(grpcServer, areamanageServer.NewAreaManageServer(svcCtx))
		sys.RegisterTenantManageServer(grpcServer, tenantmanageServer.NewTenantManageServer(svcCtx))
		sys.RegisterDataManageServer(grpcServer, datamanageServer.NewDataManageServer(svcCtx))
		sys.RegisterOpsServer(grpcServer, opsServer.NewOpsServer(svcCtx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(interceptors.Ctxs, interceptors.Error)
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
