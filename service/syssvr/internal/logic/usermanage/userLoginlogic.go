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
	UtDB *relationDB.UserThirdRepo
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
		UtDB:   relationDB.NewUserThirdRepo(ctx),
	}
}

func (l *LoginLogic) getPwd(in *sys.UserLoginReq, uc *relationDB.SysUserInfo) error {
	//根据密码类型不同做不同处理
	if in.PwdType == 0 {
		//空密码情况暂不考虑
		return errors.Password
	} else if in.PwdType == 1 {
		//明文密码，则对密码做MD5加密后再与数据库密码比对
		password1 := utils.MakePwd(in.Password, uc.UserID, false) //对密码进行md5加密
		if password1 != uc.Password {
			return errors.Password
		}
		l.Infof("用户 %d 密码验证成功 (明文密码)", uc.UserID)
	} else if in.PwdType == 2 {
		//md5加密后的密码则通过二次md5加密再对比库中的密码
		password1 := utils.MakePwd(in.Password, uc.UserID, true) //对密码进行md5加密
		if password1 != uc.Password {
			return errors.Password
		}
		l.Infof("用户 %d 密码验证成功 (MD5密码)", uc.UserID)
	} else {
		l.Errorf("用户 %d 使用了不支持的密码类型: %d", uc.UserID, in.PwdType)
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

// autoRegisterUser 自动注册用户的公共逻辑
func (l *LoginLogic) autoRegisterUser(userInfo *relationDB.SysUserInfo) (*relationDB.SysUserInfo, error) {
	userID := l.svcCtx.UserID.GetSnowflakeId()
	userInfo.UserID = userID

	err := stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		return Register(l.ctx, l.svcCtx, userInfo, tx)
	})
	if err != nil {
		l.Errorf("用户自动注册失败: UserID=%d, error=%v", userID, err)
		return nil, err
	}

	// 重新查询用户信息
	uc, err := l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{UserIDs: []int64{userID}})
	if err != nil {
		l.Errorf("查询新注册用户信息失败: UserID=%d, error=%v", userID, err)
		return nil, err
	}

	l.Infof("用户自动注册成功: UserID=%d, NickName=%s", userID, uc.NickName)
	return uc, nil
}

func (l *LoginLogic) getRet(in *sys.UserLoginReq, ui *relationDB.SysUserInfo) (*sys.UserLoginResp, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	id := genID(l.ctx, l.svcCtx.NodeID)
	now := time.Now()
	accessExpire := l.svcCtx.Config.UserToken.AccessExpire
	jwtToken, claims, err := users.GetLoginJwtToken(l.svcCtx.Config.UserToken.AccessSecret, now, accessExpire,
		ui.UserID, uc.AppCode, id, in.DeviceID)
	if err != nil {
		l.Error(err)
		return nil, errors.System.AddDetail(err)
	}
	resp := &sys.UserLoginResp{
		Info: UserInfoToPb(l.ctx, ui, l.svcCtx),
		Token: &sys.JwtToken{
			AccessToken:  jwtToken,
			AccessExpire: now.Unix() + accessExpire,
			RefreshAfter: now.Unix() + accessExpire/2,
		},
	}
	err = l.svcCtx.UserToken.Login(l.ctx, claims)
	if err != nil {
		return nil, err
	}
	l.Infof("%s getRet=%+v", utils.FuncName(), resp)
	return resp, nil
}

