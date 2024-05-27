package info

import (
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"
)

func ToAreaInfoTypes(root *sys.AreaInfo) *types.AreaInfo {
	if root == nil {
		return nil
	}
	api := &types.AreaInfo{
		CreatedTime:     root.CreatedTime,
		ProjectID:       root.ProjectID,
		ParentAreaID:    root.ParentAreaID,
		AreaID:          root.AreaID,
		AreaName:        root.AreaName,
		AreaNamePath:    root.AreaNamePath,
		LowerLevelCount: root.LowerLevelCount,
		AreaIDPath:      root.AreaIDPath,
		IsLeaf:          root.IsLeaf,
		UseBy:           root.UseBy,
		Position:        logic.ToSysPointApi(root.Position),
		Desc:            utils.ToNullString(root.Desc),
		Children:        nil,
	}
	if len(root.Children) > 0 {
		for _, child := range root.Children {
			api.Children = append(api.Children, ToAreaInfoTypes(child))
		}
	}
	return api
}
func ToAreaInfosTypes(in []*sys.AreaInfo) (ret []*types.AreaInfo) {
	for _, v := range in {
		ret = append(ret, ToAreaInfoTypes(v))
	}
	return
}
