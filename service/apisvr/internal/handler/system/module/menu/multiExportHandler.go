package menu

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/module/menu"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// 批量导出菜单
func MultiExportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MenuMultiExportReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := menu.NewMultiExportLogic(r.Context(), svcCtx)
		resp, err := l.MultiExport(&req)
		result.Http(w, r, resp, err)
	}
}
