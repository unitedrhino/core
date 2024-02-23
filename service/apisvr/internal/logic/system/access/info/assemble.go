package info

import (
	"fmt"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/access/api"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
)

func ToAccessInfoPb(in *types.AccessInfo) *sys.AccessInfo {
	if in == nil {
		return nil
	}
	return &sys.AccessInfo{
		Id:         in.ID,
		Code:       in.Code,
		Name:       in.Name,
		Group:      in.Group,
		IsNeedAuth: in.IsNeedAuth,
		Desc:       in.Desc,
	}
}

func ToAccessInfoTypes(in *sys.AccessInfo) *types.AccessInfo {
	if in == nil {
		return nil
	}
	return &types.AccessInfo{
		ID:         in.Id,
		Code:       in.Code,
		Name:       in.Name,
		Group:      in.Group,
		IsNeedAuth: in.IsNeedAuth,
		Desc:       in.Desc,
		Apis:       api.ToApiInfosTypes(in.Apis),
	}
}
func ToAccessInfosTypes(in []*sys.AccessInfo) (ret []*types.AccessInfo) {
	for _, v := range in {
		ret = append(ret, ToAccessInfoTypes(v))
	}
	return
}
func ToAccessGroupInfoTypes(in []*sys.AccessInfo) (ret []*types.AccessGroupInfo) {
	var retMap = map[string][]*types.AccessInfo{}
	for _, v := range in {
		retMap[v.Group] = append(retMap[v.Group], ToAccessInfoTypes(v))
	}
	var retList []*types.AccessGroupInfo
	var groupID int64
	for k, v := range retMap {
		groupID++
		code := fmt.Sprintf("group%d", groupID)
		retList = append(retList, &types.AccessGroupInfo{
			ID:       code,
			Code:     code,
			Name:     k,
			Children: v,
		})
	}
	return retList
}
