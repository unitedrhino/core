package dict

import (
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"
)

func ToDetailPb(in *types.DictDetail) *sys.DictDetail {
	if in == nil {
		return nil
	}
	return &sys.DictDetail{
		Id:     in.ID,
		DictID: in.DictID,
		Label:  in.Label,
		Value:  in.Value,
		Extend: in.Extend,
		Sort:   in.Sort,
		Desc:   utils.ToRpcNullString(in.Desc),
		Status: in.Status,
		Body:   utils.ToRpcNullString(in.Body),
	}
}

func ToInfoPb(in *types.DictInfo) *sys.DictInfo {
	if in == nil {
		return nil
	}
	return &sys.DictInfo{
		Id:   in.ID,
		Name: in.Name,
		Type: in.Type,
		Desc: utils.ToRpcNullString(in.Desc),
		Body: utils.ToRpcNullString(in.Body),
	}
}

func ToDetailTypes(in *sys.DictDetail) *types.DictDetail {
	if in == nil {
		return nil
	}
	return &types.DictDetail{
		ID:     in.Id,
		DictID: in.DictID,
		Label:  in.Label,
		Value:  in.Value,
		Extend: in.Extend,
		Sort:   in.Sort,
		Desc:   utils.ToNullString(in.Desc),
		Status: in.Status,
		Body:   utils.ToNullString(in.Body),
	}
}

func ToDetailsTypes(in []*sys.DictDetail) (ret []*types.DictDetail) {
	for _, v := range in {
		ret = append(ret, ToDetailTypes(v))
	}
	return
}

func ToInfoTypes(in *sys.DictInfo) *types.DictInfo {
	if in == nil {
		return nil
	}
	return &types.DictInfo{
		ID:       in.Id,
		Name:     in.Name,
		Type:     in.Type,
		Desc:     utils.ToNullString(in.Desc),
		Body:     utils.ToNullString(in.Body),
		Details:  ToDetailsTypes(in.Details),
		Children: ToInfosTypes(in.Children),
	}
}
func ToInfosTypes(in []*sys.DictInfo) (ret []*types.DictInfo) {
	for _, v := range in {
		ret = append(ret, ToInfoTypes(v))
	}
	return
}
