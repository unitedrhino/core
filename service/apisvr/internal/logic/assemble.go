package logic

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/unitedrhino/share/utils"
)

func ToTagsMap(tags []*types.Tag) map[string]string {
	if tags == nil {
		return nil
	}
	tagMap := make(map[string]string, len(tags))
	for _, tag := range tags {
		tagMap[tag.Key] = tag.Value
	}
	return tagMap
}

func ToTagsType(tags map[string]string) (retTag []*types.Tag) {
	for k, v := range tags {
		retTag = append(retTag, &types.Tag{
			Key:   k,
			Value: v,
		})
	}
	return
}

func ToSysPageRpc(in *types.PageInfo) *sys.PageInfo {
	return utils.Copy[sys.PageInfo](in)
}

func ToTimedJobPageRpc(in *types.PageInfo) *timedjob.PageInfo {
	return utils.Copy[timedjob.PageInfo](in)
}

func ToSysPointRpc(in *types.Point) *sys.Point {
	if in == nil {
		return nil
	}
	return &sys.Point{
		Longitude: in.Longitude,
		Latitude:  in.Latitude,
	}
}

func ToSysPointApi(in *sys.Point) *types.Point {
	if in == nil {
		return nil
	}
	return &types.Point{
		Longitude: in.Longitude,
		Latitude:  in.Latitude,
	}
}

func SysToWithIDTypes(in *sys.WithID) *types.WithID {
	if in == nil {
		return nil
	}
	return &types.WithID{
		ID: in.Id,
	}
}

func ToPageResp(p *types.PageInfo, total int64) types.PageResp {
	ret := types.PageResp{Total: total}
	if p == nil {
		return ret
	}
	ret.Page = p.Page
	ret.PageSize = p.Size
	return ret
}
