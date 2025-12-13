package self

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/user/self"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
)

// 用户登出
func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := self.NewLogoutLogic(r.Context(), svcCtx)
		err := l.Logout()
		result.Http(w, r, nil, err)
	}
}
