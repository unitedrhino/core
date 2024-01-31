package area

import (
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/shared/def"
)

var (
	RootNode = relationDB.SysAreaInfo{
		ProjectID:    0,
		AreaID:       def.RootNode,
		ParentAreaID: 0,
		AreaName:     "全部区域",
		Desc:         "虚拟根节点",
	}
	NotClassified = relationDB.SysAreaInfo{
		ProjectID:    0,
		AreaID:       def.NotClassified,
		ParentAreaID: 0,
		AreaName:     "未分类",
		Desc:         "虚拟根节点",
	}
)
