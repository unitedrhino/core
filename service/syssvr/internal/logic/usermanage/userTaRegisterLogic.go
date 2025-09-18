package usermanagelogic

import (
	"context"

	tenantmanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/tenantmanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/core/share/users"
	"gitee.com/unitedrhino/share/ctxs"
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
	if in.Password == "" {
		return nil, errors.Parameter.AddMsg("密码必填")
	}
	if in.UserName == "" {
		return nil, errors.Parameter.AddMsg("用户名必填")
	}
	if err := l.CheckAccount(in); err != nil {
		return nil, err
	}
	// 验证密码
	if err := l.ur.validatePassword(in.Password); err != nil {
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
		//phone, err := l.ur.verifyCaptcha(def.CaptchaTypePhone, in.CodeID, in.Code, in.Account)
		//if err != nil {
		//	return nil, err
		//}
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
	err := tenantmanagelogic.TenantCreate(l.ctx, l.svcCtx, projectPo, po, &ui, true)
	return &sys.UserRegisterResp{UserID: userID}, err
}

func (l *UserTaRegisterLogic) CheckAccount(in *sys.UserTaRegisterReq) error {
	if in.UserName == "" {
		return errors.Parameter.AddMsg("用户名必填")
	}
	err := CheckUserName(in.UserName)
	if err != nil {
		return err
	}
	_, err = relationDB.NewTenantInfoRepo(l.ctx).FindOneByFilter(ctxs.WithRoot(l.ctx), relationDB.TenantInfoFilter{
		Code: in.UserName,
	})
	if err == nil {
		return errors.Parameter.AddMsg("用户名已被使用")
	}
	if !errors.Cmp(err, errors.NotFind) {
		return err
	}

	po, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(ctxs.WithCommonTenant(l.ctx), relationDB.UserInfoFilter{
		Accounts: []string{in.UserName, in.Account},
	})
	if err == nil {
		if po.UserName.Valid && po.UserName.String == in.UserName {
			return errors.Parameter.AddMsg("用户名已被使用")
		}
		if po.Email.Valid && po.Email.String == in.Account {
			return errors.Parameter.AddMsg("邮箱已被使用")
		}
		if po.Phone.Valid && po.Phone.String == in.Account {
			return errors.Parameter.AddMsg("手机号已被使用")
		}
	}
	if !errors.Cmp(err, errors.NotFind) {
		return err
	}

	return nil
}
