package usermanagelogic

import (
	"context"
	notifymanagelogic "gitee.com/i-Things/core/service/syssvr/internal/logic/notifymanage"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
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
	//todo 防盗刷也需要做下处理
	//code = "123456" //todo debug
	if utils.SliceIn(in.Type, def.CaptchaTypePhone, def.CaptchaTypeEmail) && in.Account == "" {
		return nil, errors.Parameter.AddMsg("account需要填写")
	}

	uc := ctxs.GetUserCtx(l.ctx)
	switch in.Type {
	case def.CaptchaTypeImage:
	case def.CaptchaTypePhone:
		var imgAuth bool
		if in.Code != "" {
			account := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeImage, in.Use, in.CodeID, in.Code)
			if account == "" {
				return nil, errors.Captcha
			}
			imgAuth = true
		}
		switch in.Use {
		case def.CaptchaUseLogin, def.CaptchaUseRegister, def.CaptchaUseForgetPwd:
			if imgAuth == false {
				return nil, errors.Captcha.AddMsg("需要先填写图形验证码")
			}
		}
		var ConfigCode = def.NotifyCodeSysUserRegisterCaptcha
		if !utils.SliceIn(in.Use, def.CaptchaUseRegister) {
			ConfigCode = def.NotifyCodeSysUserLoginCaptcha
		}
		err := notifymanagelogic.SendNotifyMsg(l.ctx, l.svcCtx, notifymanagelogic.SendMsgConfig{
			Accounts:    []string{in.Account},
			AccountType: def.AccountTypePhone,
			NotifyCode:  ConfigCode,
			Type:        def.NotifyTypeSms,
			Params:      map[string]any{"code": code, "expr": def.CaptchaExpire / 60},
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
		var ConfigCode = def.NotifyCodeSysUserRegisterCaptcha
		if !utils.SliceIn(in.Use, def.CaptchaUseRegister) {
			ConfigCode = def.NotifyCodeSysUserLoginCaptcha
		}
		err := notifymanagelogic.SendNotifyMsg(l.ctx, l.svcCtx, notifymanagelogic.SendMsgConfig{
			Accounts:    []string{in.Account},
			AccountType: def.AccountTypeEmail,
			NotifyCode:  ConfigCode,
			Type:        def.NotifyTypeEmail,
			Params:      map[string]any{"code": code, "expr": def.CaptchaExpire / 60},
		})
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
