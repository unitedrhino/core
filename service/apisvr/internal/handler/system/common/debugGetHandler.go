package common

import (
	"gitee.com/i-Things/share/utils"
	"net/http"

	"gitee.com/i-Things/share/result"

	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/common"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
)

func DebugGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := common.NewDebugGetLogic(r.Context(), svcCtx)
		resp, err := l.DebugGet(r)
		l.Infof("DebugGet resp:%v err:%v", utils.Fmt(resp), err)
		result.Http(w, r, resp, err)
	}
}
