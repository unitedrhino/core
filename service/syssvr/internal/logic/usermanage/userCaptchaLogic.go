package usermanagelogic

import (
	"context"

	notifymanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/notifymanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
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
		code   = utils.Random(l.svcCtx.Config.CaptchaLen, 0)
	)
	//code = "123456" //todo debug
	if utils.SliceIn(in.Type, def.CaptchaTypePhone, def.CaptchaTypeEmail) && in.Account == "" {
		return nil, errors.Parameter.AddMsg("account需要填写")
	}
	ip := ctxs.GetUserCtxNoNil(l.ctx).IP
	switch in.Type {
	case def.CaptchaTypeImage:
	case def.CaptchaTypePhone:
		if in.Account == "" {
			return nil, errors.Parameter.AddMsg("请输入手机号")
		}
		if in.Code != "" {
			account := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeImage, in.Use, in.CodeID, in.Code)
			if account == "" {
				return nil, errors.Captcha
			}
		} else if l.svcCtx.CaptchaLimit.PhoneGet.CheckLimit(l.ctx, in.Account) {
			return nil, errors.NeedImgCaptcha
		}
		if l.svcCtx.CaptchaLimit.PhoneAccount.CheckLimit(l.ctx, in.Account) {
			return nil, errors.AccountOrIpForbidden.WithMsg("获取过于频繁,请稍后再试").AddDetail("PhoneAccount")
		}
		if ip != "" && l.svcCtx.CaptchaLimit.PhoneIp.CheckLimit(l.ctx, ip) {
			return nil, errors.AccountOrIpForbidden.WithMsg("获取过于频繁,请稍后再试").AddDetail("PhoneIp")
		}
		var ConfigCode = def.NotifyCodeSysUserRegisterCaptcha
		if in.Use == def.CaptchaUseRegister { //注册的时候要检查下是否已经注册了,如果注册了返回错误
			u, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phone: in.Account})
			if err != nil && !errors.Cmp(err, errors.NotFind) {
				return nil, err
			}
			if u != nil {
				return nil, errors.DuplicateRegister
			}
		} else {
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
		l.svcCtx.CaptchaLimit.PhoneAccount.LimitIt(l.ctx, in.Account)
		l.svcCtx.CaptchaLimit.PhoneGet.LimitIt(l.ctx, in.Account)
		if ip != "" {
			l.svcCtx.CaptchaLimit.PhoneIp.LimitIt(l.ctx, ip)
		}
	case def.CaptchaTypeEmail:
		if in.Account == "" {
			return nil, errors.Parameter.AddMsg("请输入邮箱")
		}
		if in.Code != "" {
			account := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeImage, in.Use, in.CodeID, in.Code)
			if account == "" {
				return nil, errors.Captcha
			}
		} else if l.svcCtx.CaptchaLimit.EmailGet.CheckLimit(l.ctx, in.Account) {
			return nil, errors.NeedImgCaptcha
		}
		if l.svcCtx.CaptchaLimit.EmailAccount.CheckLimit(l.ctx, in.Account) {
			return nil, errors.AccountOrIpForbidden.WithMsg("获取过于频繁,请稍后再试").AddDetail("EmailAccount")
		}
		if ip != "" && l.svcCtx.CaptchaLimit.EmailIp.CheckLimit(l.ctx, ip) {
			return nil, errors.AccountOrIpForbidden.WithMsg("获取过于频繁,请稍后再试").AddDetail("CaptchaLimit")
		}
		var ConfigCode = def.NotifyCodeSysUserRegisterCaptcha
		if in.Use == def.CaptchaUseRegister { //注册的时候要检查下是否已经注册了,如果注册了返回错误
			u, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Email: in.Account})
			if err != nil && !errors.Cmp(err, errors.NotFind) {
				return nil, err
			}
			if u != nil {
				return nil, errors.DuplicateRegister
			}
		} else {
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
		l.svcCtx.CaptchaLimit.EmailGet.LimitIt(l.ctx, in.Account)
		l.svcCtx.CaptchaLimit.EmailAccount.LimitIt(l.ctx, in.Account)
		if ip != "" {
			l.svcCtx.CaptchaLimit.EmailIp.LimitIt(l.ctx, ip)
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
