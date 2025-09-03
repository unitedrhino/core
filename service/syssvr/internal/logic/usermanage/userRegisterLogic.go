package usermanagelogic

import (
	"context"
	"database/sql"

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

func (l *UserRegisterLogic) UserRegister(in *sys.UserRegisterReq) (*sys.UserRegisterResp, error) {
	switch in.RegType {
	case users.RegDingApp:
		return l.handleDingApp(in)
	case users.RegWxMiniP:
		return l.handleWxminip(in)
	case users.RegEmail, users.RegPhone:
		return l.handleEmailOrPhone(in)
	default:
		return nil, errors.NotRealize.AddMsgf(in.RegType)
	}
}

func (l *UserRegisterLogic) handleEmailOrPhone(in *sys.UserRegisterReq) (*sys.UserRegisterResp, error) {
	thirdHandle := func(tx *stores.DB, ui *relationDB.SysUserInfo) error {
		wxOpenCode := in.Expand["wxOpenCode"]
		if wxOpenCode != "" && in.AppID != "" {
			at, err := GetWxLoginResAccessToken(l.ctx, in.AppID, wxOpenCode)
			if err != nil {
				cli, er := l.svcCtx.ThirdClientsManage.GetWxOpenClient(l.ctx, l.uc.AppCode, in.AppID)
				if er != nil {
					return er
				}
				at2, er := cli.GetOauth().GetUserAccessToken(in.Code)
				if er != nil {
					return errors.System.AddDetail(er)
				}
				at = &at2
			}
			StoreWxRegisterResAccessToken(l.ctx, in.AppID, wxOpenCode, at)
			_, err = relationDB.NewUserThirdRepo(tx).FindOneByFilter(l.ctx, relationDB.UserThirdFilter{AppID: in.AppID, UnionID: at.UnionID, OpenID: at.OpenID})
			if err == nil {
				return errors.BindAccount.WithMsg("微信已绑定其他账号")
			}
			if !errors.Cmp(err, errors.NotFind) {
				return err
			}
			err = relationDB.NewUserThirdRepo(tx).Insert(l.ctx, &relationDB.SysUserThird{
				TenantCode: "",
				AppType:    def.ThirdTypeWxOpen,
				AppID:      in.AppID,
				UserID:     ui.UserID,
				UnionID:    at.UnionID,
				OpenID:     at.OpenID,
			})
			return err
		}
		return nil
	}

	old, err := relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{WithTenant: true, Accounts: []string{in.Account}})
	if err == nil {
		if len(old.Tenants) != 0 {
			return nil, errors.DuplicateRegister
		}
		err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
			err = thirdHandle(tx, old)
			if err != nil {
				return err
			}
			err = Register(l.ctx, l.svcCtx, old, true, nil)
			return err
		})
		return &sys.UserRegisterResp{UserID: old.UserID}, err
	}
	ui := relationDB.SysUserInfo{}

	uiInit := func() {
		ui.UserID = l.svcCtx.UserID.GetSnowflakeId()
		ui.Password = utils.MakePwd(in.Password, ui.UserID, false)
		if in.Info != nil {
			ui.NickName = in.Info.NickName
			if in.Info.UserName == "" {
				ui.UserName = utils.AnyToNullString(in.Info.UserName)
			}
		}
	}

	switch in.RegType {
	case users.RegEmail:
		email := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, def.CaptchaUseRegister, in.CodeID, in.Code)
		if email == "" || email != in.Account {
			return nil, errors.Captcha
		}
		uiInit()
		ui.Email = utils.AnyToNullString(in.Account)
		ui.UserName = ui.Email
	case users.RegPhone:
		phone := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypePhone, def.CaptchaUseRegister, in.CodeID, in.Code)
		if phone == "" || phone != in.Account {
			return nil, errors.Captcha
		}
		uiInit()
		ui.Phone = utils.AnyToNullString(in.Account)
		ui.UserName = ui.Phone
	}
	if in.Password != "" {
		err := CheckPwd(l.svcCtx, in.Password)
		if err != nil {
			return nil, err
		}
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err = thirdHandle(tx, old)
		if err != nil {
			return err
		}
		err = Register(l.ctx, l.svcCtx, old, true, nil)
		return err
	})
	if err != nil {
		return nil, err
	}
	e := l.svcCtx.FastEvent.Publish(l.ctx, topics.CoreUserCreate, def.IDs{IDs: []int64{ui.UserID}})
	if e != nil {
		l.Errorf("Publish CoreUserCreate %v err:%v", ui, e)
	}
	return &sys.UserRegisterResp{
		UserID: ui.UserID,
	}, nil
}

