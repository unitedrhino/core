package svc

import (
	"gitee.com/unitedrhino/core/service/viewsvr/internal/config"
	"gitee.com/unitedrhino/core/service/viewsvr/internal/middleware"
	"gitee.com/unitedrhino/core/service/viewsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"os"
)

type ServiceContext struct {
	Config         config.Config
	SetupWare      rest.Middleware
	CheckTokenWare rest.Middleware
	DataAuthWare   rest.Middleware
	TeardownWare   rest.Middleware
	CheckApiWare   rest.Middleware
	OssClient      *oss.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	stores.InitConn(c.Database)
	logx.Must(relationDB.Migrate(c.Database))
	ossClient, err := oss.NewOssClient(c.OssConf)
	if err != nil {
		logx.Errorf("NewOss err err:%v", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		Config:         c,
		SetupWare:      middleware.NewSetupWareMiddleware().Handle,
		CheckTokenWare: middleware.NewCheckTokenWareMiddleware().Handle,
		DataAuthWare:   middleware.NewDataAuthWareMiddleware().Handle,
		TeardownWare:   middleware.NewTeardownWareMiddleware().Handle,
		CheckApiWare:   middleware.NewCheckApiWareMiddleware().Handle,
		OssClient:      ossClient,
	}
}
