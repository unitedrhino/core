package self

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/user/self"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func AreaApplyCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserAreaApplyCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := self.NewAreaApplyCreateLogic(r.Context(), svcCtx)
		err := l.AreaApplyCreate(&req)
		result.Http(w, r, nil, err)
	}
}
