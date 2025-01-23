package dept

type SyncDirection = int64

const (
	SyncDirectionFrom = 1 //上游同步到联犀(默认)
	SyncDirectionTo   = 2 //联犀同步到下游
)

type SyncMode = int64

const (
	SyncModeManual   SyncMode = 1 //手动模式
	SyncModeTiming   SyncMode = 2 //定时同步(半个小时)
	SyncModeRealTime SyncMode = 3 //实时同步
)
