package usermanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/share/users"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"

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
		if in.AppID == "" {
			return nil, errors.Parameter.AddMsgf("appID is empty")
		}
		cli, err := l.svcCtx.ThirdClientsManage.GetWxMiniClient(l.ctx, uc.AppCode, in.AppID)
		if err != nil {
			return nil, err
		}
		auth := cli.GetAuth()
		ret, er := auth.Code2SessionContext(l.ctx, in.Code)
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		if ret.ErrCode != 0 {
			return nil, errors.Parameter.AddMsgf(ret.ErrMsg)
		}
		_, err = relationDB.NewUserThirdRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserThirdFilter{
			AppType:  def.ThirdTypeWx,
			AppID:    in.AppID,
			UserID:   uc.UserID,
			UnionID:  ret.UnionID,
			OpenID:   ret.OpenID,
			WithUser: false,
		})
		if err == nil { //绑定过了
			return &sys.Empty{}, errors.BindAccount
		}
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
		err = relationDB.NewUserThirdRepo(l.ctx).Insert(l.ctx, &relationDB.SysUserThird{
			AppType: def.ThirdTypeWx,
			AppID:   in.AppID,
			UserID:  uc.UserID,
			UnionID: ret.UnionID,
			OpenID:  ret.OpenID,
		})
		return &sys.Empty{}, err
	case users.RegWxOpen:
		if in.AppID == "" {
			return nil, errors.Parameter.AddMsgf("appID is empty")
		}
		cli, err := l.svcCtx.ThirdClientsManage.GetWxOpenClient(l.ctx, uc.AppCode, in.AppID)
		if err != nil {
			return nil, err
		}
		ret, er := cli.GetOauth().GetUserAccessToken(in.Code)
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		if ret.ErrCode != 0 {
			return nil, errors.Parameter.AddMsgf(ret.ErrMsg)
		}
		_, err = relationDB.NewUserThirdRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserThirdFilter{
			AppType:  def.ThirdTypeWx,
			AppID:    in.AppID,
			UserID:   uc.UserID,
			UnionID:  ret.UnionID,
			OpenID:   ret.OpenID,
			WithUser: false,
		})
		if err == nil { //绑定过了
			return &sys.Empty{}, errors.BindAccount
		}
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
		err = relationDB.NewUserThirdRepo(l.ctx).Insert(l.ctx, &relationDB.SysUserThird{
			AppType: def.ThirdTypeWx,
			AppID:   in.AppID,
			UserID:  uc.UserID,
			UnionID: ret.UnionID,
			OpenID:  ret.OpenID,
		})
		return &sys.Empty{}, err

	case users.RegDingApp:
		if in.AppID == "" {
			return nil, errors.Parameter.AddMsgf("appID is empty")
		}
		cli, err := l.svcCtx.ThirdClientsManage.GetDingAppClient(l.ctx, uc.AppCode, in.AppID)
		if err != nil {
			return nil, err
		}
		ret, er := cli.GetUserInfoByCode(in.Code)
		if er != nil {
			return nil, errors.System.AddDetail(er)
		}
		if ret.Code != 0 {
			return nil, errors.Parameter.AddMsgf(ret.Msg)
		}
		_, err = relationDB.NewUserThirdRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserThirdFilter{
			AppType:  def.ThirdTypeDingApp,
			AppID:    in.AppID,
			UserID:   uc.UserID,
			UnionID:  ret.UserInfo.UnionId,
			OpenID:   ret.UserInfo.UserId,
			WithUser: false,
		})
		if err == nil { //绑定过了
			return &sys.Empty{}, errors.BindAccount
		}
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
		err = relationDB.NewUserThirdRepo(l.ctx).Insert(l.ctx, &relationDB.SysUserThird{
			AppType: def.ThirdTypeDingApp,
			AppID:   in.AppID,
			UserID:  uc.UserID,
			UnionID: ret.UserInfo.UnionId,
			OpenID:  ret.UserInfo.UserId,
		})
		return &sys.Empty{}, err
	}
	err = relationDB.NewUserInfoRepo(l.ctx).Update(l.ctx, ui)
	return &sys.Empty{}, err
}
