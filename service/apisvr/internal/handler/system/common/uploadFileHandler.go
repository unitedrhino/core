package common

import (
	"gitee.com/unitedrhino/share/utils"
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/common"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
)

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer utils.Recover(r.Context())
		l := common.NewUploadFileLogic(r.Context(), svcCtx, r)
		resp, err := l.UploadFile()
		result.Http(w, r, resp, err)
	}
}
