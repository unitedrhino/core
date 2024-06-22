package info

import (
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/project/info"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func ProfileReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WithCode
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := info.NewProfileReadLogic(r.Context(), svcCtx)
		resp, err := l.ProfileRead(&req)
		result.Http(w, r, resp, err)
	}
}
