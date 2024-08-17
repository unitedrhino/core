package middleware

import (
	"gitee.com/i-Things/core/service/apisvr/internal/config"
	operLog "gitee.com/i-Things/core/service/syssvr/client/log"
	"net/http"
	"sync"
)

type TeardownWareMiddleware struct {
	cfg    config.Config
	LogRpc operLog.Log
}

var respPool sync.Pool
var bufferSize = 1024

func init() {
	respPool.New = func() interface{} {
		return make([]byte, bufferSize)
	}
}

func NewTeardownWareMiddleware(cfg config.Config, LogRpc operLog.Log) *TeardownWareMiddleware {
	return &TeardownWareMiddleware{cfg: cfg, LogRpc: LogRpc}
}

func (m *TeardownWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//logx.WithContext(r.Context()).Infof("%s.Lifecycle.Before", utils.FuncName())

		next(w, r)

		//logx.WithContext(r.Context()).Infof("%s.Lifecycle.After", utils.FuncName())
	}
}
