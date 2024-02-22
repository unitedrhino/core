package dictmanagelogic

import (
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"
)

func ToDictInfoPb(in *relationDB.SysDictInfo, children []*relationDB.SysDictInfo) *sys.DictInfo {
	if in == nil {
		return nil
	}
	ret := &sys.DictInfo{
		Id:       in.ID,
		Name:     in.Name,
		Type:     in.Type,
		Desc:     utils.ToRpcNullString(in.Desc),
		ParentID: in.ParentID,
		Children: ToDictInfosPb(in.Children),
		Body:     utils.ToRpcNullString(in.Body),
		Details:  ToDictDetailsPb(in.Details),
		IdPath:   utils.GetIDPath(in.IDPath),
	}
	if children != nil {
		var idMap = map[int64][]*sys.DictInfo{}
		for _, v := range children {
			idMap[v.ParentID] = append(idMap[v.ParentID], ToDictInfoPb(v, nil))
		}
		fillDictInfoChildren(ret, idMap)
	}
	return ret
}

func fillDictInfoChildren(node *sys.DictInfo, nodeMap map[int64][]*sys.DictInfo) {
	// 找到当前节点的子节点数组
	children := nodeMap[node.Id]
	for _, child := range children {
		fillDictInfoChildren(child, nodeMap)
	}
	node.Children = children
}

func ToDictInfosPb(in []*relationDB.SysDictInfo) (list []*sys.DictInfo) {
	for _, v := range in {
		list = append(list, ToDictInfoPb(v, nil))
	}
	return
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