func (l *LoginLogic) GetUserInfo(in *sys.UserLoginReq, cfg *relationDB.SysTenantApp) (uc *relationDB.SysUserInfo, err error) {
	//cli, er := l.svcCtx.Cm.GetClients(l.ctx, "")
	//if er != nil {
	//	return nil, errors.System.AddDetail(err)
	//}
	//if !utils.SliceIn(in.LoginType, cli.Config.LoginTypes...) {
	//	l.Errorf("不支持的登录方式:%v", in.LoginType)
	//	return nil, errors.NotSupportLogin
	//}
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
		l.svcCtx.LoginLimit.PwdCaptcha.LimitIt(l.ctx, in.Account)
		if l.svcCtx.LoginLimit.PwdAccount.CheckLimit(l.ctx, in.Account) {
			return nil, errors.AccountOrIpForbidden.WithMsg("错误次数过多,请稍后再试")
		}
		ip := ctxs.GetUserCtxNoNil(l.ctx).IP
		if ip != "" && l.svcCtx.LoginLimit.PwdIp.CheckLimit(l.ctx, ip) {
			return nil, errors.AccountOrIpForbidden.WithMsg("错误次数过多,请稍后再试")
		}
		limit := func() {
			l.svcCtx.LoginLimit.PwdAccount.LimitIt(l.ctx, in.Account)
			if ip != "" {
				l.svcCtx.LoginLimit.PwdIp.LimitIt(l.ctx, ip)
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
		cli, err := l.svcCtx.ThirdClientsManage.GetDingAppClient(l.ctx, ctxs.GetAppCode(l.ctx), in.AppID)
		if err != nil {
			l.Errorf("获取钉钉客户端失败: AppID=%s, error=%v", in.AppID, err)
			return nil, err
		}
		ret, er := cli.GetUserInfoByCode(in.Code)
		if er != nil {
			l.Errorf("钉钉获取用户信息失败: Code=%s, error=%v", in.Code, er)
			return nil, errors.System.AddDetail(er)
		}
		if ret.Code != 0 {
			l.Errorf("钉钉API返回错误: Code=%d, Msg=%s", ret.Code, ret.Msg)
			return nil, errors.Parameter.AddMsgf(ret.Msg)
		}

		ut, err := l.UtDB.FindOneByFilter(l.ctx, relationDB.UserThirdFilter{
			WithUser: true, AppType: def.ThirdTypeDingApp, OpenID: ret.UserInfo.UserId, UnionID: ret.UserInfo.UnionId})
		if err != nil && !errors.Cmp(err, errors.NotFind) {
			l.Errorf("查询钉钉第三方登录信息失败: OpenID=%s, UnionID=%s, error=%v", ret.UserInfo.UserId, ret.UserInfo.UnionId, err)
			return nil, err
		}
		if err != nil && cfg.IsAutoRegister == def.True { //未注册,自动注册
			err = nil
			ui, er := cli.GetUserDetail(&request.UserDetail{
				UserId: ret.UserInfo.UserId,
			})
			if er != nil {
				return nil, errors.System.AddMsg("无法获取钉钉信息,需检查是否授权").AddDetail(er)
			}
			l.Infof("GetUserDetail ui:%v err:%v", utils.Fmt(ui), er)
			var accounts []string
			if ui.OrgEmail != "" {
				accounts = append(accounts, ui.OrgEmail)
			}
			if ui.Mobile != "" {
				accounts = append(accounts, ui.Mobile)
			}
			if len(accounts) == 0 {
				return nil, errors.AccountDisable.AddMsg("钉钉需要先绑定邮箱或手机号")
			}
			uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: accounts})
			if err == nil {
				// 更新现有用户信息
				if ui.OrgEmail != "" {
					uc.Email = sql.NullString{String: ui.OrgEmail, Valid: true}
				}
				if ui.Mobile != "" {
					uc.Phone = sql.NullString{String: ui.Mobile, Valid: true}
				}
				err = l.UiDB.Update(l.ctx, uc)
				goto end
			} else if !errors.Cmp(err, errors.NotFind) {
				return nil, err
			}
			// 创建新用户
			// 创建新用户信息
			uc = &relationDB.SysUserInfo{
				NickName: ret.UserInfo.Name,
			}
			// 设置邮箱和手机号
			if ui.OrgEmail != "" {
				uc.Email = sql.NullString{String: ui.OrgEmail, Valid: true}
			}
			if ui.Mobile != "" {
				uc.Phone = sql.NullString{String: ui.Mobile, Valid: true}
			}
			// 如果钉钉昵称为空，使用钉钉用户详情中的姓名
			if uc.NickName == "" {
				uc.NickName = ui.Name
			}
			// 处理扩展信息
			if len(ui.Extension) != 0 {
				var tags = map[string]string{}
				err = json.Unmarshal([]byte(ui.Extension), &tags)
				if err == nil {
					uc.Tags = tags
				}
			}
			// 设置第三方登录信息
			uc.Thirds = []*relationDB.SysUserThird{
				{AppType: def.ThirdTypeDingApp, AppID: in.AppID, UnionID: ui.UnionId, OpenID: ui.UserId},
			}
			// 自动注册用户
			uc, err = l.autoRegisterUser(uc)
			isRegister = true
			if err != nil {
				return nil, err
			}
			// 查询第三方登录关联
			ut, err = l.UtDB.FindOneByFilter(l.ctx, relationDB.UserThirdFilter{
				WithUser: true, AppType: def.ThirdTypeDingApp, OpenID: ret.UserInfo.UserId, UnionID: ret.UserInfo.UnionId})
			if err != nil {
				return nil, err
			}
			uc = ut.User
		} else {
			uc = ut.User
		}
	case users.RegWxOpen:
		cli, err := l.svcCtx.ThirdClientsManage.GetWxOpenClient(l.ctx, ctxs.GetAppCode(l.ctx), in.AppID)
		if err != nil {
			l.Errorf("获取微信开放平台客户端失败: AppID=%s, error=%v", in.AppID, err)
			return nil, err
		}
		at, er := cli.GetOauth().GetUserAccessToken(in.Code)
		if er != nil {
			l.Errorf("微信开放平台获取access token失败，尝试从注册缓存获取: Code=%s, error=%v", in.Code, er)
			at2, err := GetWxRegisterResAccessToken(l.ctx, in.Code)
			if err != nil {
				l.Errorf("从注册缓存获取微信token失败: Code=%s, error=%v", in.Code, err)
				return nil, errors.Default.AddDetail(er)
			}
			at = *at2
		} else {
			StoreWxLoginResAccessToken(l.ctx, in.Code, &at)
		}
		ut, err := l.UtDB.FindOneByFilter(l.ctx, relationDB.UserThirdFilter{AppType: def.ThirdTypeWx, UnionID: at.UnionID, OpenID: at.OpenID})
		if err != nil {
			l.Errorf("查询微信第三方登录信息失败: UnionID=%s, OpenID=%s, error=%v", at.UnionID, at.OpenID, err)
			return nil, err
		}
		uc = ut.User
	case users.RegWxMiniP:
		cli, err := l.svcCtx.ThirdClientsManage.GetWxMiniClient(l.ctx, ctxs.GetAppCode(l.ctx), in.AppID)
		if err != nil {
			l.Errorf("获取微信小程序客户端失败: AppID=%s, error=%v", in.AppID, err)
			return nil, err
		}
		auth := cli.GetAuth()
		ret, er := auth.Code2SessionContext(l.ctx, in.Code)
		if er != nil {
			l.Errorf("微信小程序Code2Session失败: Code=%s, error=%v", in.Code, er)
			return nil, errors.System.AddDetail(er)
		}
		if ret.ErrCode != 0 {
			l.Errorf("微信小程序API返回错误: ErrCode=%d, ErrMsg=%s", ret.ErrCode, ret.ErrMsg)
			return nil, errors.Parameter.AddMsgf(ret.ErrMsg)
		}
		ut, err := l.UtDB.FindOneByFilter(l.ctx, relationDB.UserThirdFilter{AppType: def.ThirdTypeWx, UnionID: ret.UnionID, OpenID: ret.OpenID})
		if err != nil {
			l.Errorf("查询微信小程序第三方登录信息失败: UnionID=%s, OpenID=%s, error=%v", ret.UnionID, ret.OpenID, err)
			return nil, err
		}
		uc = ut.User
	case users.RegEmail:
		email := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, def.CaptchaUseLogin, in.CodeID, in.Code)
		if email == "" || email != in.Account {
			return nil, errors.Captcha
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Emails: []string{in.Account}})
		if err != nil && !errors.Cmp(err, errors.NotFind) {
			l.Errorf("查询邮箱用户信息失败: Email=%s, error=%v", in.Account, err)
			return nil, err
		}
		if errors.Cmp(err, errors.NotFind) && cfg.IsAutoRegister == def.True { //未注册,自动注册
			err = nil
			uc = &relationDB.SysUserInfo{
				Email:    sql.NullString{Valid: true, String: email},
				UserName: sql.NullString{Valid: true, String: email},
			}
			uc, err = l.autoRegisterUser(uc)
			isRegister = true
		}
	case users.RegPhone:
		phone := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypePhone, def.CaptchaUseLogin, in.CodeID, in.Code)
		if phone == "" || phone != in.Account {
			return nil, errors.Captcha
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phones: []string{in.Account}})
		if err != nil && !errors.Cmp(err, errors.NotFind) {
			l.Errorf("查询手机号用户信息失败: Phone=%s, error=%v", in.Account, err)
			return nil, err
		}
		if errors.Cmp(err, errors.NotFind) && cfg.IsAutoRegister == def.True { //未注册,自动注册
			err = nil
			uc = &relationDB.SysUserInfo{
				Phone:    sql.NullString{Valid: true, String: phone},
				UserName: sql.NullString{Valid: true, String: phone},
			}
			uc, err = l.autoRegisterUser(uc)
			isRegister = true
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
		return nil, err
	}
	if len(cfg.LoginTypes) > 0 && !utils.SliceIn(in.LoginType, cfg.LoginTypes...) {
		return nil, errors.Parameter.WithMsgf("不支持的登录方式:%v", in.LoginType)
	}
	ui, err := l.GetUserInfo(in, cfg)
	if err == nil {
		if ui.Status != def.True {
			return nil, errors.AccountDisable
		}
		return l.getRet(in, ui)
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

// 微信登录token缓存管理
type WxTokenCache struct {
	prefix string
	ttl    int
}

func NewWxTokenCache(prefix string) *WxTokenCache {
	return &WxTokenCache{
		prefix: prefix,
		ttl:    10 * 60, // 10分钟
	}
}

func (w *WxTokenCache) getKey(code string) string {
	return fmt.Sprintf("sys:user:wxak:%s:%s", w.prefix, code)
}

func (w *WxTokenCache) Store(ctx context.Context, code string, tk *oauth.ResAccessToken) error {
	return caches.GetStore().SetexCtx(ctx, w.getKey(code), utils.MarshalNoErr(tk), w.ttl)
}

func (w *WxTokenCache) Get(ctx context.Context, code string) (*oauth.ResAccessToken, error) {
	ret, err := caches.GetStore().GetCtx(ctx, w.getKey(code))
	if err != nil {
		return nil, err
	}
	var val oauth.ResAccessToken
	err = json.Unmarshal([]byte(ret), &val)
	return &val, err
}

func (w *WxTokenCache) Delete(ctx context.Context, code string) error {
	_, err := caches.GetStore().DelCtx(ctx, w.getKey(code))
	return err
}

// 全局微信token缓存实例
var (
	wxLoginCache    = NewWxTokenCache("login")
	wxRegisterCache = NewWxTokenCache("register")
)

// 兼容性函数，保持向后兼容
func gentLoginKey(code string) string {
	return wxLoginCache.getKey(code)
}

func StoreWxLoginResAccessToken(ctx context.Context, code string, tk *oauth.ResAccessToken) error {
	return wxLoginCache.Store(ctx, code, tk)
}

func DelWxLoginResAccessToken(ctx context.Context, code string) error {
	return wxLoginCache.Delete(ctx, code)
}

func GetWxLoginResAccessToken(ctx context.Context, code string) (*oauth.ResAccessToken, error) {
	return wxLoginCache.Get(ctx, code)
}

func gentRegisterKey(code string) string {
	return wxRegisterCache.getKey(code)
}

func DelWxRegisterResAccessToken(ctx context.Context, code string) error {
	return wxRegisterCache.Delete(ctx, code)
}

func StoreWxRegisterResAccessToken(ctx context.Context, code string, tk *oauth.ResAccessToken) error {
	DelWxLoginResAccessToken(ctx, code)
	return wxRegisterCache.Store(ctx, code, tk)
}

func GetWxRegisterResAccessToken(ctx context.Context, code string) (*oauth.ResAccessToken, error) {
	return wxRegisterCache.Get(ctx, code)
}
