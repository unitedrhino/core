package self

import (
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/user/self"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func UserSearchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserSearchReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := self.NewUserSearchLogic(r.Context(), svcCtx)
		resp, err := l.UserSearch(&req)
		result.Http(w, r, resp, err)
	}
}
