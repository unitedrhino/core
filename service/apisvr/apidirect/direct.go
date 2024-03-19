package apidirect

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/config"
	"gitee.com/i-Things/core/service/apisvr/internal/handler"
	"gitee.com/i-Things/core/service/apisvr/internal/handler/system/proxy"
	"gitee.com/i-Things/core/service/apisvr/internal/repo/event/appDeviceEvent"
	"gitee.com/i-Things/core/service/apisvr/internal/repo/event/subApp"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"log"
	"net/http"
	"os"
)

type (
	Config         = config.Config
	ServiceContext = svc.ServiceContext
	ApiCtx         struct {
		Server *rest.Server
		SvcCtx *ServiceContext
	}
)

var (
	c config.Config
)

func NewApi(apiCtx ApiCtx) ApiCtx {
	conf.MustLoad("etc/api.yaml", &c)
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
	handler.RegisterWsHandlers(apiCtx.SvcCtx.Ws, ctx)
	subAppCli, err := subApp.NewSubApp(ctx.Config.Event)
	if err != nil {
		logx.Error("NewSubApp err", err)
		os.Exit(-1)
	}
	err = subAppCli.Subscribe(func(ctx1 context.Context) subApp.AppSubEvent {
		return appDeviceEvent.NewAppDeviceHandle(ctx1, ctx)
	})
	if err != nil {
		log.Fatalf("%v.subApp.Subscribe err:%v",
			utils.FuncName(), err)
	}
	return apiCtx
}
