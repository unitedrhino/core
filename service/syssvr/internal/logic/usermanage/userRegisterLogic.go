package usermanagelogic

import (
	"context"
	"database/sql"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/users"
	"gitee.com/i-Things/share/utils"
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
	return &sys.UserRegisterResp{}, nil
}

func (l *UserRegisterLogic) handleEmailOrPhone(in *sys.UserRegisterReq) (*sys.UserRegisterResp, error) {
	ui := relationDB.SysUserInfo{
		UserID: l.svcCtx.UserID.GetSnowflakeId(),
	}
	ui.Password = utils.MakePwd(in.Password, ui.UserID, false)
	if in.Info != nil {
		ui.NickName = in.Info.NickName
		if in.Info.UserName != "" {
			ui.UserName = utils.AnyToNullString(in.Info.UserName)
		}
	}

	switch in.RegType {
	case users.RegEmail:
		//email := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, def.CaptchaUseRegister, in.CodeID, in.Code)
		//if email == "" || email != in.Account {
		//	return nil, errors.Captcha
		//}
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
	err := CheckPwd(l.svcCtx, in.Password)
	if err != nil {
		return nil, err
	}
	conn := stores.GetTenantConn(l.ctx)
	err = l.FillUserInfo(&ui, conn)
	if err != nil {
		return nil, err
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
		ui, err := uidb.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{Phone: wxPhone.PhoneInfo.PurePhoneNumber})
		if err != nil {
			if !errors.Cmp(err, errors.NotFind) {
				return err
			}
			userID = l.svcCtx.UserID.GetSnowflakeId()
			ui = &relationDB.SysUserInfo{
				UserID:       userID,
				Phone:        sql.NullString{Valid: true, String: wxPhone.PhoneInfo.PurePhoneNumber},
				WechatOpenID: sql.NullString{Valid: true, String: wxUid.OpenID},
			}
		} else if !(ui.WechatUnionID.Valid || ui.WechatOpenID.Valid) {
			userID = l.svcCtx.UserID.GetSnowflakeId()
			ui = &relationDB.SysUserInfo{
				UserID:       userID,
				Phone:        sql.NullString{Valid: true, String: wxPhone.PhoneInfo.PurePhoneNumber},
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
	if err != nil || cli.MiniDing == nil {
		return nil, errors.System.AddDetail(err)
	}
	ret, err := cli.MiniDing.GetUserInfoByCode(in.Code)
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
	if in.Info != nil {
		ui.NickName = in.Info.NickName
		if in.Info.UserName != "" {
			ui.UserName = utils.AnyToNullString(in.Info.UserName)
		}
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		uidb := relationDB.NewUserInfoRepo(tx)
		_, err = uidb.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{DingTalkUserID: ret.UserInfo.UserId})
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
	ctx = ctxs.WithRoot(ctx)
	uc := ctxs.GetUserCtx(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		uidb := relationDB.NewUserInfoRepo(tx)
		cfg, err := relationDB.NewTenantConfigRepo(tx).FindOne(ctx)
		if err != nil {
			return err
		}
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
		if len(cfg.RegisterAutoCreateProject) > 0 {
			piDb := relationDB.NewProjectInfoRepo(tx)
			for _, rap := range cfg.RegisterAutoCreateProject {
				po := &relationDB.SysProjectInfo{
					ProjectID:   stores.ProjectID(svcCtx.ProjectID.GetSnowflakeId()),
					ProjectName: rap.ProjectName,
					//CompanyName: utils.ToEmptyString(in.CompanyName),
					AdminUserID:  in.UserID,
					IsSysCreated: rap.IsSysCreated,
					Desc:         "自动创建",
				}
				err = piDb.Insert(ctx, po)
				if err != nil {
					logx.WithContext(ctx).Errorf("%s.Insert err=%+v", utils.FuncName(), err)
					return err
				}
				err = relationDB.NewDataProjectRepo(tx).Insert(ctx, &relationDB.SysDataProject{
					ProjectID:  int64(po.ProjectID),
					TargetType: def.TargetUser,
					TargetID:   po.AdminUserID,
					AuthType:   def.AuthAdmin,
				})
				if rap.Areas != nil {
					aiRepo := relationDB.NewAreaInfoRepo(tx)
					for _, area := range rap.Areas {
						var areaID = svcCtx.AreaID.GetSnowflakeId()
						var areaIDPath string = cast.ToString(areaID) + "-"
						var areaNamePath = area.AreaName + "-"
						areaPo := &relationDB.SysAreaInfo{
							AreaID:       stores.AreaID(areaID),
							ParentAreaID: 1,            //创建时必填
							ProjectID:    po.ProjectID, //创建时必填
							AreaIDPath:   areaIDPath,
							AreaNamePath: areaNamePath,
							AreaName:     area.AreaName,
							IsLeaf:       def.True,
							IsSysCreated: area.IsSysCreated,
						}
						err = aiRepo.Insert(ctx, areaPo)
						if err != nil {
							logx.WithContext(ctx).Errorf("%s.Insert err=%+v", utils.FuncName(), err)
							return err
						}
					}
				}
			}
		}
		return err
	})
	return err
}
