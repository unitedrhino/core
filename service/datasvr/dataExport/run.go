package dataExport

import (
	"flag"
	"gitee.com/i-Things/core/service/datasvr/internal/config"
	"gitee.com/i-Things/core/service/datasvr/internal/handler"
	"gitee.com/i-Things/core/service/datasvr/internal/svc"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/rest"
)

func Run(server *rest.Server) *rest.Server {
	flag.Parse()
	var c config.Config
	utils.ConfMustLoad("etc/data.yaml", &c)
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	return server
}
