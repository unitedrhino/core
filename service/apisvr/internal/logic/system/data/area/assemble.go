package area

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/area/info"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
)

func ToDataAreaDetail(in []*sys.DataArea, areaMap map[int64]*sys.AreaInfo) (ret []*types.DataAreaDetail) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &types.DataAreaDetail{AuthType: v.AuthType, IsAuthChildren: v.IsAuthChildren, AreaInfo: info.ToAreaInfoTypes(areaMap[v.AreaID])})
	}
	return
}

func ToAreaPbs(in []*types.DataArea) (ret []*sys.DataArea) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &sys.DataArea{AreaID: v.AreaID, AuthType: v.AuthType, IsAuthChildren: v.IsAuthChildren})
	}
	return
}
