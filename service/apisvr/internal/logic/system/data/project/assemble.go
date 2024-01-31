package project

import (
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
)

func ToProjectApis(in []*sys.DataProject) (ret []*types.DataProject) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &types.DataProject{ProjectID: v.ProjectID, AuthType: v.AuthType})
	}
	return
}

func ToProjectPbs(in []*types.DataProject) (ret []*sys.DataProject) {
	if in == nil {
		return
	}
	for _, v := range in {
		ret = append(ret, &sys.DataProject{ProjectID: v.ProjectID, AuthType: v.AuthType})
	}
	return
}
