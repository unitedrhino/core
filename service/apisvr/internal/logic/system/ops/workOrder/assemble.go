package workOrder

import (
	"gitee.com/i-Things/core/service/apisvr/internal/types"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"
)

func ToOpsWorkOrderPb(in *types.OpsWorkOrder) *sys.OpsWorkOrder {
	return utils.Copy[sys.OpsWorkOrder](in)
}
