package common

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/common"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// 获取第三方的部门信息,Children只能获取一层,需要递归获取
func ThirdDeptReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ThirdDeptInfoReadReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := common.NewThirdDeptReadLogic(r.Context(), svcCtx)
		resp, err := l.ThirdDeptRead(&req)
		result.Http(w, r, resp, err)
	}
}
