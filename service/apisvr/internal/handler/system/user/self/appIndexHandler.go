package self

import (
	"net/http"

	"gitee.com/i-Things/core/shared/result"

	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/user/self"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
)

func AppIndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := self.NewAppIndexLogic(r.Context(), svcCtx)
		resp, err := l.AppIndex()
		result.Http(w, r, resp, err)
	}
}
