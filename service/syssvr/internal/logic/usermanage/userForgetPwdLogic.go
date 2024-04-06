package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserForgetPwdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserForgetPwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserForgetPwdLogic {
	return &UserForgetPwdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserForgetPwdLogic) UserForgetPwd(in *sys.UserForgetPwdReq) (*sys.Empty, error) {
	var account string
	var oldUi *relationDB.SysUserInfo
	switch in.Type {
	case def.CaptchaTypeEmail:
		account = l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, def.CaptchaUseForgetPwd, in.CodeID, in.Code)
		if account == "" {
			return nil, errors.Captcha
		}
		ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Emails: []string{account}})
		if err != nil {
			return nil, err
		}
		oldUi = ui
	case def.CaptchaTypePhone:
		account = l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypePhone, def.CaptchaUseForgetPwd, in.CodeID, in.Code)
		if account == "" {
			return nil, errors.Captcha
		}
		ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phones: []string{account}})
		if err != nil {
			return nil, err
		}
		oldUi = ui
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
	return &sys.Empty{}, nil
}
