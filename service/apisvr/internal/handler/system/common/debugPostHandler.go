package common

import (
	"gitee.com/unitedrhino/share/utils"
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/common"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
)

func DebugPostHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := common.NewDebugPostLogic(r.Context(), svcCtx)
		resp, err := l.DebugPost(r)
		l.Infof("DebugPost resp:%v err:%v", utils.Fmt(resp), err)
		result.Http(w, r, resp, err)
	}
}
