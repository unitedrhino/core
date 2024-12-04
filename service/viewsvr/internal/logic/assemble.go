package logic

import (
	"gitee.com/unitedrhino/core/service/viewsvr/internal/types"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
)

func ToPageInfo(info *types.PageInfo) *stores.PageInfo {
	return utils.Copy[stores.PageInfo](info)
}
