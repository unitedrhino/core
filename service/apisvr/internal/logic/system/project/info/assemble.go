package info

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

func ToProjectPb(in *types.ProjectInfo) *sys.ProjectInfo {
	return utils.Copy[sys.ProjectInfo](in)

}
