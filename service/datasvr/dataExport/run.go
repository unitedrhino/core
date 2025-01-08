package dataExport

import (
	"flag"
	"gitee.com/unitedrhino/core/service/apisvr/apidirect"
	"gitee.com/unitedrhino/core/service/datasvr/internal/config"
	"gitee.com/unitedrhino/core/service/datasvr/internal/handler"
	"gitee.com/unitedrhino/core/service/datasvr/internal/svc"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/rest"
)

func init() {
	apidirect.RegisterServer(func(server *rest.Server) error {
		flag.Parse()
		var c config.Config
		utils.ConfMustLoad("etc/data.yaml", &c)
		ctx := svc.NewServiceContext(c)
		handler.RegisterHandlers(server, ctx)
		return nil
	})
}

func Run(server *rest.Server) *rest.Server {
	flag.Parse()
	var c config.Config
	utils.ConfMustLoad("etc/data.yaml", &c)
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	return server
}
