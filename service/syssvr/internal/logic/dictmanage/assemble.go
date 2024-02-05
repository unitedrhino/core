package dictmanagelogic

import (
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"
)

func ToDictInfoPb(in *relationDB.SysDictInfo) *sys.DictInfo {
	if in == nil {
		return nil
	}
	return &sys.DictInfo{
		Id:      in.ID,
		Name:    in.Name,
		Type:    in.Type,
		Desc:    utils.ToRpcNullString(in.Desc),
		Status:  in.Status,
		Body:    utils.ToRpcNullString(in.Body),
		Details: ToDictDetailsPb(in.Details),
	}

}

func ToDictDetailsPb(in []*relationDB.SysDictDetail) []*sys.DictDetail {
	var list []*sys.DictDetail
	for _, v := range in {
		list = append(list, &sys.DictDetail{
			Id:     v.ID,
			DictID: v.DictID,
			Label:  v.Label,
			Value:  v.Value,
			Extend: v.Extend,
			Sort:   v.Sort,
			Desc:   utils.ToRpcNullString(v.Desc),
			Status: v.Status,
			Body:   utils.ToRpcNullString(v.Body),
		})
	}
	return list
}
