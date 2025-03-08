package coreExport

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/config"
	"gitee.com/unitedrhino/core/service/apisvr/internal/handler"
	"gitee.com/unitedrhino/core/service/apisvr/internal/handler/system/proxy"
	"gitee.com/unitedrhino/core/service/apisvr/internal/startup"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	_ "gitee.com/unitedrhino/core/service/datasvr/dataExport"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/services"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

type (
	ApiCtx struct {
		Server *rest.Server
		SvcCtx *svc.ServiceContext
	}
)

var (
	c config.Config
)

func NewApi(apiCtx ApiCtx) ApiCtx {
	utils.ConfMustLoad("etc/api.yaml", &c)
	apiCtx = runApi(apiCtx)
	return apiCtx
}

func runApi(apiCtx ApiCtx) ApiCtx {
	var server = apiCtx.Server
	ctx := svc.NewServiceContext(c)
	apiCtx.SvcCtx = ctx
	if server == nil {
		server = rest.MustNewServer(c.RestConf, rest.WithCustomCors(func(header http.Header) {
			header.Set("Access-Control-Allow-Headers", ctxs.HttpAllowHeader)
			header.Set("Access-Control-Allow-Origin", "*")
		}, nil, "*"),
			rest.WithNotFoundHandler(proxy.Handler(ctx)),
		)
		apiCtx.Server = server
	}
	handler.RegisterHandlers(server, ctx)
	//handler.RegisterWsHandlers(apiCtx.SvcCtx.Ws, ctx)
	startup.Init(ctx)
	services.InitApisvrs(server)
	//dataExport.Run(server)
	return apiCtx
}
