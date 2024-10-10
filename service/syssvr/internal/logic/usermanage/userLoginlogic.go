package usermanagelogic

import (
	"context"
	"database/sql"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/users"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zhaoyunxing92/dingtalk/v2/request"
	"gorm.io/gorm"
	"time"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *LoginLogic) getPwd(in *sys.UserLoginReq, uc *relationDB.SysUserInfo) error {
	//根据密码类型不同做不同处理
	if in.PwdType == 0 {
		//空密码情况暂不考虑
		return errors.UnRegister
	} else if in.PwdType == 1 {
		//明文密码，则对密码做MD5加密后再与数据库密码比对
		//uid_temp := l.svcCtx.UserID.GetSnowflakeId()
		password1 := utils.MakePwd(in.Password, uc.UserID, false) //对密码进行md5加密
		if password1 != uc.Password {
			return errors.Password
		}
	} else if in.PwdType == 2 {
		//md5加密后的密码则通过二次md5加密再对比库中的密码
		password1 := utils.MakePwd(in.Password, uc.UserID, true) //对密码进行md5加密
		if password1 != uc.Password {
			return errors.Password
		}
	} else {
		return errors.Password
	}
	return nil
}

func (l *LoginLogic) getRet(ui *relationDB.SysUserInfo, list []*conf.LoginSafeCtlInfo) (*sys.UserLoginResp, error) {
	now := time.Now()
	accessExpire := l.svcCtx.Config.UserToken.AccessExpire
	uc := ctxs.GetUserCtx(l.ctx)
	jwtToken, err := users.GetLoginJwtToken(l.svcCtx.Config.UserToken.AccessSecret, now, accessExpire,
		ui.UserID, uc.AppCode)
	if err != nil {
		l.Error(err)
		return nil, errors.System.AddDetail(err)
	}

	//登录成功，清除密码错误次数相关redis key
	l.svcCtx.PwdCheck.ClearWrongpassKeys(l.ctx, list)

	//InitCacheUserAuthProject(l.ctx, ui.UserID)
	//InitCacheUserAuthArea(l.ctx, ui.UserID)

	resp := &sys.UserLoginResp{
		Info: UserInfoToPb(l.ctx, ui, l.svcCtx),
		Token: &sys.JwtToken{
			AccessToken:  jwtToken,
			AccessExpire: now.Unix() + accessExpire,
			RefreshAfter: now.Unix() + accessExpire/2,
		},
	}
	l.Infof("%s getRet=%+v", utils.FuncName(), resp)
	return resp, nil
}

