package info

import (
	"gitee.com/unitedrhino/core/service/datasvr/internal/logic/data/staticstics/info"
	"gitee.com/unitedrhino/core/service/datasvr/internal/svc"
	"gitee.com/unitedrhino/core/service/datasvr/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func ReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer utils.Recover(r.Context())
		var req types.StaticsticsInfoReadReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := info.NewReadLogic(r.Context(), svcCtx)
		resp, err := l.Read(&req)
		result.Http(w, r, resp, err)
	}
}
