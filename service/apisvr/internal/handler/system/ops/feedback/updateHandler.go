package feedback

import (
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/ops/feedback"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// 更新帮助与反馈
func UpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpsFeedback
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := feedback.NewUpdateLogic(r.Context(), svcCtx)
		resp, err := l.Update(&req)
		result.Http(w, r, resp, err)
	}
}