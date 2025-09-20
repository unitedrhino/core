package usermanagelogic

import (
	"context"
	"database/sql"

	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/caches"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/core/share/topics"
	"gitee.com/unitedrhino/core/share/users"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"

	"time"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	uc *ctxs.UserCtx
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		uc:     ctxs.GetUserCtx(ctx),
	}
}

// verifyCaptcha 验证码验证的公共逻辑
func (l *UserRegisterLogic) verifyCaptcha(captchaType, codeID, code, account string) (string, error) {
	verified := l.svcCtx.Captcha.Verify(l.ctx, captchaType, def.CaptchaUseRegister, codeID, code)
	if verified == "" || verified != account {
		return "", errors.Captcha
	}
	return verified, nil
}

// validatePassword 密码验证的公共逻辑
func (l *UserRegisterLogic) validatePassword(password string) error {
	if password == "" {
		return nil
	}
	err := logic.CheckPwd(l.svcCtx, password)
	if err != nil {
		return err
	}
	return nil
}

// createUserInfo 创建用户信息的公共逻辑
func (l *UserRegisterLogic) createUserInfo(userID int64, nickName, userName string) relationDB.SysUserInfo {
	ui := relationDB.SysUserInfo{
		UserID:   userID,
		NickName: nickName,
	}
	if userName != "" {
		ui.UserName = utils.AnyToNullString(userName)
	}
	return ui
}

func (l *UserRegisterLogic) UserRegister(in *sys.UserRegisterReq) (*sys.UserRegisterResp, error) {
	switch in.RegType {
	case users.RegDingApp:
		return l.handleDingApp(in)
	case users.RegWxMiniP:
		return l.handleWxminiP(in)
	case users.RegEmail, users.RegPhone:
		return l.handleEmailOrPhone(in)
	default:
		return nil, errors.NotRealize.AddMsgf(in.RegType)
	}
}

func (l *UserRegisterLogic) handleEmailOrPhone(in *sys.UserRegisterReq) (*sys.UserRegisterResp, error) {

	userID := l.svcCtx.UserID.GetSnowflakeId()
	var nickName, userName string
	if in.Info != nil {
		nickName = in.Info.NickName
		userName = in.Info.UserName
	}

	ui := l.createUserInfo(userID, nickName, userName)

	// 验证密码
	if err := l.validatePassword(in.Password); err != nil {
		return nil, err
	}

	// 验证验证码并设置账户信息
	switch in.RegType {
	case users.RegEmail:
		email, err := l.verifyCaptcha(def.CaptchaTypeEmail, in.CodeID, in.Code, in.Account)
		if err != nil {
			return nil, err
		}
		ui.Email = utils.AnyToNullString(email)
		ui.UserName = ui.Email
		l.Infof("邮箱验证成功: Email=%s", email)
	case users.RegPhone:
		phone, err := l.verifyCaptcha(def.CaptchaTypePhone, in.CodeID, in.Code, in.Account)
		if err != nil {
			return nil, err
		}
		ui.Phone = utils.AnyToNullString(phone)
		ui.UserName = ui.Phone
	}

	// 设置密码
	if in.Password != "" {
		ui.Password = utils.MakePwd(in.Password, ui.UserID, false)
	}

	wxOpenCode := in.Expand["wxOpenCode"]
	if wxOpenCode != "" {
		at, err := GetWxLoginResAccessToken(l.ctx, wxOpenCode)
		if err != nil {
			cli, er := l.svcCtx.ThirdClientsManage.GetWxOpenClient(l.ctx, ctxs.GetAppCode(l.ctx), in.AppID)
			if er != nil {
				return nil, er
			}
			at2, er := cli.GetOauth().GetUserAccessToken(in.Code)
			if er != nil {
				return nil, errors.System.AddDetail(er)
			}
			at = &at2
		}
		StoreWxRegisterResAccessToken(l.ctx, wxOpenCode, at)
		_, err = relationDB.NewUserThirdRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserThirdFilter{AppType: def.ThirdTypeWx, UnionID: at.UnionID, OpenID: at.OpenID})
		if err == nil {
			return nil, errors.BindAccount.WithMsg("微信已绑定其他账号")
		}
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
		ui.Thirds = []*relationDB.SysUserThird{
			{AppType: def.ThirdTypeWx, AppID: in.AppID, UserID: ui.UserID, UnionID: at.UnionID, OpenID: at.OpenID},
		}
	}
	conn := stores.GetTenantConn(l.ctx)
	err := l.FillUserInfo(&ui, conn)
	if err != nil {
		return nil, err
	}

	// 发布用户创建事件
	e := l.svcCtx.FastEvent.Publish(l.ctx, topics.CoreUserCreate, def.IDs{IDs: []int64{ui.UserID}})
	if e != nil {
		l.Errorf("发布用户创建事件失败: UserID=%d, error=%v", ui.UserID, e)
	}

	return &sys.UserRegisterResp{
		UserID: ui.UserID,
	}, nil
}

