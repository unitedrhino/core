package info

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/access/info"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"io"
	"net/http"
)

func MultiImportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AccessMultiImportReq
		req.Module = r.PostFormValue("module")
		f, _, err := r.FormFile("file")
		if err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("读取文件错误").AddDetail(err.Error()))
			return
		}
		defer f.Close()
		ff, err := io.ReadAll(f)
		if err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("读取文件错误").AddDetail(err.Error()))
		}
		l := info.NewMultiImportLogic(r.Context(), svcCtx)
		resp, err := l.MultiImport(&req, ff)
		result.Http(w, r, resp, err)
	}
}
