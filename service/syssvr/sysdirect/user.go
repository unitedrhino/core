package sysdirect

import (
	"gitee.com/unitedrhino/core/service/syssvr/client/ops"
	client "gitee.com/unitedrhino/core/service/syssvr/client/usermanage"
	opsServer "gitee.com/unitedrhino/core/service/syssvr/internal/server/ops"
	server "gitee.com/unitedrhino/core/service/syssvr/internal/server/usermanage"

	clientNotify "gitee.com/unitedrhino/core/service/syssvr/client/notifymanage"
	serverNotify "gitee.com/unitedrhino/core/service/syssvr/internal/server/notifymanage"

	clientDeptM "gitee.com/unitedrhino/core/service/syssvr/client/departmentmanage"
	serverDeptM "gitee.com/unitedrhino/core/service/syssvr/internal/server/departmentmanage"

	clientRole "gitee.com/unitedrhino/core/service/syssvr/client/rolemanage"
	serverRole "gitee.com/unitedrhino/core/service/syssvr/internal/server/rolemanage"

	clientAccess "gitee.com/unitedrhino/core/service/syssvr/client/accessmanage"
	serverAccess "gitee.com/unitedrhino/core/service/syssvr/internal/server/accessmanage"

	clientData "gitee.com/unitedrhino/core/service/syssvr/client/datamanage"
	serverData "gitee.com/unitedrhino/core/service/syssvr/internal/server/datamanage"

	clientDict "gitee.com/unitedrhino/core/service/syssvr/client/dictmanage"
	serverDict "gitee.com/unitedrhino/core/service/syssvr/internal/server/dictmanage"

	clientModule "gitee.com/unitedrhino/core/service/syssvr/client/modulemanage"
	serverModule "gitee.com/unitedrhino/core/service/syssvr/internal/server/modulemanage"

	clientLog "gitee.com/unitedrhino/core/service/syssvr/client/log"
	serverLog "gitee.com/unitedrhino/core/service/syssvr/internal/server/log"

	clientCommon "gitee.com/unitedrhino/core/service/syssvr/client/common"
	serverCommon "gitee.com/unitedrhino/core/service/syssvr/internal/server/common"

	clientApp "gitee.com/unitedrhino/core/service/syssvr/client/appmanage"
	serverApp "gitee.com/unitedrhino/core/service/syssvr/internal/server/appmanage"

	clientTenant "gitee.com/unitedrhino/core/service/syssvr/client/tenantmanage"
	serverTenant "gitee.com/unitedrhino/core/service/syssvr/internal/server/tenantmanage"

	clientProject "gitee.com/unitedrhino/core/service/syssvr/client/projectmanage"
	serverProject "gitee.com/unitedrhino/core/service/syssvr/internal/server/projectmanage"

	clientArea "gitee.com/unitedrhino/core/service/syssvr/client/areamanage"
	serverArea "gitee.com/unitedrhino/core/service/syssvr/internal/server/areamanage"
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

func NewOps(runSvr bool) ops.Ops {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return ops.NewDirectOps(svcCtx, opsServer.NewOpsServer(svcCtx))
}

func NewNotify(runSvr bool) clientNotify.NotifyManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientNotify.NewDirectNotifyManage(svcCtx, serverNotify.NewNotifyManageServer(svcCtx))
}

func NewDeptM(runSvr bool) clientDeptM.DepartmentManage {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return clientDeptM.NewDirectDepartmentManage(svcCtx, serverDeptM.NewDepartmentManageServer(svcCtx))
}
