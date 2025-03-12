package logic

import "gitee.com/unitedrhino/core/service/datasvr/internal/types"

func ToPageResp(p *types.PageInfo, total int64) types.PageResp {
	ret := types.PageResp{Total: total}
	if p == nil {
		return ret
	}
	ret.Page = p.Page
	ret.PageSize = p.Size
	return ret
}
