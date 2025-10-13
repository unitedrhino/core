package module

type Purpose = int64

const (
	PurposeNormal   = 1 //无需特殊处理
	PurposePlatform = 2 //那么只有default租户可以看,然后平台模块http头里不用传租户号
	PurposeProject  = 3 //需要选择项目,默认选择第一个
)
