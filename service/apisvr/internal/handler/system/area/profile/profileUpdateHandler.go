package profile

import (
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/area/profile"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func ProfileUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AreaProfile
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := profile.NewProfileUpdateLogic(r.Context(), svcCtx)
		err := l.ProfileUpdate(&req)
		result.Http(w, r, nil, err)
	}
}