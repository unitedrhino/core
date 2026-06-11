package usermanagelogic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/topics"
	"gitee.com/unitedrhino/core/share/users"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/clients/huaweiCli"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"

	"github.com/silenceper/wechat/v2/officialaccount/oauth"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zhaoyunxing92/dingtalk/v2/request"
	"gorm.io/gorm"
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

// syncUserEmailFromGoogle 登录成功后将 Google 返回的邮箱同步到用户表，便于 info.email 回显
func (l *LoginLogic) syncUserEmailFromGoogle(uc *relationDB.SysUserInfo, email string) {
	if uc == nil || email == "" {
		return
	}
	if uc.Email.Valid && uc.Email.String == email {
		return
	}
	oldEmail := uc.Email
	uc.Email = sql.NullString{Valid: true, String: email}
	if err := l.UiDB.Update(l.ctx, uc); err != nil {
		uc.Email = oldEmail
		l.Errorf("%s sync email userID=%d email=%s err=%v", utils.FuncName(), uc.UserID, email, err)
	}
}

func (l *LoginLogic) getPwd(in *sys.UserLoginReq, uc *relationDB.SysUserInfo) error {
	//根据密码类型不同做不同处理
	if in.PwdType == 0 {
		return errors.Parameter.WithMsg("账号密码登录需指定 pwdType：1 明文密码，2 MD5 密码")
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

var randID atomic.Uint32

func genID(ctx context.Context, nodeID int64) string {
	var token = uint32(nodeID) & 0xff
	token += randID.Add(1) << 8 & 0xfff00
	return fmt.Sprintf("%x", token)
}

func GenLoginResp(ctx context.Context, svcCtx *svc.ServiceContext, deviceID string, ui *relationDB.SysUserInfo) (*sys.UserLoginResp, error) {
	uc := ctxs.GetUserCtx(ctx)
	id := genID(ctx, svcCtx.NodeID)
	now := time.Now()
	accessExpire := svcCtx.Config.UserToken.AccessExpire
	jwtToken, claims, err := users.GetLoginJwtToken(svcCtx.Config.UserToken.AccessSecret, now,
		ui.UserID, uc.AppCode, id, deviceID)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return nil, errors.System.AddDetail(err)
	}
	resp := &sys.UserLoginResp{
		Info: UserInfoToPb(ctx, ui, svcCtx),
		Token: &sys.JwtToken{
			AccessToken:  jwtToken,
			AccessExpire: now.Unix() + accessExpire,
			RefreshAfter: now.Unix() + accessExpire/2,
		},
	}
	// 登录前 ctx 可能为 platform 等租户，与 ui.TenantCode 不一致会导致租户配置查询 tenant_code 条件冲突
	tenantCode := string(ui.TenantCode)
	if tenantCode == def.TenantCodePlateform {
		tenantCode = def.TenantCodeDefault
	}
	origTenantCode := uc.TenantCode
	uc.TenantCode = tenantCode
	defer func() { uc.TenantCode = origTenantCode }()
	tc, err := svcCtx.TenantConfigCache.GetData(ctx, tenantCode)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s  err=%s", utils.FuncName(), err.Error())
		return nil, err
	}
	err = svcCtx.UserToken.Login(ctx, claims, accessExpire, tc.IsSsl == def.True)
	if err != nil {
		return nil, err
	}
	logx.WithContext(ctx).Infof("%s GenLoginResp=%+v", utils.FuncName(), resp)
	return resp, nil
}

