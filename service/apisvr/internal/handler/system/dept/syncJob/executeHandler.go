package syncJob

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/dept/syncJob"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// 执行同步任务
func ExecuteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeptSyncJobExecuteReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := syncJob.NewExecuteLogic(r.Context(), svcCtx)
		err := l.Execute(&req)
		result.Http(w, r, nil, err)
	}
}
