package data

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/user/data"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// 获取区域权限列表
func AreaIndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserDataAreaIndexReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := data.NewAreaIndexLogic(r.Context(), svcCtx)
		resp, err := l.AreaIndex(&req)
		result.Http(w, r, resp, err)
	}
}
