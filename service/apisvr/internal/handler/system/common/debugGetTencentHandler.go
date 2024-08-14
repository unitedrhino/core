package common

import (
	"net/http"

	"gitee.com/i-Things/share/result"

	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/common"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
)

func DebugGetTencentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := common.NewDebugGetTencentLogic(r.Context(), svcCtx)
		err := l.DebugGetTencent()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(r.Header.Get("echostr")))
		result.Http(w, r, nil, err)
	}
}