func (l *LoginLogic) GetUserInfo(in *sys.UserLoginReq) (uc *relationDB.SysUserInfo, err error) {
	cli, er := l.svcCtx.Cm.GetClients(l.ctx, "")
	if er != nil {
		return nil, errors.System.AddDetail(er)
	}
	if in.LoginType != users.RegJwt && !utils.SliceIn(in.LoginType, cli.Config.LoginTypes...) {
		l.Errorf("不支持的登录方式:%v", in.LoginType)
		return nil, errors.NotSupportLogin
	}
	var isRegister bool
	switch in.LoginType {
	case users.RegPwd:
		if in.Code != "" {
			if l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeImage, def.CaptchaUseLogin, in.CodeID, in.Code) == "" {
				return nil, errors.Captcha
			}
		} else if l.svcCtx.LoginLimit.PwdCaptcha.CheckLimit(l.ctx, in.Account) {
			return nil, errors.NeedImgCaptcha
		}
		if l.svcCtx.LoginLimit.PwdAccount.CheckLimit(l.ctx, in.Account) {
			return nil, errors.AccountOrIpForbidden.WithMsg("错误次数过多,请稍后再试")
		}
		ip := ctxs.GetUserCtxNoNil(l.ctx).IP
		if ip != "" && l.svcCtx.LoginLimit.PwdIp.CheckLimit(l.ctx, ip) {
			return nil, errors.AccountOrIpForbidden.WithMsg("错误次数过多,请稍后再试")
		}
		l.svcCtx.LoginLimit.PwdCaptcha.LimitIt(l.ctx, in.Account)
		limit := func() {
			l.svcCtx.LoginLimit.PwdAccount.LimitIt(l.ctx, in.Account)
			if ip != "" {
				l.svcCtx.LoginLimit.PwdIp.LimitIt(l.ctx, in.Ip)
			}
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: []string{in.Account}})
		if err != nil {
			limit()
			return nil, err
		}
		if err = l.getPwd(in, uc); err != nil {
			limit()
			return nil, err
		}
		l.svcCtx.LoginLimit.PwdAccount.CleanLimit(l.ctx, in.Account)
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

		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{DingTalkUserID: ret.UserInfo.UserId, DingTalkUnionID: ret.UserInfo.UnionId})
		if errors.Cmp(err, errors.NotFind) && cli.Config.IsAutoRegister == def.True { //未注册,自动注册
			err = nil
			userID := l.svcCtx.UserID.GetSnowflakeId()
			uc = &relationDB.SysUserInfo{
				UserID:         userID,
				DingTalkUserID: sql.NullString{Valid: true, String: ret.UserInfo.UserId},
				NickName:       ret.UserInfo.Name,
			}
			if ret.UserInfo.UnionId != "" {
				uc.DingTalkUnionID = sql.NullString{Valid: true, String: ret.UserInfo.UnionId}
			}
			ui, er := cli.DingMini.GetUserDetail(&request.UserDetail{
				UserId: ret.UserInfo.UserId,
			})
			l.Infof("GetUserDetail ui:%v err:%v", utils.Fmt(ui), er)
			if er == nil {
				var accounts []string
				if ui.OrgEmail != "" {
					accounts = append(accounts, ui.OrgEmail)
				}
				if ui.Mobile != "" {
					accounts = append(accounts, ui.Mobile)
				}
				uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: accounts})
				if err == nil {
					if ui.OrgEmail != "" {
						uc.Email = sql.NullString{String: ui.OrgEmail, Valid: true}
					}
					if ui.Mobile != "" {
						uc.Phone = sql.NullString{String: ui.Mobile, Valid: true}
					}
					if uc.NickName == "" {
						uc.NickName = ui.Name
					}
					uc.DingTalkUserID = sql.NullString{Valid: true, String: ret.UserInfo.UserId}
					if ret.UserInfo.UnionId != "" {
						uc.DingTalkUnionID = sql.NullString{Valid: true, String: ret.UserInfo.UnionId}
					}
					err = l.UiDB.Update(l.ctx, uc)
					goto end
				}
			}
			uc = &relationDB.SysUserInfo{
				UserID:         userID,
				DingTalkUserID: sql.NullString{Valid: true, String: ret.UserInfo.UserId},
				NickName:       ret.UserInfo.Name,
			}
			if ret.UserInfo.UnionId != "" {
				uc.DingTalkUnionID = sql.NullString{Valid: true, String: ret.UserInfo.UnionId}
			}
			if ui.OrgEmail != "" {
				uc.Email = sql.NullString{String: ui.OrgEmail, Valid: true}
			}
			if ui.Mobile != "" {
				uc.Phone = sql.NullString{String: ui.Mobile, Valid: true}
			}
			if len(ui.Extension) != 0 {
				var tags = map[string]string{}
				err = json.Unmarshal([]byte(ui.Extension), &tags)
				if err == nil {
					uc.Tags = tags
				}
			}
			err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
				return Register(l.ctx, l.svcCtx, uc, tx)
			})
			isRegister = true
			if err != nil {
				return nil, err
			}
			uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{DingTalkUserID: ret.UserInfo.UserId, DingTalkUnionID: ret.UserInfo.UnionId, WithRoles: true, WithTenant: true})
		}
	case users.RegWxOpen:
		if cli.WxOfficial == nil {
			return nil, errors.System.AddDetail(er)
		}
		at, er := cli.WxOfficial.GetOauth().GetUserAccessToken(in.Code)
		if er != nil {
			at2, err := GetWxRegisterResAccessToken(l.ctx, in.Code)
			if err != nil {
				return nil, errors.Default.AddDetail(er)
			}
			at = *at2
		} else {
			StoreWxLoginResAccessToken(l.ctx, in.Code, &at)
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
				UserID:   userID,
				Email:    sql.NullString{Valid: true, String: email},
				UserName: sql.NullString{Valid: true, String: email},
			}
			err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
				return Register(l.ctx, l.svcCtx, uc, tx)
			})
			isRegister = true
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
				UserID:   userID,
				Phone:    sql.NullString{Valid: true, String: phone},
				UserName: sql.NullString{Valid: true, String: phone},
			}
			err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
				return Register(l.ctx, l.svcCtx, uc, tx)
			})
			isRegister = true
			if err != nil {
				return nil, err
			}
			uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phones: []string{in.Account}})
		}
	case users.RegGoogle:
		if cli.Google == nil {
			return nil, errors.System.AddDetail(er)
		}
		gUser, er := cli.Google.ResolveUser(l.ctx, in.Code)
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		tenantCode := ctxs.GetUserCtxNoNil(l.ctx).TenantCode
		var isGoogleAutoRegister bool
		err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
			uidb := relationDB.NewUserInfoRepo(tx)
			uc, err = findOrBindGoogleUser(l.ctx, uidb, tenantCode, gUser)
			if err == nil {
				return nil
			}
			if !errors.Cmp(err, errors.NotFind) {
				return err
			}
			if cli.Config.IsAutoRegister != def.True {
				return err
			}
			userID := l.svcCtx.UserID.GetSnowflakeId()
			googleUser := &relationDB.SysUserInfo{
				UserID:       userID,
				GoogleUserID: sql.NullString{Valid: true, String: gUser.ID},
				NickName:     gUser.Name,
				Email:        sql.NullString{Valid: true, String: gUser.Email},
				HeadImg:      gUser.Picture,
			}
			uc = googleUser
			applyOAuthLoginAccount(uc)
			if err = Register(l.ctx, l.svcCtx, uc, tx); err != nil {
				if errors.Cmp(err, errors.Duplicate) {
					duplicateErr := err
					uc, err = uidb.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{
						TenantCode:   tenantCode,
						GoogleUserID: gUser.ID,
					})
					if err == nil {
						return nil
					}
					return duplicateErr
				}
				return err
			}
			isGoogleAutoRegister = true
			return nil
		})
		isRegister = isGoogleAutoRegister
		if err != nil {
			return nil, err
		}
		if uc != nil {
			freshGoogleUser, findErr := l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{
				TenantCode:   tenantCode,
				GoogleUserID: gUser.ID,
			})
			if findErr != nil {
				return nil, findErr
			} else {
				uc = freshGoogleUser
			}
		}
		if err == nil {
			l.syncUserEmailFromGoogle(uc, gUser.Email)
		}
	case users.RegGithub:
		if cli.Github == nil {
			return nil, errors.System.AddDetail(er)
		}
		token, er := cli.Github.ExchangeCode(l.ctx, in.Code, "")
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		ghUser, er := cli.Github.GetUserInfo(l.ctx, token)
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{GithubUserID: cast.ToString(ghUser.ID)})
		if errors.Cmp(err, errors.NotFind) && cli.Config.IsAutoRegister == def.True { //未注册,自动注册
			err = nil
			userID := l.svcCtx.UserID.GetSnowflakeId()
			uc = &relationDB.SysUserInfo{
				UserID:       userID,
				GithubUserID: sql.NullString{Valid: true, String: cast.ToString(ghUser.ID)},
				NickName:     ghUser.Name,
				Email:        sql.NullString{Valid: true, String: ghUser.Email},
				HeadImg:      ghUser.AvatarURL,
			}
			err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
				return Register(l.ctx, l.svcCtx, uc, tx)
			})
			isRegister = true
			if err != nil {
				return nil, err
			}
			uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{GithubUserID: cast.ToString(ghUser.ID)})
		}
	case users.RegApple:
		if cli.Apple == nil {
			return nil, errors.System.AddMsg("Apple登录未配置或私钥无效，请检查租户应用 Apple 配置")
		}
		aUser, _, er := cli.Apple.ExchangeCode(l.ctx, in.Code)
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		tenantCode := ctxs.GetUserCtxNoNil(l.ctx).TenantCode
		var isAppleAutoRegister bool
		err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
			uidb := relationDB.NewUserInfoRepo(tx)
			uc, err = findOrBindAppleUser(l.ctx, uidb, tenantCode, aUser)
			if err == nil {
				return nil
			}
			if !errors.Cmp(err, errors.NotFind) {
				return err
			}
			if cli.Config.IsAutoRegister != def.True {
				return err
			}
			userID := l.svcCtx.UserID.GetSnowflakeId()
			appleUser := &relationDB.SysUserInfo{
				UserID:      userID,
				AppleUserID: sql.NullString{Valid: true, String: aUser.Sub},
				Email:       sql.NullString{Valid: true, String: aUser.Email},
			}
			uc = appleUser
			applyOAuthLoginAccount(uc)
			if err = Register(l.ctx, l.svcCtx, uc, tx); err != nil {
				if errors.Cmp(err, errors.Duplicate) {
					duplicateErr := err
					uc, err = uidb.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{
						TenantCode:  tenantCode,
						AppleUserID: aUser.Sub,
					})
					if err == nil {
						return nil
					}
					return duplicateErr
				}
				return err
			}
			isAppleAutoRegister = true
			return nil
		})
		isRegister = isAppleAutoRegister
		if err != nil {
			return nil, err
		}
		if uc != nil {
			freshAppleUser, findErr := l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{
				TenantCode:  tenantCode,
				AppleUserID: aUser.Sub,
			})
			if findErr != nil {
				return nil, findErr
			} else {
				uc = freshAppleUser
			}
		}
	case users.RegHuawei:
		if cli.Huawei == nil {
			return nil, errors.System.AddDetail(er)
		}
		// cc:= huawei.NewHuaweiClient(l.ctx,&conf.ThirdConf{AppID: cli.Config.Huawei.AppID, AppSecret: cli.Config.Huawei.AppSecret})
		// cc.GetPhoneNumberByCode(l.ctx, in.Code)
		loginResult, er2 := cli.Huawei.QuickLoginByCode(l.ctx, in.Code)
		if er2 != nil {
			// 如果调用失败，尝试从注册缓存读取
			loginResult, er2 = GetHuaweiRegisterResult(l.ctx, in.Code)
			if er2 != nil {
				return nil, errors.System.AddDetail(er2)
			}
		} else {
			// 调用成功，缓存结果
			StoreHuaweiLoginResult(l.ctx, in.Code, loginResult)
		}
		l.Infof("loginResult=%v,err=%v", utils.Fmt(loginResult), er2)
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{
			HuaweiUnionID: loginResult.UserInfo.UnionID,
			HuaweiOpenID:  loginResult.UserInfo.OpenID,
		})
		// if errors.Cmp(err, errors.NotFind) && cli.Config.IsAutoRegister == def.True { //未注册,自动注册
		// 	err = nil
		// 	if len(loginResult.PhoneNumber) == 0 {
		// 		return nil, errors.UnBindAccount.AddDetail("华为账号未绑定手机号")
		// 	}
		// 	userID := l.svcCtx.UserID.GetSnowflakeId()
		// 	uc = &relationDB.SysUserInfo{
		// 		UserID:        userID,
		// 		HuaweiUnionID: sql.NullString{Valid: true, String: loginResult.UnionID},
		// 		HuaweiOpenID:  sql.NullString{Valid: true, String: loginResult.OpenID},
		// 		Phone:         sql.NullString{Valid: true, String: loginResult.PhoneNumber},
		// 	}
		// 	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 		return Register(l.ctx, l.svcCtx, uc, tx)
		// 	})
		// 	isRegister = true
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{
		// 		HuaweiUnionID: loginResult.UnionID,
		// 		HuaweiOpenID:  loginResult.OpenID,
		// 	})
		// }
	case users.RegJwt:
		account, er2 := ParseThirdJwt(in.Code, l.svcCtx.Config.ThirdJwtSecret)
		if er2 != nil {
			return nil, er2
		}
		if utils.IsPhone(account) {
			uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phones: []string{account}})
		} else if utils.IsEmail(account) {
			uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Emails: []string{account}})
		} else {
			uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: []string{account}})
		}
		if errors.Cmp(err, errors.NotFind) && cli.Config.IsAutoRegister == def.True {
			err = nil
			userID := l.svcCtx.UserID.GetSnowflakeId()
			uc = &relationDB.SysUserInfo{UserID: userID}
			if utils.IsPhone(account) {
				uc.Phone = sql.NullString{Valid: true, String: account}
				uc.UserName = sql.NullString{Valid: true, String: account}
			} else if utils.IsEmail(account) {
				uc.Email = sql.NullString{Valid: true, String: account}
				uc.UserName = sql.NullString{Valid: true, String: account}
			} else {
				uc.UserName = sql.NullString{Valid: true, String: account}
			}
			err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
				return Register(l.ctx, l.svcCtx, uc, tx)
			})
			isRegister = true
			if err != nil {
				return nil, err
			}
			if utils.IsPhone(account) {
				uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phones: []string{account}})
			} else if utils.IsEmail(account) {
				uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Emails: []string{account}})
			} else {
				uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: []string{account}})
			}
		}
	default:
		l.Error("%s LoginType=%s not support", utils.FuncName(), in.LoginType)
		return nil, errors.Parameter
	}
