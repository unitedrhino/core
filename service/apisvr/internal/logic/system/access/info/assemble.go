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
func ToAccessModuleInfoTypes(in []*sys.AccessInfo) (ret []*types.AccessModuleInfo) {
	var retMap = map[string]map[string][]*types.AccessInfo{}
	for _, v := range in {
		_, ok := retMap[v.Module]
		if !ok {
			retMap[v.Module] = map[string][]*types.AccessInfo{}
		}

		retMap[v.Module][v.Group] = append(retMap[v.Module][v.Group], ToAccessInfoTypes(v))
	}
	var retList []*types.AccessModuleInfo
	var moduleID int64

	for k, v := range retMap {
		moduleID++
		code := fmt.Sprintf("module%d", moduleID)
		var groups []*types.AccessGroupInfo
		var groupID int64
		for gk, gv := range v {
			groupID++
			code := fmt.Sprintf("group%d", groupID)
			groups = append(groups, &types.AccessGroupInfo{
				ID:       code,
				Code:     code,
				Name:     gk,
				Children: gv,
			})
		}
		retList = append(retList, &types.AccessModuleInfo{
			ID:       code,
			Code:     code,
			Name:     k,
			Children: groups,
		})
	}
	return retList
}
