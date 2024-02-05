package sysdirect

import (
	client "gitee.com/i-Things/core/service/syssvr/client/usermanage"
	server "gitee.com/i-Things/core/service/syssvr/internal/server/usermanage"

	clientRole "gitee.com/i-Things/core/service/syssvr/client/rolemanage"
	serverRole "gitee.com/i-Things/core/service/syssvr/internal/server/rolemanage"

	clientAccess "gitee.com/i-Things/core/service/syssvr/client/accessmanage"
	serverAccess "gitee.com/i-Things/core/service/syssvr/internal/server/accessmanage"

	clientData "gitee.com/i-Things/core/service/syssvr/client/datamanage"
	serverData "gitee.com/i-Things/core/service/syssvr/internal/server/datamanage"

	clientDict "gitee.com/i-Things/core/service/syssvr/client/dictmanage"
	serverDict "gitee.com/i-Things/core/service/syssvr/internal/server/dictmanage"

	clientModule "gitee.com/i-Things/core/service/syssvr/client/modulemanage"
	serverModule "gitee.com/i-Things/core/service/syssvr/internal/server/modulemanage"

	clientLog "gitee.com/i-Things/core/service/syssvr/client/log"
	serverLog "gitee.com/i-Things/core/service/syssvr/internal/server/log"

	clientCommon "gitee.com/i-Things/core/service/syssvr/client/common"
	serverCommon "gitee.com/i-Things/core/service/syssvr/internal/server/common"

	clientApp "gitee.com/i-Things/core/service/syssvr/client/appmanage"
	serverApp "gitee.com/i-Things/core/service/syssvr/internal/server/appmanage"

	clientTenant "gitee.com/i-Things/core/service/syssvr/client/tenantmanage"
	serverTenant "gitee.com/i-Things/core/service/syssvr/internal/server/tenantmanage"

	clientProject "gitee.com/i-Things/core/service/syssvr/client/projectmanage"
	serverProject "gitee.com/i-Things/core/service/syssvr/internal/server/projectmanage"

	clientArea "gitee.com/i-Things/core/service/syssvr/client/areamanage"
	serverArea "gitee.com/i-Things/core/service/syssvr/internal/server/areamanage"
)

func NewUser(runSvr bool) client.UserManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return client.NewDirectUserManage(svcCtx, server.NewUserManageServer(svcCtx))
}

func NewRole(runSvr bool) clientRole.RoleManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientRole.NewDirectRoleManage(svcCtx, serverRole.NewRoleManageServer(svcCtx))
}
func NewAccess(runSvr bool) clientAccess.AccessManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientAccess.NewDirectAccessManage(svcCtx, serverAccess.NewAccessManageServer(svcCtx))
}

func NewData(runSvr bool) clientData.DataManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientData.NewDirectDataManage(svcCtx, serverData.NewDataManageServer(svcCtx))
}

func NewDict(runSvr bool) clientDict.DictManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientDict.NewDirectDictManage(svcCtx, serverDict.NewDictManageServer(svcCtx))
}

func NewModule(runSvr bool) clientModule.ModuleManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientModule.NewDirectModuleManage(svcCtx, serverModule.NewModuleManageServer(svcCtx))
}

func NewCommon(runSvr bool) clientCommon.Common {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientCommon.NewDirectCommon(svcCtx, serverCommon.NewCommonServer(svcCtx))
}

func NewLog(runSvr bool) clientLog.Log {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientLog.NewDirectLog(svcCtx, serverLog.NewLogServer(svcCtx))
}

func NewApp(runSvr bool) clientApp.AppManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientApp.NewDirectAppManage(svcCtx, serverApp.NewAppManageServer(svcCtx))
}

func NewTenantManage(runSvr bool) clientTenant.TenantManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientTenant.NewDirectTenantManage(svcCtx, serverTenant.NewTenantManageServer(svcCtx))
}

func NewProjectManage(runSvr bool) clientProject.ProjectManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientProject.NewDirectProjectManage(svcCtx, serverProject.NewProjectManageServer(svcCtx))
}
func NewAreaManage(runSvr bool) clientArea.AreaManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientArea.NewDirectAreaManage(svcCtx, serverArea.NewAreaManageServer(svcCtx))
}
