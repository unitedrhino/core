package config

import (
	"net/http"

	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/notify/config"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TestSendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NotifyConfigTestSendReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := config.NewTestSendLogic(r.Context(), svcCtx)
		err := l.TestSend(&req)
		result.Http(w, r, nil, err)
	}
}