func (l *LoginLogic) GetUserInfo(in *sys.UserLoginReq) (uc *relationDB.SysUserInfo, err error) {
	cli, er := l.svcCtx.Cm.GetClients(l.ctx, "")
	if er != nil {
		return nil, errors.System.AddDetail(err)
	}
	if !utils.SliceIn(in.LoginType, cli.Config.LoginTypes...) {
		l.Errorf("不支持的登录方式:%v", in.LoginType)
		//return nil, errors.NotSupportLogin
	}
	switch in.LoginType {
	case users.RegPwd:
		if l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeImage, def.CaptchaUseLogin, in.CodeID, in.Code) == "" {
			return nil, errors.Captcha
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: []string{in.Account}})
		if err != nil {
			return nil, err
		}
		if err = l.getPwd(in, uc); err != nil {
			return nil, err
		}
	case users.RegDingApp:
		if cli.DingMini == nil {
			return nil, errors.System.AddDetail(err)
		}
		ret, er := cli.DingMini.GetUserInfoByCode(in.Code)
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		if ret.Code != 0 {
			return nil, errors.Parameter.AddMsgf(ret.Msg)
		}

		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{DingTalkUserID: ret.UserInfo.UserId})
		if errors.Cmp(err, errors.NotFind) && cli.Config.IsAutoRegister == def.True { //未注册,自动注册
			err = nil
			userID := l.svcCtx.UserID.GetSnowflakeId()
			uc = &relationDB.SysUserInfo{
				UserID:         userID,
				DingTalkUserID: sql.NullString{Valid: true, String: ret.UserInfo.UserId},
				NickName:       ret.UserInfo.Name,
			}
			ui, er := cli.DingMini.GetUserDetail(&request.UserDetail{
				UserId: ret.UserInfo.UserId,
			})
			l.Infof("GetUserDetail ui:%v err:%v", utils.Fmt(ui), er)
			if er == nil {
				if ui.OrgEmail != "" {
					uc.Email = sql.NullString{String: ui.OrgEmail, Valid: true}
					uc.UserName = uc.Email
				}
				if ui.Mobile != "" {
					uc.Phone = sql.NullString{String: ui.Mobile, Valid: true}
					uc.UserName = uc.Phone
				}

			}
			err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
				return Register(l.ctx, l.svcCtx, uc, tx)
			})
			if err != nil {
				return nil, err
			}
			uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{DingTalkUserID: ret.UserInfo.UserId, WithRoles: true, WithTenant: true})
		}
	case users.RegWxOpen:
		if cli.WxOfficial == nil {
			return nil, errors.System.AddDetail(er)
		}
		at, er := cli.WxOfficial.GetOauth().GetUserAccessToken(in.Code)
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{WechatUnionID: at.UnionID, WechatOpenID: at.OpenID})
	case users.RegWxMiniP:
		if cli.MiniProgram == nil {
			return nil, errors.System.AddDetail(er)
		}
		auth := cli.MiniProgram.GetAuth()
		ret, er := auth.Code2SessionContext(l.ctx, in.Code)
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		if ret.ErrCode != 0 {
			return nil, errors.Parameter.AddMsgf(ret.ErrMsg)
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{WechatUnionID: ret.UnionID, WechatOpenID: ret.OpenID})
	case users.RegEmail:
		email := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, def.CaptchaUseLogin, in.CodeID, in.Code)
		if email == "" || email != in.Account {
			return nil, errors.Captcha
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Emails: []string{in.Account}})
		if errors.Cmp(err, errors.NotFind) && cli.Config.IsAutoRegister == def.True { //未注册,自动注册
			err = nil
			userID := l.svcCtx.UserID.GetSnowflakeId()
			uc = &relationDB.SysUserInfo{
				UserID: userID,
				Email:  sql.NullString{Valid: true, String: email},
			}
			err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
				return Register(l.ctx, l.svcCtx, uc, tx)
			})
			uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Emails: []string{in.Account}})
		}
	case users.RegPhone:
		phone := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypePhone, def.CaptchaUseLogin, in.CodeID, in.Code)
		if phone == "" || phone != in.Account {
			return nil, errors.Captcha
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phones: []string{in.Account}})
		if errors.Cmp(err, errors.NotFind) && cli.Config.IsAutoRegister == def.True { //未注册,自动注册
			err = nil
			userID := l.svcCtx.UserID.GetSnowflakeId()
			uc = &relationDB.SysUserInfo{
				UserID: userID,
				Phone:  sql.NullString{Valid: true, String: phone},
			}
			err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
				return Register(l.ctx, l.svcCtx, uc, tx)
			})
			uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phones: []string{in.Account}})
		}
	default:
		l.Error("%s LoginType=%s not support", utils.FuncName(), in.LoginType)
		return nil, errors.Parameter
	}
	l.Infof("%s uc=%#v err=%+v", utils.FuncName(), uc, err)
	return uc, err
}

func (l *LoginLogic) UserLogin(in *sys.UserLoginReq) (*sys.UserLoginResp, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
	//检查账号是否冻结
	list := l.svcCtx.Config.WrongPasswordCounter.ParseWrongPassConf(in.Account, in.Ip)
	if len(list) > 0 {
		forbiddenTime, f := l.svcCtx.PwdCheck.CheckAccountOrIpForbidden(l.ctx, list)
		if f {
			return nil, errors.Default.AddMsgf("%s %d 分钟", errors.AccountOrIpForbidden.Error(), forbiddenTime/60)
		}
	}
	uc := ctxs.GetUserCtx(l.ctx)
	cfg, err := relationDB.NewTenantAppRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantAppFilter{AppCodes: []string{uc.AppCode}})
	if err != nil {
		return nil, err
	}
	if len(cfg.LoginTypes) > 0 && !utils.SliceIn(in.LoginType, cfg.LoginTypes...) {
		return nil, errors.Parameter.WithMsgf("不支持的升级方式:%v", in.LoginType)
	}
	ui, err := l.GetUserInfo(in)
	if err == nil {
		if ui.Status != def.True {
			return nil, errors.AccountDisable
		}
		return l.getRet(ui, list)
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, errors.UnRegister
	}
	if errors.Cmp(err, errors.Password) {
		ret, err := l.svcCtx.PwdCheck.CheckPasswordTimes(l.ctx, list)
		if err != nil {
			return nil, err
		}
		if ret {
			return nil, errors.UseCaptcha
		}
		return nil, errors.Password
	}
	return nil, err
}

func GetAccount(ui *relationDB.SysUserInfo) string {
	var account = ui.UserName.String
	if account == "" {
		account = ui.Phone.String
	}
	if account == "" {
		account = ui.Email.String
	}
	if account == "" {
		account = cast.ToString(ui.UserID)
	}
	return account
}
