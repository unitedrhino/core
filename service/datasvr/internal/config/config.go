package config

import (
	"gitee.com/unitedrhino/share/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Database   conf.Database
	CacheRedis cache.ClusterConf
	SysRpc     conf.RpcClientConf `json:",optional"`
}
