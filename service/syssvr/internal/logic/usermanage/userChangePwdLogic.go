package usermanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/share/users"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserChangePwdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserChangePwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserChangePwdLogic {
	return &UserChangePwdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserChangePwdLogic) UserChangePwd(in *sys.UserChangePwdReq) (*sys.Empty, error) {
	var account string
	uc := ctxs.GetUserCtx(l.ctx)
	var oldUi *relationDB.SysUserInfo
	switch in.Type {
	case def.CaptchaTypeEmail:
		account = l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, def.CaptchaUseChangePwd, in.CodeID, in.Code)
		if account == "" {
			return nil, errors.Captcha
		}
		ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOne(l.ctx, uc.UserID)
		if err != nil {
			return nil, err
		}
		if ui.Email.String != account {
			return nil, errors.UnBindAccount
		}
		oldUi = ui
	case def.CaptchaTypePhone:
		account = l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypePhone, def.CaptchaUseChangePwd, in.CodeID, in.Code)
		if account == "" {
			return nil, errors.Captcha
		}
		ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOne(l.ctx, uc.UserID)
		if err != nil {
			return nil, err
		}
		if ui.Phone.String != account {
			return nil, errors.UnBindAccount
		}
		oldUi = ui
	case users.RegPwd:
		if in.Code != "" {
			if l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeImage, def.CaptchaUseChangePwd, in.CodeID, in.Code) == "" {
				return nil, errors.Captcha
			}
		} else if l.svcCtx.LoginLimit.PwdCaptcha.CheckLimit(l.ctx, "changePwd:"+uc.Account) {
			return nil, errors.NeedImgCaptcha
		}
		l.svcCtx.LoginLimit.PwdCaptcha.LimitIt(l.ctx, "changePwd:"+uc.Account)
		ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOne(l.ctx, uc.UserID)
		if err != nil {
			return nil, err
		}
		if ui.Password != "" {
			//md5加密后的密码则通过二次md5加密再对比库中的密码
			password1 := utils.MakePwd(in.OldPassword, uc.UserID, true) //对密码进行md5加密
			if password1 != ui.Password {
				return nil, errors.Password
			}
		}
		oldUi = ui
	}
	if oldUi.UserID != uc.UserID {
		return nil, errors.Permissions.AddMsgf("只能修改自己的密码")
	}
	err := CheckPwd(l.svcCtx, in.Password)
	if err != nil {
		return nil, err
	}
	oldUi.Password = utils.MakePwd(in.Password, oldUi.UserID, false)
	err = relationDB.NewUserInfoRepo(l.ctx).Update(l.ctx, oldUi)
	if err != nil {
		return nil, err
	}
	e := l.svcCtx.UserToken.KickedOut(l.ctx, oldUi.UserID)
	if e != nil {
		l.Error(e)
	}
	return &sys.Empty{}, nil
}
