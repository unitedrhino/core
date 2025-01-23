package ops

type WorkOrderStatus = int64

const (
	WorkOrderStatusWait     WorkOrderStatus = 1 //待处理
	WorkOrderStatusHandling WorkOrderStatus = 2 //处理中
	WorkOrderStatusFinished WorkOrderStatus = 3 //处理完成
)

type WorkOrderType = string

const (
	WorkOrderTypeDeviceMaintenance = "deviceMaintenance" //设备维修工单
	WorkOrderTypeSceneAlarm        = "sceneAlarm"        //场景自动化报警
)
