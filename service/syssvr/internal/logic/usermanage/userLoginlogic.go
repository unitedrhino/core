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
	"gitee.com/unitedrhino/core/share/dataType"
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
	"gorm.io/gorm"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB  *relationDB.UserInfoRepo
	UttDB *relationDB.UserThirdRepo
	UtDB  *relationDB.UserTenantRepo
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
		UttDB:  relationDB.NewUserThirdRepo(ctx),
		UtDB:   relationDB.NewUserTenantRepo(ctx),
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

var randID atomic.Uint32

func genID(ctx context.Context, nodeID int64) string {
	var token = uint32(nodeID) & 0xff
	token += randID.Add(1) << 8 & 0xfff00
	return fmt.Sprintf("%x", token)
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

func (l *LoginLogic) GetUserInfo(in *sys.UserLoginReq) (uc *relationDB.SysUserInfo, err error) {
	//cli, er := l.svcCtx.Cm.GetClients(l.ctx, "")
	//if er != nil {
	//	return nil, errors.System.AddDetail(err)
	//}
	//if !utils.SliceIn(in.LoginType, cli.Config.LoginTypes...) {
	//	l.Errorf("不支持的登录方式:%v", in.LoginType)
	//	return nil, errors.NotSupportLogin
	//}
	ucc := ctxs.GetUserCtx(l.ctx)
	ta, err := relationDB.NewTenantAppRepo(l.ctx).FindOneByFilter(ctxs.CommonWithDefault(l.ctx),
		relationDB.TenantAppFilter{AppCode: ucc.AppCode, WithApp: true})
	if err != nil {
		return nil, errors.NotEnable.WithMsg("未启用应用")
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
				l.svcCtx.LoginLimit.PwdIp.LimitIt(l.ctx, in.Ip)
			}
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: []string{in.Account}, WithTenants: true, TenantStatus: def.True})
		if err != nil {
			limit()
			return nil, err
		}
		if err = l.getPwd(in, uc); err != nil {
			limit()
			return nil, err
		}
		l.svcCtx.LoginLimit.PwdAccount.CleanLimit(l.ctx, in.Account)
	case users.RegDingApp: //todo
		//if cli.DingMini == nil {
		//	return nil, errors.System.AddDetail(err)
		//}
		//ret, er := cli.DingMini.GetUserInfoByCode(in.Code)
		//if er != nil {
		//	return nil, errors.System.AddDetail(er)
		//}
		//if ret.Code != 0 {
		//	return nil, errors.Parameter.AddMsgf(ret.Msg)
		//}
		//
		//ut, err := l.UtDB.FindOneByFilter(l.ctx, relationDB.UserTenantFilter{
		//	DingTalkUserID: ret.UserInfo.UserId, DingTalkUnionID: ret.UserInfo.UnionId, WithUser: true})
		//if errors.Cmp(err, errors.NotFind) && cli.Config.IsAutoRegister == def.True { //未注册,自动注册
		//	err = nil
		//	userID := l.svcCtx.UserID.GetSnowflakeId()
		//	uc = &relationDB.SysUserInfo{
		//		UserID:         userID,
		//		DingTalkUserID: sql.NullString{Valid: true, String: ret.UserInfo.UserId},
		//		NickName:       ret.UserInfo.Name,
		//	}
		//	if ret.UserInfo.UnionId != "" {
		//		uc.DingTalkUnionID = sql.NullString{Valid: true, String: ret.UserInfo.UnionId}
		//	}
		//	ui, er := cli.DingMini.GetUserDetail(&request.UserDetail{
		//		UserId: ret.UserInfo.UserId,
		//	})
		//	l.Infof("GetUserDetail ui:%v err:%v", utils.Fmt(ui), er)
		//	if er == nil {
		//		var accounts []string
		//		if ui.OrgEmail != "" {
		//			accounts = append(accounts, ui.OrgEmail)
		//		}
		//		if ui.Mobile != "" {
		//			accounts = append(accounts, ui.Mobile)
		//		}
		//		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Accounts: accounts})
		//		if err == nil {
		//			if ui.OrgEmail != "" {
		//				uc.Email = sql.NullString{String: ui.OrgEmail, Valid: true}
		//			}
		//			if ui.Mobile != "" {
		//				uc.Phone = sql.NullString{String: ui.Mobile, Valid: true}
		//			}
		//			if uc.NickName == "" {
		//				uc.NickName = ui.Name
		//			}
		//			uc.DingTalkUserID = sql.NullString{Valid: true, String: ret.UserInfo.UserId}
		//			if ret.UserInfo.UnionId != "" {
		//				uc.DingTalkUnionID = sql.NullString{Valid: true, String: ret.UserInfo.UnionId}
		//			}
		//			err = l.UiDB.Update(l.ctx, uc)
		//			goto end
		//		}
		//	}
		//	uc = &relationDB.SysUserInfo{
		//		UserID:         userID,
		//		DingTalkUserID: sql.NullString{Valid: true, String: ret.UserInfo.UserId},
		//		NickName:       ret.UserInfo.Name,
		//	}
		//	if ret.UserInfo.UnionId != "" {
		//		uc.DingTalkUnionID = sql.NullString{Valid: true, String: ret.UserInfo.UnionId}
		//	}
		//	if ui.OrgEmail != "" {
		//		uc.Email = sql.NullString{String: ui.OrgEmail, Valid: true}
		//	}
		//	if ui.Mobile != "" {
		//		uc.Phone = sql.NullString{String: ui.Mobile, Valid: true}
		//	}
		//	if len(ui.Extension) != 0 {
		//		var tags = map[string]string{}
		//		err = json.Unmarshal([]byte(ui.Extension), &tags)
		//		if err == nil {
		//			uc.Tags = tags
		//		}
		//	}
		//	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		//		return Register(l.ctx, l.svcCtx, uc, tx)
		//	})
		//	isRegister = true
		//	if err != nil {
		//		return nil, err
		//	}
		//	uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{DingTalkUserID: ret.UserInfo.UserId, DingTalkUnionID: ret.UserInfo.UnionId, WithRoles: true, WithTenants: true,TenantStatus: def.True})
		//}
	case users.RegWxOpen: //todo
		//if cli.WxOfficial == nil {
		//	return nil, errors.System.AddDetail(er)
		//}
		//at, er := cli.WxOfficial.GetOauth().GetUserAccessToken(in.Code)
		//if er != nil {
		//	at2, err := GetWxRegisterResAccessToken(l.ctx, in.Code)
		//	if err != nil {
		//		return nil, errors.Default.AddDetail(er)
		//	}
		//	at = *at2
		//} else {
		//	StoreWxLoginResAccessToken(l.ctx, in.Code, &at)
		//}
		//uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{WechatUnionID: at.UnionID, WechatOpenID: at.OpenID})
	case users.RegWxMiniP: //todo
		//cli, err := l.svcCtx.ThirdClientsManage.GetWxMiniClient(l.ctx, ctxs.GetUserCtxNoNil(l.ctx).AppCode, in.AppID)
		//if err == nil {
		//	return nil, err
		//}
		//auth := cli.GetAuth()
		//ret, er := auth.Code2SessionContext(l.ctx, in.Code)
		//if er != nil {
		//	return nil, errors.System.AddDetail(er)
		//}
		//if ret.ErrCode != 0 {
		//	return nil, errors.Parameter.AddMsgf(ret.ErrMsg)
		//}
		//var tk string
		//
		//uis, err := l.UttDB.FindByFilter(ctxs.WithDefaultRoot(l.ctx), relationDB.UserThirdFilter{WithUser: true, TenantCode: tk,
		//	AppType: def.ThirdTypeWxMiniP, AppID: in.AppID, UnionID: ret.UnionID, OpenID: ret.OpenID}, nil)
		//if err != nil {
		//	return nil, err
		//}
		//if len(uis) == 0 {
		//	return nil, errors.UnRegister
		//}
		//uc = uis[0].User
		//var tenants []string
		//
		//var ts []*relationDB.SysUserTenant
		//for _, u := range uis {
		//	tenants = append(tenants, string(u.TenantCode))
		//}

	case users.RegEmail:
		email := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, def.CaptchaUseLogin, in.CodeID, in.Code)
		if email == "" || email != in.Account {
			return nil, errors.Captcha
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{WithTenants: true, TenantStatus: def.True, Emails: []string{in.Account}})
		if errors.Cmp(err, errors.NotFind) && ta.IsAutoRegister == def.True { //未注册,自动注册
			err = nil
			userID := l.svcCtx.UserID.GetSnowflakeId()
			uc = &relationDB.SysUserInfo{
				UserID:   userID,
				Email:    sql.NullString{Valid: true, String: email},
				UserName: sql.NullString{Valid: true, String: email},
				Tenants: []*relationDB.SysUserTenant{
					{
						TenantCode: dataType.TenantCode(ucc.TenantCode),
						UserID:     userID,
					},
				},
			}
			err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
				return Register(l.ctx, l.svcCtx, uc, tx)
			})
			isRegister = true
			uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{WithTenants: true, TenantStatus: def.True, Emails: []string{in.Account}})
		}
	case users.RegPhone:
		phone := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypePhone, def.CaptchaUseLogin, in.CodeID, in.Code)
		if phone == "" || phone != in.Account {
			return nil, errors.Captcha
		}
		uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{WithTenants: true, TenantStatus: def.True, Phones: []string{in.Account}})
		if errors.Cmp(err, errors.NotFind) && ta.IsAutoRegister == def.True { //未注册,自动注册
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
			uc, err = l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{WithTenants: true, TenantStatus: def.True, Phones: []string{in.Account}})
		}
	default:
		l.Error("%s LoginType=%s not support", utils.FuncName(), in.LoginType)
		return nil, errors.Parameter
	}
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
	ui, err := l.GetUserInfo(in)
	if err == nil {
		//if ui.Status != def.True {
		//	return nil, errors.AccountDisable
		//}
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
