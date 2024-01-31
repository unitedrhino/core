package logic

import (
	"gitee.com/i-Things/core/service/viewsvr/internal/types"
	"gitee.com/i-Things/core/shared/def"
)

func ToPageInfo(info *types.PageInfo) *def.PageInfo {
	if info == nil {
		return nil
	}
	return &def.PageInfo{
		Page: info.Page,
		Size: info.Size,
	}
}