end:
	l.Infof("%s uc=%#v err=%+v", utils.FuncName(), uc, err)
	if isRegister && err == nil {
		e := l.svcCtx.FastEvent.Publish(l.ctx, topics.CoreUserCreate, def.IDs{IDs: []int64{uc.UserID}})
		if e != nil {
			l.Errorf("Publish CoreUserCreate %v err:%v", uc, e)
		}
	}
	return uc, err
}

func (l *LoginLogic) UserLogin(in *sys.UserLoginReq) (*sys.UserLoginResp, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
	uc := ctxs.GetUserCtx(l.ctx)
	cfg, err := relationDB.NewTenantAppRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantAppFilter{AppCodes: []string{uc.AppCode}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Permissions.WithMsg("本租户无该应用")
		}
		return nil, err
	}
	if in.LoginType != users.RegJwt && len(cfg.LoginTypes) > 0 && !utils.SliceIn(in.LoginType, cfg.LoginTypes...) {
		return nil, errors.Parameter.WithMsgf("不支持的登录方式:%v", in.LoginType)
	}
	ui, err := l.GetUserInfo(in)
	if err == nil {
		if ui != nil && ui.Status != def.True {
			return nil, errors.AccountDisable
		}
		return GenLoginResp(l.ctx, l.svcCtx, in.DeviceID, ui)
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, errors.UnRegister
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

func gentLoginKey(code string) string {
	return fmt.Sprintf("sys:user:wxak:login:%s", code)
}

func StoreWxLoginResAccessToken(ctx context.Context, code string, tk *oauth.ResAccessToken) error {
	return caches.GetStore().SetexCtx(ctx, gentLoginKey(code), utils.MarshalNoErr(tk), 10*60)
}

func DelWxLoginResAccessToken(ctx context.Context, code string) error {
	_, err := caches.GetStore().DelCtx(ctx, gentLoginKey(code))
	return err
}

func GetWxLoginResAccessToken(ctx context.Context, code string) (*oauth.ResAccessToken, error) {
	ret, err := caches.GetStore().GetCtx(ctx, gentLoginKey(code))
	if err != nil {
		return nil, err
	}
	var val oauth.ResAccessToken
	err = json.Unmarshal([]byte(ret), &val)
	return &val, err
}

func gentRegisterKey(code string) string {
	return fmt.Sprintf("sys:user:wxak:register:%s", code)
}

func DelWxRegisterResAccessToken(ctx context.Context, code string) error {
	_, err := caches.GetStore().DelCtx(ctx, gentRegisterKey(code))
	return err
}

func StoreWxRegisterResAccessToken(ctx context.Context, code string, tk *oauth.ResAccessToken) error {
	DelWxLoginResAccessToken(ctx, code)
	return caches.GetStore().SetexCtx(ctx, gentRegisterKey(code), utils.MarshalNoErr(tk), 10*60)
}
func GetWxRegisterResAccessToken(ctx context.Context, code string) (*oauth.ResAccessToken, error) {
	ret, err := caches.GetStore().GetCtx(ctx, gentRegisterKey(code))
	if err != nil {
		return nil, err
	}
	var val oauth.ResAccessToken
	err = json.Unmarshal([]byte(ret), &val)
	return &val, err
}

// 华为登录缓存 key
func genHuaweiLoginKey(code string) string {
	return fmt.Sprintf("sys:user:huawei:login:%s", code)
}

// 存储华为登录结果到缓存
func StoreHuaweiLoginResult(ctx context.Context, code string, result *huaweiCli.HuaweiQuickLoginResult) error {
	return caches.GetStore().SetexCtx(ctx, genHuaweiLoginKey(code), utils.MarshalNoErr(result), 10*60)
}

// 从缓存获取华为登录结果
func GetHuaweiLoginResult(ctx context.Context, code string) (*huaweiCli.HuaweiQuickLoginResult, error) {
	ret, err := caches.GetStore().GetCtx(ctx, genHuaweiLoginKey(code))
	if err != nil {
		return nil, err
	}
	var val huaweiCli.HuaweiQuickLoginResult
	err = json.Unmarshal([]byte(ret), &val)
	return &val, err
}

// 删除华为登录缓存
func DelHuaweiLoginResult(ctx context.Context, code string) error {
	_, err := caches.GetStore().DelCtx(ctx, genHuaweiLoginKey(code))
	return err
}

// 华为注册缓存 key
func genHuaweiRegisterKey(code string) string {
	return fmt.Sprintf("sys:user:huawei:register:%s", code)
}

// 存储华为注册结果到缓存
func StoreHuaweiRegisterResult(ctx context.Context, code string, result *huaweiCli.HuaweiQuickLoginResult) error {
	DelHuaweiLoginResult(ctx, code) // 删除登录缓存
	return caches.GetStore().SetexCtx(ctx, genHuaweiRegisterKey(code), utils.MarshalNoErr(result), 10*60)
}

// 从缓存获取华为注册结果
func GetHuaweiRegisterResult(ctx context.Context, code string) (*huaweiCli.HuaweiQuickLoginResult, error) {
	ret, err := caches.GetStore().GetCtx(ctx, genHuaweiRegisterKey(code))
	if err != nil {
		return nil, err
	}
	var val huaweiCli.HuaweiQuickLoginResult
	err = json.Unmarshal([]byte(ret), &val)
	return &val, err
}

// 删除华为注册缓存
func DelHuaweiRegisterResult(ctx context.Context, code string) error {
	_, err := caches.GetStore().DelCtx(ctx, genHuaweiRegisterKey(code))
	return err
}
