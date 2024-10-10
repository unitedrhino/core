package logic

import (
	"gitee.com/unitedrhino/core/service/viewsvr/internal/types"
	"gitee.com/unitedrhino/share/stores"
)

func ToPageInfo(info *types.PageInfo) *stores.PageInfo {
	if info == nil {
		return nil
	}
	return &stores.PageInfo{
		Page: info.Page,
		Size: info.Size,
	}
}
