package usermanagelogic

import (
	"bytes"
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"text/template"

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
		var ConfigCode = def.CaptchaUseRegister
		if !utils.SliceIn(in.Use, def.CaptchaUseRegister) {
			count, err := relationDB.NewUserInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.UserInfoFilter{Phones: []string{in.Account}})
			if err != nil {
				return nil, err
			}
			if count == 0 && in.Use == def.CaptchaUseLogin {
				return nil, errors.UnRegister
			}
			ConfigCode = def.NotifyCodeSysUserLoginCaptcha
		}
		c, err := relationDB.NewTenantNotifyTemplateRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantNotifyConfigFilter{
			ConfigCode: ConfigCode,
			Type:       def.NotifyTypeSms,
		})
		if err != nil {
			if errors.Cmp(err, errors.NotFind) {
				return nil, errors.NotEnable
			}
			return nil, err
		}
		err = l.svcCtx.Sms.SendSms(clients.SendSmsParam{
			PhoneNumbers:  in.Account,
			SignName:      c.Template.SignName,
			TemplateCode:  c.Template.Code,
			TemplateParam: map[string]any{"code": code},
		})
		if err != nil {
			return nil, err
		}
	case def.CaptchaTypeEmail:
		if uc == nil {
			account := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, in.Use, in.CodeID, in.Code)
			if account == "" {
				return nil, errors.Captcha
			}
		}
		var ConfigCode = def.CaptchaUseRegister
		if !utils.SliceIn(in.Use, def.CaptchaUseRegister) {
			count, err := relationDB.NewUserInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.UserInfoFilter{Emails: []string{in.Account}})
			if err != nil {
				return nil, err
			}
			if count == 0 && in.Use == def.CaptchaUseLogin {
				return nil, errors.UnRegister
			}
			ConfigCode = def.NotifyCodeSysUserLoginCaptcha
		}
		c, err := relationDB.NewTenantNotifyTemplateRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantNotifyConfigFilter{
			ConfigCode: ConfigCode,
			Type:       def.NotifyTypeEmail,
		})
		if err != nil {
			if errors.Cmp(err, errors.NotFind) {
				return nil, errors.NotEnable
			}
			return nil, err
		}
		tc, err := relationDB.NewTenantConfigRepo(l.ctx).FindOne(l.ctx)
		if err != nil {
			return nil, err
		}
		tmpl, err := template.New(c.Template.Code).Parse(c.Template.Body)
		if err != nil {
			return nil, errors.System.AddMsg("模版解析失败").AddDetail(err)
		}
		buffer := &bytes.Buffer{}
		err = tmpl.Execute(buffer, map[string]any{"code": code, "expr": def.CaptchaExpire / 60})
		if err != nil {
			return nil, errors.System.AddMsg("模版匹配失败").AddDetail(err)
		}
		err = utils.SenEmail(conf.Email{
			From:     tc.Email.From,
			Host:     tc.Email.Host,
			Secret:   tc.Email.Secret,
			Nickname: tc.Email.Nickname,
			Port:     tc.Email.Port,
			IsSSL:    tc.Email.IsSSL == def.True,
		}, []string{in.Account}, c.Template.Name,
			buffer.String())
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
