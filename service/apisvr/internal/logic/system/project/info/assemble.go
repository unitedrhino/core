package info

import (
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"
)

func ToProjectPb(in *types.ProjectInfo) *sys.ProjectInfo {
	if in == nil {
		return nil
	}
	return &sys.ProjectInfo{
		ProjectID:   in.ProjectID,
		ProjectName: in.ProjectName,
		AdminUserID: in.AdminUserID,
		Position:    logic.ToSysPointRpc(in.Position),
		Desc:        utils.ToRpcNullString(in.Desc),
	}
}
