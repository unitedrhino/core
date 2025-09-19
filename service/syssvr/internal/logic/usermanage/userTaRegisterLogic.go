package usermanagelogic

import (
	"context"

	tenantmanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/tenantmanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/core/share/users"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserTaRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	ur     *UserRegisterLogic
	logx.Logger
}

func NewUserTaRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserTaRegisterLogic {
	return &UserTaRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ur:     NewUserRegisterLogic(ctx, svcCtx),
	}
}

func (l *UserTaRegisterLogic) UserTaRegister(in *sys.UserTaRegisterReq) (*sys.UserRegisterResp, error) {
	if !l.svcCtx.Config.TenantConfig.IsEnableTaRegister {
		return &sys.UserRegisterResp{}, errors.Permissions.AddMsg("未开启租户管理员注册")
	}
	err := tenantmanagelogic.CheckAdmin(l.ctx, l.svcCtx, in.UserName, in.NickName, in.Password, in.Account)
	if err != nil {
		return nil, err
	}
	userID := l.svcCtx.UserID.GetSnowflakeId()

	ui := l.ur.createUserInfo(userID, in.NickName, in.UserName)
	ui.Password = utils.MakePwd(in.Password, ui.UserID, false)
	// 验证验证码并设置账户信息
	switch in.RegType {
	case users.RegEmail:
		email, err := l.ur.verifyCaptcha(def.CaptchaTypeEmail, in.CodeID, in.Code, in.Account)
		if err != nil {
			return nil, err
		}
		ui.Email = utils.AnyToNullString(email)
	case users.RegPhone:
		_, err := l.ur.verifyCaptcha(def.CaptchaTypePhone, in.CodeID, in.Code, in.Account)
		if err != nil {
			return nil, err
		}
		ui.Phone = utils.AnyToNullString(in.Account)
	default:
		return nil, errors.Parameter.AddMsg("只支持手机号或邮箱的校验方式")
	}

	projectPo := &relationDB.SysProjectInfo{
		TenantCode:   dataType.TenantCode(in.Code),
		ProjectID:    dataType.ProjectID(l.svcCtx.ProjectID.GetSnowflakeId()),
		ProjectName:  "默认项目",
		AdminUserID:  ui.UserID,
		IsSysCreated: def.True,
	}

	po := &relationDB.SysTenantInfo{
		Code:             dataType.TenantCode(ui.UserName.String),
		Name:             ui.UserName.String,
		AdminUserID:      ui.UserID,
		DefaultProjectID: int64(projectPo.ProjectID),
		UserCount:        1,
	}
	ui.IsTenantAdmin = def.True
	err = tenantmanagelogic.TenantCreate(l.ctx, l.svcCtx, projectPo, po, &ui)
	return &sys.UserRegisterResp{UserID: userID}, err
}
