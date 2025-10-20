package data

import (
	"net/http"

	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/user/data"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取项目权限列表
func ProjectIndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserDataProjectIndexReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := data.NewProjectIndexLogic(r.Context(), svcCtx)
		resp, err := l.ProjectIndex(&req)
		result.Http(w, r, resp, err)
	}
}
