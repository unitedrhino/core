package self

import (
	"context"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/role"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/tools"
	"gitee.com/unitedrhino/share/utils"
	"github.com/mssola/user_agent"

	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.UserLoginReq) (resp *types.UserLoginResp, err error) {
	uc := ctxs.GetUserCtx(l.ctx)
	ua := user_agent.New(ctxs.GetUserCtx(l.ctx).Os)
	browser, _ := ua.Browser()
	os := ua.OS()

	uResp, err := l.svcCtx.UserRpc.UserLogin(l.ctx, &sys.UserLoginReq{
		Account:   req.Account,
		PwdType:   req.PwdType,
		Password:  req.Password,
		LoginType: req.LoginType,
		Code:      req.Code,
		CodeID:    req.CodeID,
		Ip:        ctxs.GetUserCtx(l.ctx).IP,
	})

	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.Login req=%v err=%+v", utils.FuncName(), utils.Fmt(req), er)
		//登录失败记录
		l.svcCtx.LogRpc.LoginLogCreate(l.ctx, &sys.LoginLogCreateReq{
			UserID:        0,
			AppCode:       uc.AppCode,
			UserName:      req.Account, //todo 这里也要调整
			IpAddr:        ctxs.GetUserCtx(l.ctx).IP,
			LoginLocation: tools.GetCityByIp(ctxs.GetUserCtx(l.ctx).IP),
			Browser:       browser,
			Os:            os,
			Msg:           er.Error(),
			Code:          er.Code,
		})
		return nil, er
	}
	//登录成功记录

	_, err = l.svcCtx.LogRpc.LoginLogCreate(l.ctx, &sys.LoginLogCreateReq{
		AppCode:       uc.AppCode,
		UserID:        uResp.Info.UserID,
		UserName:      uResp.Info.UserName,
		IpAddr:        ctxs.GetUserCtx(l.ctx).IP,
		LoginLocation: tools.GetCityByIp(ctxs.GetUserCtx(l.ctx).IP),
		Browser:       browser,
		Os:            os,
		Msg:           "登录成功",
		Code:          errors.OK.GetCode(),
	})
	info, err := l.svcCtx.UserRpc.UserRoleIndex(l.ctx, &sys.UserRoleIndexReq{
		UserID: uResp.Info.UserID,
	})
	if err != nil {
		return nil, err
	}
	var (
		roles []*types.RoleInfo
	)

	roles = role.ToRoleInfosTypes(info.List)
	return &types.UserLoginResp{
		Info: types.UserInfo{
			UserID:      uResp.Info.UserID,
			UserName:    uResp.Info.UserName,
			Password:    "",
			Email:       utils.ToNullString(uResp.Info.Email),
			Phone:       utils.ToNullString(uResp.Info.Phone),
			LastIP:      uResp.Info.LastIP,
			RegIP:       uResp.Info.RegIP,
			NickName:    uResp.Info.NickName,
			City:        uResp.Info.City,
			Country:     uResp.Info.Country,
			Province:    uResp.Info.Province,
			Language:    uResp.Info.Language,
			HeadImg:     uResp.Info.HeadImg,
			CreatedTime: uResp.Info.CreatedTime,
			Role:        uResp.Info.Role,
			Sex:         uResp.Info.Sex,
			IsAllData:   uResp.Info.IsAllData,
		},
		Roles: roles,
		Token: types.JwtToken{
			AccessToken:  uResp.Token.AccessToken,
			AccessExpire: uResp.Token.AccessExpire,
			RefreshAfter: uResp.Token.RefreshAfter,
		},
	}, nil
}
