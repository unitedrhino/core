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

	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"

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
	ui := relationDB.SysUserInfo{
		UserID: l.svcCtx.UserID.GetSnowflakeId(),
	}
	ui.Password = utils.MakePwd(in.Password, ui.UserID, false)
	if in.Info != nil {
		ui.NickName = in.Info.NickName
		if in.Info.UserName == "" {
			ui.UserName = utils.AnyToNullString(in.Info.UserName)
		}
	}

	switch in.RegType {
	case users.RegEmail:
		email := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, def.CaptchaUseRegister, in.CodeID, in.Code)
		if email == "" || email != in.Account {
			return nil, errors.Captcha
		}
		ui.Email = utils.AnyToNullString(in.Account)
		ui.UserName = ui.Email
	case users.RegPhone:
		phone := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypePhone, def.CaptchaUseRegister, in.CodeID, in.Code)
		if phone == "" || phone != in.Account {
			return nil, errors.Captcha
		}
		ui.Phone = utils.AnyToNullString(in.Account)
		ui.UserName = ui.Phone
	}
	if in.Password != "" {
		err := CheckPwd(l.svcCtx, in.Password)
		if err != nil {
			return nil, err
		}
	}

	wxOpenCode := in.Expand["wxOpenCode"]
	if wxOpenCode != "" {
		at, err := GetWxLoginResAccessToken(l.ctx, wxOpenCode)
		if err != nil {
			cli, er := l.svcCtx.Cm.GetClients(l.ctx, "")
			if er != nil {
				return nil, errors.System.AddDetail(er)
			}
			if cli.WxOfficial == nil {
				return nil, errors.System.AddDetail(er)
			}
			at2, er := cli.WxOfficial.GetOauth().GetUserAccessToken(in.Code)
			if er != nil {
				return nil, errors.System.AddDetail(er)
			}
			at = &at2
		}
		StoreWxRegisterResAccessToken(l.ctx, wxOpenCode, at)
		_, err = relationDB.NewUserInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserInfoFilter{WechatUnionID: at.UnionID, WechatOpenID: at.OpenID})
		if err == nil {
			return nil, errors.BindAccount.WithMsg("微信已绑定其他账号")
		}
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
		ui.WechatOpenID = sql.NullString{Valid: true, String: at.OpenID}
		if at.UnionID != "" {
			ui.WechatUnionID = sql.NullString{Valid: true, String: at.UnionID}
		}
	}
	conn := stores.GetTenantConn(l.ctx)
	err := l.FillUserInfo(&ui, conn)
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
	cli, err := l.svcCtx.Cm.GetClients(l.ctx, "")
	if err != nil || cli.MiniProgram == nil {
		return nil, errors.System.AddDetail(err)
	}
	auth := cli.MiniProgram.GetAuth()

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
	conn := stores.GetTenantConn(l.ctx)
	err = conn.Transaction(func(tx *gorm.DB) error {
		uidb := relationDB.NewUserInfoRepo(tx)
		_, err = uidb.FindOneByFilter(l.ctx,
			relationDB.UserInfoFilter{WechatUnionID: wxUid.UnionID, WechatOpenID: wxUid.OpenID})
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
				UserID:       userID,
				Phone:        sql.NullString{Valid: true, String: wxPhone.PhoneInfo.PurePhoneNumber},
				UserName:     sql.NullString{Valid: true, String: wxPhone.PhoneInfo.PurePhoneNumber},
				WechatOpenID: sql.NullString{Valid: true, String: wxUid.OpenID},
			}
		} else if !(ui.WechatUnionID.Valid || ui.WechatOpenID.Valid) {
			userID = l.svcCtx.UserID.GetSnowflakeId()
			ui = &relationDB.SysUserInfo{
				UserID:       userID,
				Phone:        sql.NullString{Valid: true, String: wxPhone.PhoneInfo.PurePhoneNumber},
				UserName:     sql.NullString{Valid: true, String: wxPhone.PhoneInfo.PurePhoneNumber},
				WechatOpenID: sql.NullString{Valid: true, String: wxUid.OpenID},
			}
		}
		userID = ui.UserID
		phone, err := uidb.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phone: wxPhone.PhoneInfo.PurePhoneNumber})
		if err == nil { //手机号注册过,绑定微信信息
			phone.WechatOpenID = sql.NullString{Valid: true, String: wxUid.OpenID}
			if wxUid.UnionID != "" {
				phone.WechatUnionID = sql.NullString{Valid: true, String: wxUid.UnionID}
			}
			return uidb.Update(l.ctx, phone)
		} else if !errors.Cmp(err, errors.NotFind) { //如果是数据库错误,则直接返回
			return err
		}

		if wxUid.UnionID != "" {
			ui.WechatUnionID = sql.NullString{Valid: true, String: wxUid.UnionID}
		}
		if in.Info != nil {
			ui.NickName = in.Info.NickName
			if in.Info.UserName != "" {
				ui.UserName = utils.AnyToNullString(in.Info.UserName)
			}
		}
		err = l.FillUserInfo(ui, tx)
		return err
	})

	return &sys.UserRegisterResp{UserID: userID}, err
}

func (l *UserRegisterLogic) handleDingApp(in *sys.UserRegisterReq) (*sys.UserRegisterResp, error) {
	cli, err := l.svcCtx.Cm.GetClients(l.ctx, "")
	if err != nil || cli.DingMini == nil {
		return nil, errors.System.AddDetail(err)
	}
	ret, err := cli.DingMini.GetUserInfoByCode(in.Code)
	if err != nil {
		l.Errorf("%v.Code2SessionContext err:%v", err)
		if ret.Code != 0 {
			return nil, errors.System.AddDetail(ret.Msg)
		}
		return nil, errors.System.AddDetail(err)
	} else if ret.Code != 0 {
		return nil, errors.Parameter.AddDetail(ret.Msg)
	}
	userID := l.svcCtx.UserID.GetSnowflakeId()
	ui := relationDB.SysUserInfo{
		UserID:         userID,
		DingTalkUserID: sql.NullString{Valid: true, String: ret.UserInfo.UserId},
		NickName:       ret.UserInfo.Name,
	}
	if ret.UserInfo.UnionId != "" {
		ui.DingTalkUnionID = sql.NullString{Valid: true, String: ret.UserInfo.UnionId}
	}
	if in.Info != nil {
		ui.NickName = in.Info.NickName
		if in.Info.UserName != "" {
			ui.UserName = utils.AnyToNullString(in.Info.UserName)
		}
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		uidb := relationDB.NewUserInfoRepo(tx)
		_, err = uidb.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{DingTalkUserID: ret.UserInfo.UserId, DingTalkUnionID: ret.UserInfo.UnionId})
		if err == nil { //已经注册过
			return errors.DuplicateRegister
		}
		if !errors.Cmp(err, errors.NotFind) {
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
		return err
	}
	err = tx.Transaction(func(tx *gorm.DB) error {
		uidb := relationDB.NewUserInfoRepo(tx)
		in.RegIP = uc.IP
		in.Role = cfg.RegisterRoleID
		err = uidb.Insert(ctx, in)
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
