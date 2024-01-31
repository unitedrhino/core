package common

import (
	"net/http"

	"gitee.com/i-Things/share/result"

	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/common"
	"gitee.com/i-Things/core/service/apisvr/internal/svc"
)

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := common.NewUploadFileLogic(r.Context(), svcCtx, r)
		resp, err := l.UploadFile()
		result.Http(w, r, resp, err)
	}
}
