package opslogic

import (
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/dataType"
)

func ToOpsWorkOrderPo(in *sys.OpsWorkOrder) *relationDB.SysOpsWorkOrder {
	if in == nil {
		return nil
	}
	return &relationDB.SysOpsWorkOrder{
		ID:          in.Id,
		AreaID:      dataType.AreaID(in.AreaID),
		RaiseUserID: in.RaiseUserID,
		IssueDesc:   in.IssueDesc,
		Number:      in.Number,
		Type:        in.Type,
		Params:      in.Params,
		Status:      in.Status,
	}
}
