package self

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/user/self"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
)

func ModuleIndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := self.NewModuleIndexLogic(r.Context(), svcCtx)
		resp, err := l.ModuleIndex()
		result.Http(w, r, resp, err)
	}
}
