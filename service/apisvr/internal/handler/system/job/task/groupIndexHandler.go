package task

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/job/task"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func GroupIndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TimedTaskGroupIndexReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := task.NewGroupIndexLogic(r.Context(), svcCtx)
		resp, err := l.GroupIndex(&req)
		result.Http(w, r, resp, err)
	}
}
