package dataExport

import (
	"flag"
	"gitee.com/i-Things/core/service/datasvr/internal/config"
	"gitee.com/i-Things/core/service/datasvr/internal/handler"
	"gitee.com/i-Things/core/service/datasvr/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/data.yaml", "the config file")

func Run(server *rest.Server) *rest.Server {
	flag.Parse()
	var c config.Config
	err := conf.Load(*configFile, &c)
	if err != nil {
		return server
	}
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	return server
}
