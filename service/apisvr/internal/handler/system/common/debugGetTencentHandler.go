package common

import (
	"net/http"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
)

func DebugGetTencentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(r.Header.Get("echostr")))
		//result.Http(w, r, nil, err)
	}
}
