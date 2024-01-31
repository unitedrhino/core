package usermanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/ctxs"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCaptchaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCaptchaLogic {
	return &UserCaptchaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserCaptchaLogic) UserCaptcha(in *sys.UserCaptchaReq) (*sys.UserCaptchaResp, error) {
	var (
		codeID = utils.Random(20, 1)
		code   = utils.Random(6, 0)
	)
	if utils.SliceIn(in.Type, def.CaptchaTypePhone, def.CaptchaTypeEmail) && in.Account == "" {
		return nil, errors.Parameter.AddMsg("account需要填写")
	}
	uc := ctxs.GetUserCtx(l.ctx)
	switch in.Type {
	case def.CaptchaTypeImage:
	case def.CaptchaTypePhone:
		if uc == nil {
			account := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypePhone, in.Use, in.CodeID, in.Code)
			if account == "" {
				return nil, errors.Captcha
			}
		}
		if !utils.SliceIn(in.Use, def.CaptchaUseRegister) {
			count, err := relationDB.NewUserInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.UserInfoFilter{Phones: []string{in.Account}})
			if err != nil {
				return nil, err
			}
			if count == 0 && in.Use == def.CaptchaUseLogin {
				return nil, errors.UnRegister
			}
		}
		code = "123456"
	case def.CaptchaTypeEmail:
		if uc == nil {
			account := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, in.Use, in.CodeID, in.Code)
			if account == "" {
				return nil, errors.Captcha
			}
		}
		if !utils.SliceIn(in.Use, def.CaptchaUseRegister) {
			count, err := relationDB.NewUserInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.UserInfoFilter{Emails: []string{in.Account}})
			if err != nil {
				return nil, err
			}
			if count == 0 && in.Use == def.CaptchaUseLogin {
				return nil, errors.UnRegister
			}
		}
		c, err := relationDB.NewTenantConfigRepo(l.ctx).FindOne(l.ctx)
		if err != nil {
			return nil, err
		}
		code = "123456"
		err = utils.SenEmail(conf.Email{
			From:     c.Email.From,
			Host:     c.Email.Host,
			Secret:   c.Email.Secret,
			Nickname: c.Email.Nickname,
			Port:     c.Email.Port,
			IsSSL:    c.Email.IsSSL == def.True,
		}, []string{in.Account}, "验证码校验",
			fmt.Sprintf("您的验证码为：%s，有效期为%d分钟", code, def.CaptchaExpire/60))
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.Parameter.AddMsgf("不支持的验证方式:%v", in.Type)
	}
	err := l.svcCtx.Captcha.Store(l.ctx, in.Type, in.Use, codeID, code, in.Account, def.CaptchaExpire)
	if err != nil {
		return nil, err
	}
	l.Infof("code:%v codeID:%v", code, codeID)
	return &sys.UserCaptchaResp{Code: code, CodeID: codeID, Expire: def.CaptchaExpire}, nil
}
