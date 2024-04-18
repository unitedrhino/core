package common

import (
	"net/http"

	"gitee.com/i-Things/share/result"

	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/common"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
)

func WebsocketSubscribeIndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := common.NewWebsocketSubscribeIndexLogic(r.Context(), svcCtx)
		resp, err := l.WebsocketSubscribeIndex()
		result.Http(w, r, resp, err)
	}
}