func (l *UserRegisterLogic) handleWxminip(in *sys.UserRegisterReq) (*sys.UserRegisterResp, error) {
	cli, err := l.svcCtx.ThirdClientsManage.GetWxMiniClient(l.ctx, l.uc.AppCode, in.AppID)
	if err != nil {
		return nil, err
	}
	auth := cli.GetAuth()
	wxUid, err := auth.Code2SessionContext(l.ctx, in.Code)
	if err != nil {
		l.Errorf("%v.Code2SessionContext err:%v", err)
		if wxUid.ErrCode != 0 {
			return nil, errors.System.AddDetail(wxUid.ErrMsg)
		}
		return nil, errors.System.AddDetail(err)
	} else if wxUid.ErrCode != 0 {
		return nil, errors.Parameter.AddDetail(wxUid.ErrMsg)
	}
	if in.Expand == nil || in.Expand["phoneCode"] == "" {
		return nil, errors.Parameter.AddMsg("微信小程序注册需要填写expand.phoneCode")
	}
	phoneCode := in.Expand["phoneCode"]
	wxPhone, err := auth.GetPhoneNumber(phoneCode)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	if wxPhone.ErrCode != 0 {
		return nil, errors.Parameter.AddDetail(wxPhone.ErrMsg)
	}
	var userID int64
	var isCreateUi bool
	conn := stores.GetTenantConn(l.ctx)
	err = conn.Transaction(func(tx *gorm.DB) error {
		uidb := relationDB.NewUserInfoRepo(tx)
		_, err := relationDB.NewUserThirdRepo(tx).FindOneByFilter(l.ctx,
			relationDB.UserThirdFilter{AppType: def.ThirdTypeWxMiniP, AppID: in.AppID, UnionID: wxUid.UnionID, OpenID: wxUid.OpenID})
		if err == nil { //已经注册过
			return errors.DuplicateRegister
		}
		if !errors.Cmp(err, errors.NotFind) {
			return err
		}
		ui, err := uidb.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phones: []string{wxPhone.PhoneInfo.PurePhoneNumber}})
		if err != nil {
			if !errors.Cmp(err, errors.NotFind) {
				return err
			}
			userID = l.svcCtx.UserID.GetSnowflakeId()
			ui = &relationDB.SysUserInfo{
				UserID:   userID,
				Phone:    sql.NullString{Valid: true, String: wxPhone.PhoneInfo.PurePhoneNumber},
				UserName: sql.NullString{Valid: true, String: wxPhone.PhoneInfo.PurePhoneNumber},
			}
			if in.Info != nil {
				ui.NickName = in.Info.NickName
				if in.Info.UserName != "" {
					ui.UserName = utils.AnyToNullString(in.Info.UserName)
				}
			}
			isCreateUi = true
		}
		err = relationDB.NewUserThirdRepo(tx).Insert(l.ctx, &relationDB.SysUserThird{
			AppType: def.ThirdTypeWxMiniP,
			AppID:   in.AppID,
			UserID:  ui.UserID,
			UnionID: wxUid.UnionID,
			OpenID:  wxUid.OpenID,
		})
		if err != nil {
			return err
		}

		err = l.FillUserInfo(ui, isCreateUi, tx)
		return err
	})
	return &sys.UserRegisterResp{UserID: userID}, err
}

