package self

import (
	"context"

	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/role"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/user"
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

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.UserRegisterReq) (resp *types.UserLoginResp, err error) {
	uResp, err := l.svcCtx.UserRpc.UserRegister(l.ctx, &sys.UserRegisterReq{
		RegType:  req.RegType,
		Account:  req.Account,
		Code:     req.Code,
		CodeID:   req.CodeID,
		Password: req.Password,
		Expand:   req.Expand,
		Info:     user.UserInfoToRpc(req.Info),
	})
	if err != nil {
		return nil, err
	}
	if req.IsWithLogin {
		//登录成功记录
		uc := ctxs.GetUserCtx(l.ctx)
		ua := user_agent.New(ctxs.GetUserCtx(l.ctx).Os)
		browser, _ := ua.Browser()
		os := ua.OS()
		_, err = l.svcCtx.LogRpc.LoginLogCreate(l.ctx, &sys.LoginLogCreateReq{
			AppCode:       uc.AppCode,
			UserID:        uResp.Info.UserID,
			UserName:      uResp.Info.UserName,
			IpAddr:        ctxs.GetUserCtx(l.ctx).IP,
			LoginLocation: tools.GetCityByIp(ctxs.GetUserCtx(l.ctx).IP),
			Browser:       browser,
			Os:            os,
			Msg:           "注册自动登录",
			Code:          errors.OK.GetCode(),
		})
	}

	info, err := l.svcCtx.UserRpc.UserRoleIndex(l.ctx, &sys.UserRoleIndexReq{
		UserID: uResp.Info.UserID,
	})
	if err != nil {
		return nil, err
	}
	var (
		roles []*types.RoleInfo
	)
	if len(uResp.Info.Password) != 0 {
		uResp.Info.Password = "xxxx"
	}
	roles = role.ToRoleInfosTypes(info.List)
	return &types.UserLoginResp{
		Info: types.UserInfo{
			UserID:      uResp.Info.UserID,
			UserName:    uResp.Info.UserName,
			Password:    uResp.Info.Password,
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
		Token: utils.Copy2[types.JwtToken](uResp.Token),
	}, nil
}
