package area

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/area/info"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"
)

func ToDataAreaDetail(ctx context.Context, svcCtx *svc.ServiceContext, in []*sys.DataArea, areaMap map[int64]*sys.AreaInfo) (ret []*types.DataAreaDetail) {
	if in == nil {
		return
	}
	for _, v := range in {
		var ui *types.UserCore
		if svcCtx != nil && v.TargetType == def.TargetUser {
			u, err := svcCtx.UserCache.GetData(ctx, v.TargetID)
			if err != nil {
				continue
			}
			ui = utils.Copy[types.UserCore](u)
		}
		ret = append(ret, &types.DataAreaDetail{User: ui, TargetType: v.TargetType, TargetID: v.TargetID, AuthType: v.AuthType, IsAuthChildren: v.IsAuthChildren, AreaInfo: info.ToAreaInfoTypes(areaMap[v.AreaID])})
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