func (l *UserRegisterLogic) handleDingApp(in *sys.UserRegisterReq) (*sys.UserRegisterResp, error) {
	cli, err := l.svcCtx.ThirdClientsManage.GetDingAppClient(l.ctx, l.uc.AppCode, in.AppID)
	if err != nil {
		return nil, err
	}
	ret, err := cli.GetUserInfoByCode(in.Code)
	if err != nil {
		l.Errorf("%v.Code2SessionContext err:%v", err)
		if ret.Code != 0 {
			return nil, errors.System.AddDetail(ret.Msg)
		}
		return nil, errors.System.AddDetail(err)
	} else if ret.Code != 0 {
		return nil, errors.Parameter.AddDetail(ret.Msg)
	}

	var userID int64
	var isCreateUi bool
	conn := stores.GetTenantConn(l.ctx)
	err = conn.Transaction(func(tx *gorm.DB) error {
		_, err := relationDB.NewUserThirdRepo(tx).FindOneByFilter(l.ctx,
			relationDB.UserThirdFilter{AppType: def.ThirdTypeDingApp, AppID: in.AppID, UnionID: ret.UserInfo.UnionId, OpenID: ret.UserInfo.UserId})
		if err == nil { //已经注册过
			return errors.DuplicateRegister
		}
		if !errors.Cmp(err, errors.NotFind) {
			return err
		}
		userID = l.svcCtx.UserID.GetSnowflakeId()
		ui := &relationDB.SysUserInfo{
			UserID:   userID,
			NickName: ret.UserInfo.Name,
		}
		if in.Info != nil {
			ui.NickName = in.Info.NickName
			if in.Info.UserName != "" {
				ui.UserName = utils.AnyToNullString(in.Info.UserName)
			}
		}
		isCreateUi = true
		err = relationDB.NewUserThirdRepo(tx).Insert(l.ctx, &relationDB.SysUserThird{
			AppType: def.ThirdTypeDingApp,
			AppID:   in.AppID,
			UserID:  ui.UserID,
			UnionID: ret.UserInfo.UnionId,
			OpenID:  ret.UserInfo.UserId,
		})
		if err != nil {
			return err
		}
		err = l.FillUserInfo(ui, isCreateUi, tx)
		return err
	})
	return &sys.UserRegisterResp{UserID: userID}, err
}

func (l *UserRegisterLogic) FillUserInfo(in *relationDB.SysUserInfo, isCreateUi bool, tx *gorm.DB) error {
	err := Register(l.ctx, l.svcCtx, in, isCreateUi, tx)
	if err != nil && errors.Cmp(err, errors.Duplicate) { //已经注册过
		return errors.DuplicateRegister
	}
	return err
}

func Register(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.SysUserInfo, isCreateUi bool, tx *gorm.DB) error {
	ctx = ctxs.WithAdmin(ctx)
	uc := ctxs.GetUserCtx(ctx)
	cfg, err := svcCtx.TenantConfigCache.GetData(ctx, uc.TenantCode)
	if err != nil {
		return err
	}
	if tx == nil {
		tx = stores.GetTenantConn(ctx)
	}
	err = tx.Transaction(func(tx *gorm.DB) error {
		if isCreateUi {
			uidb := relationDB.NewUserInfoRepo(tx)
			in.RegIP = uc.IP
			err = uidb.Insert(ctx, in)
			if err != nil {
				return err
			}
		}
		err = relationDB.NewUserTenantRepo(tx).Insert(ctx, &relationDB.SysUserTenant{UserID: in.UserID, Status: def.True, Tags: make(map[string]string)})
		if err != nil {
			return err
		}
		err = relationDB.NewUserRoleRepo(tx).Insert(ctx, &relationDB.SysUserRole{
			UserID: in.UserID,
			RoleID: cfg.RegisterRoleID,
		})
		if err != nil {
			return err
		}
		return err
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
			for i := 3; i > 0; i-- { //三次重试
				err := stores.GetTenantConn(ctx).Transaction(func(tx *gorm.DB) error {
					if len(pis) > 0 {
						piDb := relationDB.NewProjectInfoRepo(tx)
						err := piDb.MultiInsert(ctx, pis)
						if err != nil {
							logx.WithContext(ctx).Error(err)
							return err
						}
					}
					if len(dps) > 0 {
						err := relationDB.NewDataProjectRepo(tx).MultiInsert(ctx, dps)
						if err != nil {
							logx.WithContext(ctx).Error(err)
							return err
						}
					}
					if len(ais) > 0 {
						aiRepo := relationDB.NewAreaInfoRepo(tx)
						err := aiRepo.MultiInsert(ctx, ais)
						if err != nil {
							logx.WithContext(ctx).Error(err)
							return err
						}
					}
					return nil
				})
				if err == nil {
					return
				}
			}

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
