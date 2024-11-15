package project

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/data/project"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// 批量删除授权项目权限
func MultiDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DataProjectMultiDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := project.NewMultiDeleteLogic(r.Context(), svcCtx)
		err := l.MultiDelete(&req)
		result.Http(w, r, nil, err)
	}
}
