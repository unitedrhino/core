package self

import (
	"net/http"

	"gitee.com/i-Things/share/result"

	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/user/self"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
)

func AccessTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := self.NewAccessTreeLogic(r.Context(), svcCtx)
		resp, err := l.AccessTree()
		result.Http(w, r, resp, err)
	}
}
