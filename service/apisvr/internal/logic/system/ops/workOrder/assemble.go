package workOrder

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

func ToOpsWorkOrderPb(in *types.OpsWorkOrder) *sys.OpsWorkOrder {
	return utils.Copy[sys.OpsWorkOrder](in)
}
