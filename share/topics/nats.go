package topics

const (
	CoreProjectInfoDelete  = "server.core.project.info.delete"
	CoreAreaInfoDelete     = "server.core.area.info.delete"
	CoreOpsWorkOrderFinish = "server.core.ops.workOrder.finish"
	CoreSyncHalfHour       = "server.core.sync.halfHour" //半小时统计
	CoreSyncDay            = "server.core.sync.day"      //一天

	CoreTenantCreate = "server.core.tenant.create"

	CoreUserDelete     = "server.core.user.delete"
	CoreUserCreate     = "server.core.user.create"
	CoreUserUpdate     = "server.core.user.update"
	CoreProjectDelete  = "server.core.project.delete"
	CoreApiUserPublish = "server.core.api.user.publish.%v"

	ServerCacheKeySysUserInfo          = "cache:sys:user:info"
	ServerCacheKeySysUserTokenInfo     = "cache:sys:userToken:info"
	ServerCacheKeySysProjectInfo       = "cache:sys:project:info"
	ServerCacheKeySysAccessApi         = "cache:sys:access:api"
	ServerCacheKeySysRoleAccess        = "cache:sys:role:access"
	ServerCacheKeySysAreaInfo          = "cache:sys:area:info"
	ServerCacheKeySysTenantInfo        = "cache:sys:tenant:info"
	ServerCacheKeySysTenantConfig      = "cache:sys:tenant:config"
	ServerCacheKeySysTenantOpenWebhook = "cache:sys:tenant:open:webhook"
)
