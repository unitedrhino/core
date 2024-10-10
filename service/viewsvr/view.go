package main

import (
	"flag"
	"fmt"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/viewsvr/internal/config"
	"gitee.com/unitedrhino/core/service/viewsvr/internal/handler"
	"gitee.com/unitedrhino/core/service/viewsvr/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func main() {
	flag.Parse()

	var c config.Config
	utils.ConfMustLoad("etc/view.yaml", &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	server.PrintRoutes()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
