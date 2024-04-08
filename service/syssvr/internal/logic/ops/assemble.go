package opslogic

import (
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/stores"
)

func ToOpsWorkOrderPo(in *sys.OpsWorkOrder) *relationDB.SysOpsWorkOrder {
	if in == nil {
		return nil
	}
	return &relationDB.SysOpsWorkOrder{
		ID:          in.Id,
		AreaID:      stores.AreaID(in.AreaID),
		RaiseUserID: in.RaiseUserID,
		IssueDesc:   in.IssueDesc,
		Number:      in.Number,
		Type:        in.Type,
		Params:      in.Params,
		Status:      in.Status,
	}
}
