package timedschedulerdirect

import (
	"flag"
	"fmt"
	"gitee.com/unitedrhino/core/service/timed/timedschedulersvr/internal/config"
	schedulerServer "gitee.com/unitedrhino/core/service/timed/timedschedulersvr/internal/server/timedscheduler"
	"gitee.com/unitedrhino/core/service/timed/timedschedulersvr/internal/startup"
	"gitee.com/unitedrhino/core/service/timed/timedschedulersvr/internal/svc"
	"gitee.com/unitedrhino/core/service/timed/timedschedulersvr/pb/timedscheduler"
	"gitee.com/unitedrhino/share/interceptors"
	"gitee.com/unitedrhino/share/services"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"sync"
)

type Config = config.Config

var (
	svcCtx     *svc.ServiceContext
	svcOnce    sync.Once
	runSvrOnce sync.Once
	c          config.Config
)

func GetSvcCtx() *svc.ServiceContext {
	svcOnce.Do(func() {
		flag.Parse()
		utils.ConfMustLoad("etc/timedscheduler.yaml", &c)
		svcCtx = svc.NewServiceContext(c)
		startup.Init(svcCtx)
		logx.Infof("enabled timedschedulersvr")
	})
	return svcCtx
}

// RunServer 如果是直连模式,同时提供Grpc的能力
func RunServer(svcCtx *svc.ServiceContext) {
	runSvrOnce.Do(func() {
		go Run(svcCtx)
	})
}

func Run(svcCtx *svc.ServiceContext) {
	c := svcCtx.Config
	s := services.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		timedscheduler.RegisterTimedschedulerServer(grpcServer, schedulerServer.NewTimedschedulerServer(svcCtx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(interceptors.Ctxs, interceptors.Error)
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
