package self

import (
	"net/http"

	"gitee.com/i-Things/share/result"

	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/user/self"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
)

func MessageStatisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := self.NewMessageStatisticsLogic(r.Context(), svcCtx)
		resp, err := l.MessageStatistics()
		result.Http(w, r, resp, err)
	}
}