func (l *UserRegisterLogic) handleWxminiP(in *sys.UserRegisterReq) (*sys.UserRegisterResp, error) {

	cli, err := l.svcCtx.ThirdClientsManage.GetWxMiniClient(l.ctx, ctxs.GetUserCtxNoNil(l.ctx).AppCode, in.AppID)
	if err != nil {
		l.Errorf("获取微信小程序客户端失败: AppID=%s, error=%v", in.AppID, err)
		return nil, err
	}
	auth := cli.GetAuth()

	wxUid, err := auth.Code2SessionContext(l.ctx, in.Code)
	if err != nil {
		l.Errorf("微信小程序Code2Session失败: Code=%s, error=%v", in.Code, err)
		return nil, errors.System.AddDetail(err)
	}
	if wxUid.ErrCode != 0 {
		l.Errorf("微信小程序API返回错误: ErrCode=%d, ErrMsg=%s", wxUid.ErrCode, wxUid.ErrMsg)
		return nil, errors.Parameter.AddDetail(wxUid.ErrMsg)
	}

	if in.Expand == nil || in.Expand["phoneCode"] == "" {
		return nil, errors.Parameter.AddMsg("微信小程序注册需要填写expand.phoneCode")
	}
	phoneCode := in.Expand["phoneCode"]
	wxPhone, err := auth.GetPhoneNumber(phoneCode)
	if err != nil {
		l.Errorf("获取微信小程序手机号失败: PhoneCode=%s, error=%v", phoneCode, err)
		return nil, errors.System.AddDetail(err)
	}
	if wxPhone.ErrCode != 0 {
		l.Errorf("微信小程序手机号API返回错误: ErrCode=%d, ErrMsg=%s", wxPhone.ErrCode, wxPhone.ErrMsg)
		return nil, errors.Parameter.AddDetail(wxPhone.ErrMsg)
	}

	// 验证密码
	if err := l.validatePassword(in.Password); err != nil {
		return nil, err
	}
	var userID int64
	conn := stores.GetTenantConn(l.ctx)
	err = conn.Transaction(func(tx *gorm.DB) error {
		utdb := relationDB.NewUserThirdRepo(tx)
		uidb := relationDB.NewUserInfoRepo(tx)

		// 检查微信账号是否已注册
		_, err = utdb.FindOneByFilter(l.ctx,
			relationDB.UserThirdFilter{AppType: def.ThirdTypeWx, AppID: in.AppID, UnionID: wxUid.UnionID, OpenID: wxUid.OpenID})
		if err == nil { //已经注册过
			return errors.DuplicateRegister
		}
		if err != nil && !errors.Cmp(err, errors.NotFind) {
			l.Errorf("查询微信第三方登录信息失败: UnionID=%s, OpenID=%s, error=%v", wxUid.UnionID, wxUid.OpenID, err)
			return err
		}

		// 检查手机号是否已注册
		ui, err := uidb.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{WithThird: true, Phones: []string{wxPhone.PhoneInfo.PurePhoneNumber}})
		if err != nil {
			if !errors.Cmp(err, errors.NotFind) {
				l.Errorf("查询手机号用户信息失败: Phone=%s, error=%v", wxPhone.PhoneInfo.PurePhoneNumber, err)
				return err
			}
			//没有注册过，创建新用户
			userID = l.svcCtx.UserID.GetSnowflakeId()
			ui = &relationDB.SysUserInfo{
				UserID:   userID,
				Phone:    sql.NullString{Valid: true, String: wxPhone.PhoneInfo.PurePhoneNumber},
				UserName: sql.NullString{Valid: true, String: wxPhone.PhoneInfo.PurePhoneNumber},
				Thirds:   []*relationDB.SysUserThird{{AppType: def.ThirdTypeWx, AppID: in.AppID, UserID: userID, UnionID: wxUid.UnionID, OpenID: wxUid.OpenID}},
			}
			if in.Password != "" {
				ui.Password = utils.MakePwd(in.Password, ui.UserID, false)
			}
			return l.FillUserInfo(ui, tx)
		} else { //有人注册过,绑定即可
			err = utdb.Insert(l.ctx, &relationDB.SysUserThird{AppType: def.ThirdTypeWx, AppID: in.AppID, UserID: ui.UserID, UnionID: wxUid.UnionID, OpenID: wxUid.OpenID})
			if err != nil {
				l.Errorf("绑定微信失败: UserID=%d, error=%v", ui.UserID, err)
				return err
			}
			userID = ui.UserID
			return nil
		}
	})

	return &sys.UserRegisterResp{UserID: userID}, err
}

