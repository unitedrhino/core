// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.1
// Source: sys.proto

package departmentmanage

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AccessInfo                            = sys.AccessInfo
	AccessInfoIndexReq                    = sys.AccessInfoIndexReq
	AccessInfoIndexResp                   = sys.AccessInfoIndexResp
	AccessInfoMultiImportReq              = sys.AccessInfoMultiImportReq
	AccessInfoMultiImportResp             = sys.AccessInfoMultiImportResp
	ApiInfo                               = sys.ApiInfo
	ApiInfoIndexReq                       = sys.ApiInfoIndexReq
	ApiInfoIndexResp                      = sys.ApiInfoIndexResp
	AppInfo                               = sys.AppInfo
	AppInfoIndexReq                       = sys.AppInfoIndexReq
	AppInfoIndexResp                      = sys.AppInfoIndexResp
	AppModuleIndexReq                     = sys.AppModuleIndexReq
	AppModuleIndexResp                    = sys.AppModuleIndexResp
	AppModuleMultiUpdateReq               = sys.AppModuleMultiUpdateReq
	AreaInfo                              = sys.AreaInfo
	AreaInfoIndexReq                      = sys.AreaInfoIndexReq
	AreaInfoIndexResp                     = sys.AreaInfoIndexResp
	AreaInfoReadReq                       = sys.AreaInfoReadReq
	AreaProfile                           = sys.AreaProfile
	AreaProfileIndexReq                   = sys.AreaProfileIndexReq
	AreaProfileIndexResp                  = sys.AreaProfileIndexResp
	AreaProfileReadReq                    = sys.AreaProfileReadReq
	AreaWithID                            = sys.AreaWithID
	AuthApiInfo                           = sys.AuthApiInfo
	CompareInt64                          = sys.CompareInt64
	CompareString                         = sys.CompareString
	ConfigResp                            = sys.ConfigResp
	DataArea                              = sys.DataArea
	DataAreaIndexReq                      = sys.DataAreaIndexReq
	DataAreaIndexResp                     = sys.DataAreaIndexResp
	DataAreaMultiDeleteReq                = sys.DataAreaMultiDeleteReq
	DataAreaMultiUpdateReq                = sys.DataAreaMultiUpdateReq
	DataProject                           = sys.DataProject
	DataProjectDeleteReq                  = sys.DataProjectDeleteReq
	DataProjectIndexReq                   = sys.DataProjectIndexReq
	DataProjectIndexResp                  = sys.DataProjectIndexResp
	DataProjectMultiDeleteReq             = sys.DataProjectMultiDeleteReq
	DataProjectMultiSaveReq               = sys.DataProjectMultiSaveReq
	DataProjectSaveReq                    = sys.DataProjectSaveReq
	DateRange                             = sys.DateRange
	DeptInfo                              = sys.DeptInfo
	DeptInfoIndexReq                      = sys.DeptInfoIndexReq
	DeptInfoIndexResp                     = sys.DeptInfoIndexResp
	DeptInfoReadReq                       = sys.DeptInfoReadReq
	DeptRoleIndexReq                      = sys.DeptRoleIndexReq
	DeptRoleIndexResp                     = sys.DeptRoleIndexResp
	DeptRoleMultiSaveReq                  = sys.DeptRoleMultiSaveReq
	DeptSyncJob                           = sys.DeptSyncJob
	DeptSyncJobExecuteReq                 = sys.DeptSyncJobExecuteReq
	DeptSyncJobExecuteResp                = sys.DeptSyncJobExecuteResp
	DeptSyncJobIndexReq                   = sys.DeptSyncJobIndexReq
	DeptSyncJobIndexResp                  = sys.DeptSyncJobIndexResp
	DeptSyncJobReadReq                    = sys.DeptSyncJobReadReq
	DeptUser                              = sys.DeptUser
	DeptUserIndexReq                      = sys.DeptUserIndexReq
	DeptUserIndexResp                     = sys.DeptUserIndexResp
	DeptUserMultiCreateReq                = sys.DeptUserMultiCreateReq
	DeptUserMultiDeleteReq                = sys.DeptUserMultiDeleteReq
	DictDetail                            = sys.DictDetail
	DictDetailIndexReq                    = sys.DictDetailIndexReq
	DictDetailIndexResp                   = sys.DictDetailIndexResp
	DictDetailMultiCreateReq              = sys.DictDetailMultiCreateReq
	DictDetailReadReq                     = sys.DictDetailReadReq
	DictInfo                              = sys.DictInfo
	DictInfoIndexReq                      = sys.DictInfoIndexReq
	DictInfoIndexResp                     = sys.DictInfoIndexResp
	DictInfoReadReq                       = sys.DictInfoReadReq
	Empty                                 = sys.Empty
	IDList                                = sys.IDList
	JwtToken                              = sys.JwtToken
	LoginLogCreateReq                     = sys.LoginLogCreateReq
	LoginLogIndexReq                      = sys.LoginLogIndexReq
	LoginLogIndexResp                     = sys.LoginLogIndexResp
	LoginLogInfo                          = sys.LoginLogInfo
	Map                                   = sys.Map
	MenuInfo                              = sys.MenuInfo
	MenuInfoIndexReq                      = sys.MenuInfoIndexReq
	MenuInfoIndexResp                     = sys.MenuInfoIndexResp
	MenuMultiExportReq                    = sys.MenuMultiExportReq
	MenuMultiExportResp                   = sys.MenuMultiExportResp
	MenuMultiImportReq                    = sys.MenuMultiImportReq
	MenuMultiImportResp                   = sys.MenuMultiImportResp
	MessageInfo                           = sys.MessageInfo
	MessageInfoIndexReq                   = sys.MessageInfoIndexReq
	MessageInfoIndexResp                  = sys.MessageInfoIndexResp
	MessageInfoSendReq                    = sys.MessageInfoSendReq
	ModuleInfo                            = sys.ModuleInfo
	ModuleInfoIndexReq                    = sys.ModuleInfoIndexReq
	ModuleInfoIndexResp                   = sys.ModuleInfoIndexResp
	NotifyChannel                         = sys.NotifyChannel
	NotifyChannelIndexReq                 = sys.NotifyChannelIndexReq
	NotifyChannelIndexResp                = sys.NotifyChannelIndexResp
	NotifyConfig                          = sys.NotifyConfig
	NotifyConfigIndexReq                  = sys.NotifyConfigIndexReq
	NotifyConfigIndexResp                 = sys.NotifyConfigIndexResp
	NotifyConfigSendReq                   = sys.NotifyConfigSendReq
	NotifyConfigTemplate                  = sys.NotifyConfigTemplate
	NotifyConfigTemplateDeleteReq         = sys.NotifyConfigTemplateDeleteReq
	NotifyConfigTemplateIndexReq          = sys.NotifyConfigTemplateIndexReq
	NotifyConfigTemplateIndexResp         = sys.NotifyConfigTemplateIndexResp
	NotifyTemplate                        = sys.NotifyTemplate
	NotifyTemplateIndexReq                = sys.NotifyTemplateIndexReq
	NotifyTemplateIndexResp               = sys.NotifyTemplateIndexResp
	OpenAccess                            = sys.OpenAccess
	OpenAccessIndexReq                    = sys.OpenAccessIndexReq
	OpenAccessIndexResp                   = sys.OpenAccessIndexResp
	OperLogCreateReq                      = sys.OperLogCreateReq
	OperLogIndexReq                       = sys.OperLogIndexReq
	OperLogIndexResp                      = sys.OperLogIndexResp
	OperLogInfo                           = sys.OperLogInfo
	OpsFeedback                           = sys.OpsFeedback
	OpsFeedbackIndexReq                   = sys.OpsFeedbackIndexReq
	OpsFeedbackIndexResp                  = sys.OpsFeedbackIndexResp
	OpsWorkOrder                          = sys.OpsWorkOrder
	OpsWorkOrderIndexReq                  = sys.OpsWorkOrderIndexReq
	OpsWorkOrderIndexResp                 = sys.OpsWorkOrderIndexResp
	PageInfo                              = sys.PageInfo
	PageInfo_OrderBy                      = sys.PageInfo_OrderBy
	Point                                 = sys.Point
	ProjectAuth                           = sys.ProjectAuth
	ProjectInfo                           = sys.ProjectInfo
	ProjectInfoIndexReq                   = sys.ProjectInfoIndexReq
	ProjectInfoIndexResp                  = sys.ProjectInfoIndexResp
	ProjectProfile                        = sys.ProjectProfile
	ProjectProfileIndexReq                = sys.ProjectProfileIndexReq
	ProjectProfileIndexResp               = sys.ProjectProfileIndexResp
	ProjectProfileReadReq                 = sys.ProjectProfileReadReq
	ProjectWithID                         = sys.ProjectWithID
	QRCodeReadReq                         = sys.QRCodeReadReq
	QRCodeReadResp                        = sys.QRCodeReadResp
	RoleAccessIndexReq                    = sys.RoleAccessIndexReq
	RoleAccessIndexResp                   = sys.RoleAccessIndexResp
	RoleAccessMultiUpdateReq              = sys.RoleAccessMultiUpdateReq
	RoleApiAuthReq                        = sys.RoleApiAuthReq
	RoleApiAuthResp                       = sys.RoleApiAuthResp
	RoleAppIndexReq                       = sys.RoleAppIndexReq
	RoleAppIndexResp                      = sys.RoleAppIndexResp
	RoleAppMultiUpdateReq                 = sys.RoleAppMultiUpdateReq
	RoleAppUpdateReq                      = sys.RoleAppUpdateReq
	RoleInfo                              = sys.RoleInfo
	RoleInfoIndexReq                      = sys.RoleInfoIndexReq
	RoleInfoIndexResp                     = sys.RoleInfoIndexResp
	RoleMenuIndexReq                      = sys.RoleMenuIndexReq
	RoleMenuIndexResp                     = sys.RoleMenuIndexResp
	RoleMenuMultiUpdateReq                = sys.RoleMenuMultiUpdateReq
	RoleModuleIndexReq                    = sys.RoleModuleIndexReq
	RoleModuleIndexResp                   = sys.RoleModuleIndexResp
	RoleModuleMultiUpdateReq              = sys.RoleModuleMultiUpdateReq
	SendOption                            = sys.SendOption
	ServiceInfo                           = sys.ServiceInfo
	SlotInfo                              = sys.SlotInfo
	SlotInfoIndexReq                      = sys.SlotInfoIndexReq
	SlotInfoIndexResp                     = sys.SlotInfoIndexResp
	TenantAccessIndexReq                  = sys.TenantAccessIndexReq
	TenantAccessIndexResp                 = sys.TenantAccessIndexResp
	TenantAccessMultiSaveReq              = sys.TenantAccessMultiSaveReq
	TenantAgreement                       = sys.TenantAgreement
	TenantAgreementIndexReq               = sys.TenantAgreementIndexReq
	TenantAgreementIndexResp              = sys.TenantAgreementIndexResp
	TenantAppIndexReq                     = sys.TenantAppIndexReq
	TenantAppIndexResp                    = sys.TenantAppIndexResp
	TenantAppInfo                         = sys.TenantAppInfo
	TenantAppMenu                         = sys.TenantAppMenu
	TenantAppMenuIndexReq                 = sys.TenantAppMenuIndexReq
	TenantAppMenuIndexResp                = sys.TenantAppMenuIndexResp
	TenantAppModule                       = sys.TenantAppModule
	TenantAppMultiUpdateReq               = sys.TenantAppMultiUpdateReq
	TenantAppWithIDOrCode                 = sys.TenantAppWithIDOrCode
	TenantConfig                          = sys.TenantConfig
	TenantConfigRegisterAutoCreateArea    = sys.TenantConfigRegisterAutoCreateArea
	TenantConfigRegisterAutoCreateProject = sys.TenantConfigRegisterAutoCreateProject
	TenantInfo                            = sys.TenantInfo
	TenantInfoCreateReq                   = sys.TenantInfoCreateReq
	TenantInfoIndexReq                    = sys.TenantInfoIndexReq
	TenantInfoIndexResp                   = sys.TenantInfoIndexResp
	TenantModuleCreateReq                 = sys.TenantModuleCreateReq
	TenantModuleIndexReq                  = sys.TenantModuleIndexReq
	TenantModuleIndexResp                 = sys.TenantModuleIndexResp
	TenantModuleWithIDOrCode              = sys.TenantModuleWithIDOrCode
	TenantOpenCheckTokenReq               = sys.TenantOpenCheckTokenReq
	TenantOpenCheckTokenResp              = sys.TenantOpenCheckTokenResp
	TenantOpenWebHook                     = sys.TenantOpenWebHook
	ThirdApp                              = sys.ThirdApp
	ThirdAppConfig                        = sys.ThirdAppConfig
	ThirdDeptInfoIndexReq                 = sys.ThirdDeptInfoIndexReq
	ThirdDeptInfoReadReq                  = sys.ThirdDeptInfoReadReq
	ThirdEmail                            = sys.ThirdEmail
	ThirdSms                              = sys.ThirdSms
	TimeRange                             = sys.TimeRange
	UserAreaApplyCreateReq                = sys.UserAreaApplyCreateReq
	UserAreaApplyDealReq                  = sys.UserAreaApplyDealReq
	UserAreaApplyIndexReq                 = sys.UserAreaApplyIndexReq
	UserAreaApplyIndexResp                = sys.UserAreaApplyIndexResp
	UserAreaApplyInfo                     = sys.UserAreaApplyInfo
	UserBindAccountReq                    = sys.UserBindAccountReq
	UserCaptchaReq                        = sys.UserCaptchaReq
	UserCaptchaResp                       = sys.UserCaptchaResp
	UserChangePwdReq                      = sys.UserChangePwdReq
	UserCheckTokenReq                     = sys.UserCheckTokenReq
	UserCheckTokenResp                    = sys.UserCheckTokenResp
	UserCodeToUserIDReq                   = sys.UserCodeToUserIDReq
	UserCodeToUserIDResp                  = sys.UserCodeToUserIDResp
	UserCreateResp                        = sys.UserCreateResp
	UserDeptIndexReq                      = sys.UserDeptIndexReq
	UserDeptIndexResp                     = sys.UserDeptIndexResp
	UserDeptMultiSaveReq                  = sys.UserDeptMultiSaveReq
	UserForgetPwdReq                      = sys.UserForgetPwdReq
	UserInfo                              = sys.UserInfo
	UserInfoCreateReq                     = sys.UserInfoCreateReq
	UserInfoDeleteReq                     = sys.UserInfoDeleteReq
	UserInfoIndexReq                      = sys.UserInfoIndexReq
	UserInfoIndexResp                     = sys.UserInfoIndexResp
	UserInfoReadReq                       = sys.UserInfoReadReq
	UserInfoUpdateReq                     = sys.UserInfoUpdateReq
	UserLoginReq                          = sys.UserLoginReq
	UserLoginResp                         = sys.UserLoginResp
	UserMessage                           = sys.UserMessage
	UserMessageIndexReq                   = sys.UserMessageIndexReq
	UserMessageIndexResp                  = sys.UserMessageIndexResp
	UserMessageStatistics                 = sys.UserMessageStatistics
	UserMessageStatisticsResp             = sys.UserMessageStatisticsResp
	UserProfile                           = sys.UserProfile
	UserProfileIndexReq                   = sys.UserProfileIndexReq
	UserProfileIndexResp                  = sys.UserProfileIndexResp
	UserRegisterReq                       = sys.UserRegisterReq
	UserRegisterResp                      = sys.UserRegisterResp
	UserRoleIndexReq                      = sys.UserRoleIndexReq
	UserRoleIndexResp                     = sys.UserRoleIndexResp
	UserRoleMultiUpdateReq                = sys.UserRoleMultiUpdateReq
	WeatherAir                            = sys.WeatherAir
	WeatherReadReq                        = sys.WeatherReadReq
	WeatherReadResp                       = sys.WeatherReadResp
	WithAppCodeID                         = sys.WithAppCodeID
	WithCode                              = sys.WithCode
	WithID                                = sys.WithID
	WithIDCode                            = sys.WithIDCode

	DepartmentManage interface {
		DeptInfoRead(ctx context.Context, in *DeptInfoReadReq, opts ...grpc.CallOption) (*DeptInfo, error)
		DeptInfoCreate(ctx context.Context, in *DeptInfo, opts ...grpc.CallOption) (*WithID, error)
		DeptInfoIndex(ctx context.Context, in *DeptInfoIndexReq, opts ...grpc.CallOption) (*DeptInfoIndexResp, error)
		DeptInfoUpdate(ctx context.Context, in *DeptInfo, opts ...grpc.CallOption) (*Empty, error)
		DeptInfoDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error)
		DeptUserIndex(ctx context.Context, in *DeptUserIndexReq, opts ...grpc.CallOption) (*DeptUserIndexResp, error)
		DeptUserMultiDelete(ctx context.Context, in *DeptUserMultiDeleteReq, opts ...grpc.CallOption) (*Empty, error)
		DeptUserMultiCreate(ctx context.Context, in *DeptUserMultiCreateReq, opts ...grpc.CallOption) (*Empty, error)
		DeptRoleIndex(ctx context.Context, in *DeptRoleIndexReq, opts ...grpc.CallOption) (*DeptRoleIndexResp, error)
		DeptRoleMultiDelete(ctx context.Context, in *DeptRoleMultiSaveReq, opts ...grpc.CallOption) (*Empty, error)
		DeptRoleMultiCreate(ctx context.Context, in *DeptRoleMultiSaveReq, opts ...grpc.CallOption) (*Empty, error)
		DeptSyncJobExecute(ctx context.Context, in *DeptSyncJobExecuteReq, opts ...grpc.CallOption) (*DeptSyncJobExecuteResp, error)
		DeptSyncJobRead(ctx context.Context, in *DeptSyncJobReadReq, opts ...grpc.CallOption) (*DeptSyncJob, error)
		DeptSyncJobCreate(ctx context.Context, in *DeptSyncJob, opts ...grpc.CallOption) (*WithID, error)
		DeptSyncJobIndex(ctx context.Context, in *DeptSyncJobIndexReq, opts ...grpc.CallOption) (*DeptSyncJobIndexResp, error)
		DeptSyncJobUpdate(ctx context.Context, in *DeptSyncJob, opts ...grpc.CallOption) (*Empty, error)
		DeptSyncJobDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error)
	}

	defaultDepartmentManage struct {
		cli zrpc.Client
	}

	directDepartmentManage struct {
		svcCtx *svc.ServiceContext
		svr    sys.DepartmentManageServer
	}
)

