package logic

import (
	"gitee.com/i-Things/core/service/viewsvr/internal/types"
	"gitee.com/i-Things/share/stores"
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