func (l *UserRegisterLogic) handleDingApp(in *sys.UserRegisterReq) (*sys.UserRegisterResp, error) {

	cli, err := l.svcCtx.Cm.GetClients(l.ctx, "")
	if err != nil || cli.DingMini == nil {
		l.Errorf("获取钉钉客户端失败: error=%v", err)
		return nil, errors.System.AddDetail(err)
	}
	ret, err := cli.DingMini.GetUserInfoByCode(in.Code)
	if err != nil {
		l.Errorf("钉钉获取用户信息失败: Code=%s, error=%v", in.Code, err)
		return nil, errors.System.AddDetail(err)
	}
	if ret.Code != 0 {
		l.Errorf("钉钉API返回错误: Code=%d, Msg=%s", ret.Code, ret.Msg)
		return nil, errors.Parameter.AddDetail(ret.Msg)
	}

	userID := l.svcCtx.UserID.GetSnowflakeId()
	var nickName, userName string
	if in.Info != nil {
		nickName = in.Info.NickName
		userName = in.Info.UserName
	}
	if nickName == "" {
		nickName = ret.UserInfo.Name
	}

	ui := l.createUserInfo(userID, nickName, userName)
	ui.Thirds = []*relationDB.SysUserThird{
		{AppType: def.ThirdTypeDingApp, AppID: in.AppID, UserID: userID, OpenID: ret.UserInfo.UserId, UnionID: ret.UserInfo.UnionId},
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		uidb := relationDB.NewUserThirdRepo(tx)

		// 检查钉钉账号是否已注册
		_, err = uidb.FindOneByFilter(l.ctx, relationDB.UserThirdFilter{
			AppType: def.ThirdTypeDingApp, OpenID: ret.UserInfo.UserId, UnionID: ret.UserInfo.UnionId})
		if err == nil { //已经注册过
			return errors.DuplicateRegister
		}
		if err != nil && !errors.Cmp(err, errors.NotFind) {
			l.Errorf("查询钉钉第三方登录信息失败: OpenID=%s, UnionID=%s, error=%v", ret.UserInfo.UserId, ret.UserInfo.UnionId, err)
			return err
		}

		return l.FillUserInfo(&ui, tx)
	})
	return &sys.UserRegisterResp{UserID: userID}, err
}

func (l *UserRegisterLogic) FillUserInfo(in *relationDB.SysUserInfo, tx *gorm.DB) error {
	err := Register(l.ctx, l.svcCtx, in, tx)
	if err != nil && errors.Cmp(err, errors.Duplicate) { //已经注册过
		return errors.DuplicateRegister
	}
	return err
}