func NewDepartmentManage(cli zrpc.Client) DepartmentManage {
	return &defaultDepartmentManage{
		cli: cli,
	}
}

func NewDirectDepartmentManage(svcCtx *svc.ServiceContext, svr sys.DepartmentManageServer) DepartmentManage {
	return &directDepartmentManage{
		svr:    svr,
		svcCtx: svcCtx,
	}
}

func (m *defaultDepartmentManage) DeptInfoRead(ctx context.Context, in *DeptInfoReadReq, opts ...grpc.CallOption) (*DeptInfo, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptInfoRead(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptInfoRead(ctx context.Context, in *DeptInfoReadReq, opts ...grpc.CallOption) (*DeptInfo, error) {
	return d.svr.DeptInfoRead(ctx, in)
}

func (m *defaultDepartmentManage) DeptInfoCreate(ctx context.Context, in *DeptInfo, opts ...grpc.CallOption) (*WithID, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptInfoCreate(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptInfoCreate(ctx context.Context, in *DeptInfo, opts ...grpc.CallOption) (*WithID, error) {
	return d.svr.DeptInfoCreate(ctx, in)
}

func (m *defaultDepartmentManage) DeptInfoIndex(ctx context.Context, in *DeptInfoIndexReq, opts ...grpc.CallOption) (*DeptInfoIndexResp, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptInfoIndex(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptInfoIndex(ctx context.Context, in *DeptInfoIndexReq, opts ...grpc.CallOption) (*DeptInfoIndexResp, error) {
	return d.svr.DeptInfoIndex(ctx, in)
}

func (m *defaultDepartmentManage) DeptInfoUpdate(ctx context.Context, in *DeptInfo, opts ...grpc.CallOption) (*Empty, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptInfoUpdate(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptInfoUpdate(ctx context.Context, in *DeptInfo, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeptInfoUpdate(ctx, in)
}

func (m *defaultDepartmentManage) DeptInfoDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptInfoDelete(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptInfoDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeptInfoDelete(ctx, in)
}

func (m *defaultDepartmentManage) DeptUserIndex(ctx context.Context, in *DeptUserIndexReq, opts ...grpc.CallOption) (*DeptUserIndexResp, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptUserIndex(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptUserIndex(ctx context.Context, in *DeptUserIndexReq, opts ...grpc.CallOption) (*DeptUserIndexResp, error) {
	return d.svr.DeptUserIndex(ctx, in)
}

func (m *defaultDepartmentManage) DeptUserMultiDelete(ctx context.Context, in *DeptUserMultiDeleteReq, opts ...grpc.CallOption) (*Empty, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptUserMultiDelete(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptUserMultiDelete(ctx context.Context, in *DeptUserMultiDeleteReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeptUserMultiDelete(ctx, in)
}

func (m *defaultDepartmentManage) DeptUserMultiCreate(ctx context.Context, in *DeptUserMultiCreateReq, opts ...grpc.CallOption) (*Empty, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptUserMultiCreate(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptUserMultiCreate(ctx context.Context, in *DeptUserMultiCreateReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeptUserMultiCreate(ctx, in)
}

func (m *defaultDepartmentManage) DeptRoleIndex(ctx context.Context, in *DeptRoleIndexReq, opts ...grpc.CallOption) (*DeptRoleIndexResp, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptRoleIndex(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptRoleIndex(ctx context.Context, in *DeptRoleIndexReq, opts ...grpc.CallOption) (*DeptRoleIndexResp, error) {
	return d.svr.DeptRoleIndex(ctx, in)
}

func (m *defaultDepartmentManage) DeptRoleMultiDelete(ctx context.Context, in *DeptRoleMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptRoleMultiDelete(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptRoleMultiDelete(ctx context.Context, in *DeptRoleMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeptRoleMultiDelete(ctx, in)
}

func (m *defaultDepartmentManage) DeptRoleMultiCreate(ctx context.Context, in *DeptRoleMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptRoleMultiCreate(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptRoleMultiCreate(ctx context.Context, in *DeptRoleMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeptRoleMultiCreate(ctx, in)
}

func (m *defaultDepartmentManage) DeptSyncJobExecute(ctx context.Context, in *DeptSyncJobExecuteReq, opts ...grpc.CallOption) (*DeptSyncJobExecuteResp, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptSyncJobExecute(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptSyncJobExecute(ctx context.Context, in *DeptSyncJobExecuteReq, opts ...grpc.CallOption) (*DeptSyncJobExecuteResp, error) {
	return d.svr.DeptSyncJobExecute(ctx, in)
}

func (m *defaultDepartmentManage) DeptSyncJobRead(ctx context.Context, in *DeptSyncJobReadReq, opts ...grpc.CallOption) (*DeptSyncJob, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptSyncJobRead(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptSyncJobRead(ctx context.Context, in *DeptSyncJobReadReq, opts ...grpc.CallOption) (*DeptSyncJob, error) {
	return d.svr.DeptSyncJobRead(ctx, in)
}

func (m *defaultDepartmentManage) DeptSyncJobCreate(ctx context.Context, in *DeptSyncJob, opts ...grpc.CallOption) (*WithID, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptSyncJobCreate(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptSyncJobCreate(ctx context.Context, in *DeptSyncJob, opts ...grpc.CallOption) (*WithID, error) {
	return d.svr.DeptSyncJobCreate(ctx, in)
}

func (m *defaultDepartmentManage) DeptSyncJobIndex(ctx context.Context, in *DeptSyncJobIndexReq, opts ...grpc.CallOption) (*DeptSyncJobIndexResp, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptSyncJobIndex(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptSyncJobIndex(ctx context.Context, in *DeptSyncJobIndexReq, opts ...grpc.CallOption) (*DeptSyncJobIndexResp, error) {
	return d.svr.DeptSyncJobIndex(ctx, in)
}

func (m *defaultDepartmentManage) DeptSyncJobUpdate(ctx context.Context, in *DeptSyncJob, opts ...grpc.CallOption) (*Empty, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptSyncJobUpdate(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptSyncJobUpdate(ctx context.Context, in *DeptSyncJob, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeptSyncJobUpdate(ctx, in)
}

func (m *defaultDepartmentManage) DeptSyncJobDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	client := sys.NewDepartmentManageClient(m.cli.Conn())
	return client.DeptSyncJobDelete(ctx, in, opts...)
}

func (d *directDepartmentManage) DeptSyncJobDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeptSyncJobDelete(ctx, in)
}
