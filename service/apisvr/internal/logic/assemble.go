package logic

import (
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/core/service/timed/timedjobsvr/pb/timedjob"
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
	if in == nil {
		return nil
	}
	return &sys.PageInfo{
		Page: in.Page,
		Size: in.Size,
	}
}

func ToTimedJobPageRpc(in *types.PageInfo) *timedjob.PageInfo {
	if in == nil {
		return nil
	}
	return &timedjob.PageInfo{
		Page: in.Page,
		Size: in.Size,
	}
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