func Register(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.SysUserInfo, tx *gorm.DB) error {
	ctx = ctxs.WithAdmin(ctx)
	uc := ctxs.GetUserCtx(ctx)
	cfg, err := svcCtx.TenantConfigCache.GetData(ctx, uc.TenantCode)
	if err != nil {
		logx.WithContext(ctx).Errorf("获取租户配置失败: TenantCode=%s, error=%v", uc.TenantCode, err)
		return err
	}

	err = tx.Transaction(func(tx *gorm.DB) error {
		uidb := relationDB.NewUserInfoRepo(tx)
		in.RegIP = uc.IP

		// 插入用户信息
		err = uidb.Insert(ctx, in)
		if err != nil {
			logx.WithContext(ctx).Errorf("插入用户信息失败: UserID=%d, error=%v", in.UserID, err)
			return err
		}

		// 插入用户角色关联
		err = relationDB.NewUserRoleRepo(tx).Insert(ctx, &relationDB.SysUserRole{
			UserID: in.UserID,
			RoleID: cfg.RegisterRoleID,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("插入用户角色关联失败: UserID=%d, RoleID=%d, error=%v", in.UserID, cfg.RegisterRoleID, err)
			return err
		}

		return nil
	})
	if err == nil && len(cfg.RegisterAutoCreateProject) > 0 {
		ctxs.GoNewCtx(ctx, func(ctx context.Context) {
			var pis []*relationDB.SysProjectInfo
			var dps []*relationDB.SysDataProject
			var ais []*relationDB.SysAreaInfo
			for _, rap := range cfg.RegisterAutoCreateProject {
				po := relationDB.SysProjectInfo{
					ProjectID:   dataType.ProjectID(svcCtx.ProjectID.GetSnowflakeId()),
					ProjectName: rap.ProjectName,
					//CompanyName: utils.ToEmptyString(in.CompanyName),
					AdminUserID:  in.UserID,
					AreaCount:    int64(len(rap.Areas)),
					UserCount:    1,
					IsSysCreated: rap.IsSysCreated,
					Desc:         "自动创建",
				}
				pis = append(pis, &po)
				dps = append(dps, &relationDB.SysDataProject{
					ProjectID:  int64(po.ProjectID),
					TargetType: def.TargetUser,
					TargetID:   po.AdminUserID,
					AuthType:   def.AuthAdmin,
				})
				if rap.Areas != nil {
					for _, area := range rap.Areas {
						var areaID = svcCtx.AreaID.GetSnowflakeId()
						var areaIDPath string = cast.ToString(areaID) + "-"
						var areaNamePath = area.AreaName + "-"
						areaPo := relationDB.SysAreaInfo{
							AreaID:       dataType.AreaID(areaID),
							ParentAreaID: def.RootNode, //创建时必填
							ProjectID:    po.ProjectID, //创建时必填
							AreaIDPath:   areaIDPath,
							AreaNamePath: areaNamePath,
							AreaName:     area.AreaName,
							AreaImg:      area.AreaImg,
							IsLeaf:       def.True,
							IsSysCreated: area.IsSysCreated,
						}
						ais = append(ais, &areaPo)
					}
				}
			}
			// 三次重试创建项目
			for i := 3; i > 0; i-- {
				err := stores.GetTenantConn(ctx).Transaction(func(tx *gorm.DB) error {
					if len(pis) > 0 {
						piDb := relationDB.NewProjectInfoRepo(tx)
						err := piDb.MultiInsert(ctx, pis)
						if err != nil {
							logx.WithContext(ctx).Errorf("批量插入项目信息失败: UserID=%d, error=%v", in.UserID, err)
							return err
						}
					}
					if len(dps) > 0 {
						err := relationDB.NewDataProjectRepo(tx).MultiInsert(ctx, dps)
						if err != nil {
							logx.WithContext(ctx).Errorf("批量插入数据项目失败: UserID=%d, error=%v", in.UserID, err)
							return err
						}
					}
					if len(ais) > 0 {
						aiRepo := relationDB.NewAreaInfoRepo(tx)
						err := aiRepo.MultiInsert(ctx, ais)
						if err != nil {
							logx.WithContext(ctx).Errorf("批量插入区域信息失败: UserID=%d, error=%v", in.UserID, err)
							return err
						}
					}
					return nil
				})
				if err == nil {
					logx.WithContext(ctx).Infof("用户自动创建项目成功: UserID=%d", in.UserID)
					return
				}
			}
			logx.WithContext(ctx).Errorf("用户自动创建项目最终失败: UserID=%d", in.UserID)

		})
	}
	return err
}

func Init() {
	ctx := ctxs.WithRoot(context.Background())
	utils.Go(ctx, func() {
		t := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-t.C:
				tenantCodes, err := caches.GetTenantCodes(ctx)
				if err != nil {
					logx.WithContext(ctx).Error(err)
					continue
				}
				for _, code := range tenantCodes {
					err = stores.WithNoDebug(ctx, relationDB.NewTenantInfoRepo).UpdateUserCount(ctx, code)
					if err != nil {
						logx.WithContext(ctx).Error(err)
					}
				}
			}
		}
	})
}
