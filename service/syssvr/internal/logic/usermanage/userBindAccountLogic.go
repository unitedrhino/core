package usermanagelogic

import (
	"context"
	"database/sql"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/users"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserBindAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserBindAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserBindAccountLogic {
	return &UserBindAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserBindAccountLogic) UserBindAccount(in *sys.UserBindAccountReq) (*sys.Empty, error) {
	cli, er := l.svcCtx.Cm.GetClients(l.ctx, "")
	if er != nil {
		return nil, errors.System.AddDetail(er)
	}
	if !utils.SliceIn(in.Type, cli.Config.LoginTypes...) {
		l.Errorf("不支持的登录方式:%v", in.Type)
		return nil, errors.NotSupportLogin
	}
	uc := ctxs.GetUserCtx(l.ctx)
	ui, err := relationDB.NewUserInfoRepo(l.ctx).FindOne(l.ctx, uc.UserID)
	if err != nil {
		return nil, err
	}
	switch in.Type {
	case users.RegEmail:
		if ui.Email.Valid {
			return &sys.Empty{}, errors.BindAccount
		}
		email := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypeEmail, def.CaptchaUseBindAccount, in.CodeID, in.Code)
		if email == "" || email != in.Account {
			return nil, errors.Captcha
		}
		ui.Email = utils.AnyToNullString(in.Account)
	case users.RegPhone:
		phone := l.svcCtx.Captcha.Verify(l.ctx, def.CaptchaTypePhone, def.CaptchaUseBindAccount, in.CodeID, in.Code)
		if phone == "" || phone != in.Account {
			return nil, errors.Captcha
		}
		ui.Phone = utils.AnyToNullString(in.Account)
	case users.RegWxMiniP:
		if ui.WechatUnionID.Valid || ui.WechatOpenID.Valid {
			return &sys.Empty{}, errors.BindAccount
		}
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
		if ret.UnionID != "" {
			ui.WechatUnionID = sql.NullString{ret.UnionID, true}
		}
		ui.WechatOpenID = sql.NullString{ret.OpenID, true}
	case users.RegWxOpen:
		if ui.WechatUnionID.Valid || ui.WechatOpenID.Valid {
			return &sys.Empty{}, errors.BindAccount
		}
		if cli.WxOfficial == nil {
			return nil, errors.System.AddDetail(er)
		}
		at, er := cli.WxOfficial.GetOauth().GetUserAccessToken(in.Code)
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		if at.UnionID != "" {
			ui.WechatUnionID = sql.NullString{at.UnionID, true}
		}
		ui.WechatOpenID = sql.NullString{at.OpenID, true}

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
		ui.DingTalkUserID = sql.NullString{ret.UserInfo.UserId, true}
	}
	err = relationDB.NewUserInfoRepo(l.ctx).Update(l.ctx, ui)
	return &sys.Empty{}, err
}
