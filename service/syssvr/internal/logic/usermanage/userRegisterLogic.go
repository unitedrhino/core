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
		email := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, def.CaptchaUseRegister, in.CodeID, in.Code)
		if email == "" || email != in.Account {
			return nil, errors.Captcha
		}
		ui.Email = utils.AnyToNullString(in.Account)
	case users.RegPhone:
		phone := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypePhone, def.CaptchaUseRegister, in.CodeID, in.Code)
		if phone == "" || phone != in.Account {
			return nil, errors.Captcha
		}
		ui.Phone = utils.AnyToNullString(in.Account)
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
	ret, err := auth.Code2SessionContext(l.ctx, in.Code)
	if err != nil {
		l.Errorf("%v.Code2SessionContext err:%v", err)
		if ret.ErrCode != 0 {
			return nil, errors.System.AddDetail(ret.ErrMsg)
		}
		return nil, errors.System.AddDetail(err)
	} else if ret.ErrCode != 0 {
		return nil, errors.Parameter.AddDetail(ret.ErrMsg)
	}
	userID := l.svcCtx.UserID.GetSnowflakeId()

	conn := stores.GetTenantConn(l.ctx)
	err = conn.Transaction(func(tx *gorm.DB) error {
		uidb := relationDB.NewUserInfoRepo(tx)
		_, err = uidb.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{WechatUnionID: ret.UnionID})
		if err == nil { //已经注册过
			return errors.DuplicateRegister
		}
		if !errors.Cmp(err, errors.NotFind) {
			return err
		}
		ui := relationDB.SysUserInfo{
			UserID:        userID,
			WechatUnionID: sql.NullString{Valid: true, String: ret.UnionID},
		}
		if in.Info != nil {
			ui.NickName = in.Info.NickName
			if in.Info.UserName != "" {
				ui.UserName = utils.AnyToNullString(in.Info.UserName)
			}
		}
		err = l.FillUserInfo(&ui, tx)
		return err
	})

	return &sys.UserRegisterResp{UserID: userID}, err
}

func (l *UserRegisterLogic) handleDingApp(in *sys.UserRegisterReq) (*sys.UserRegisterResp, error) {
	cli, err := l.svcCtx.Cm.GetClients(l.ctx, "")
	if err != nil || cli.DingTalk == nil {
		return nil, errors.System.AddDetail(err)
	}
	ret, err := cli.DingTalk.GetUserInfoByCode(in.Code)
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
		return Register(l.ctx, l.svcCtx, &ui, tx)
	})
	return &sys.UserRegisterResp{UserID: userID}, err
}

func (l *UserRegisterLogic) FillUserInfo(in *relationDB.SysUserInfo, tx *gorm.DB) error {
	err := tx.Transaction(func(tx *gorm.DB) error {
		cfg, err := relationDB.NewTenantConfigRepo(tx).FindOne(l.ctx)
		if err != nil {
			return err
		}
		in.RegIP = l.uc.IP
		in.Role = cfg.RegisterRoleID
		uidb := relationDB.NewUserInfoRepo(tx)
		err = uidb.Insert(l.ctx, in)
		if err != nil {
			return err
		}
		err = relationDB.NewUserRoleRepo(tx).Insert(l.ctx, &relationDB.SysUserRole{
			UserID: in.UserID,
			RoleID: cfg.RegisterRoleID,
		})
		return err
	})
	if err != nil && errors.Cmp(err, errors.Duplicate) { //已经注册过
		return errors.DuplicateRegister
	}
	return err
}

func Register(ctx context.Context, svcCtx *svc.ServiceContext, in *relationDB.SysUserInfo, tx *gorm.DB) error {
	err := tx.Transaction(func(tx *gorm.DB) error {
		uidb := relationDB.NewUserInfoRepo(tx)
		cfg, err := relationDB.NewTenantConfigRepo(tx).FindOne(ctx)
		if err != nil {
			return err
		}
		in.RegIP = ctxs.GetUserCtx(ctx).IP
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
		if cfg.RegisterCreateProject != def.True {
			return nil
		}
		po := &relationDB.SysProjectInfo{
			ProjectID:   stores.ProjectID(svcCtx.ProjectID.GetSnowflakeId()),
			ProjectName: in.UserName.String + "的小屋",
			//CompanyName: utils.ToEmptyString(in.CompanyName),
			AdminUserID: in.UserID,
			//Region:      utils.ToEmptyString(in.Region),
			//Address:     utils.ToEmptyString(in.Address),
		}
		err = relationDB.NewProjectInfoRepo(tx).Insert(ctx, po)
		return err
	})
	return err
}
