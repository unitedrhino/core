package usermanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

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

func (l *UserChangePwdLogic) UserChangePwd(in *sys.UserChangePwdReq) (*sys.Response, error) {
	var account string
	uc := ctxs.GetUserCtx(l.ctx)
	var oldUi *relationDB.SysUserInfo
	switch in.Type {
	case def.CaptchaTypeEmail:
		account = l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, def.CaptchaUseChangePwd, in.CodeID, in.Code)
		if account == "" {
			return nil, errors.Captcha
		}
		ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Emails: []string{account}})
		if err != nil {
			return nil, err
		}
		oldUi = ui
	case def.CaptchaTypePhone:
		account = l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypePhone, def.CaptchaUseChangePwd, in.CodeID, in.Code)
		if account == "" {
			return nil, errors.Captcha
		}
		ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phones: []string{account}})
		if err != nil {
			return nil, err
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
	return &sys.Response{}, nil
}
