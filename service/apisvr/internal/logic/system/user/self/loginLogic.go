package self

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/apisvr/internal/logic/system/role"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/core/shared/ctxs"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/gogf/gf/v2/encoding/gcharset"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/mssola/user_agent"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

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

// GetCityByIp 获取ip所属城市
func GetCityByIp(ip string) string {
	if ip == "" {
		return ""
	}
	if ip == "[::1]" || ip == "127.0.0.1" {
		return "内网IP"
	}
	url := "http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip
	bytes := g.Client().GetBytes(context.TODO(), url)
	src := string(bytes)
	srcCharset := "GBK"
	tmp, _ := gcharset.ToUTF8(srcCharset, src)
	json, err := gjson.DecodeToJson(tmp)
	if err != nil {
		return ""
	}
	if json.Get("code").Int() == 0 {
		city := fmt.Sprintf("%s %s", json.Get("pro").String(), json.Get("city").String())
		return city
	} else {
		return ""
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
		l.Errorf("%s.rpc.Login req=%v err=%+v", utils.FuncName(), req, er)
		//登录失败记录
		l.svcCtx.LogRpc.LoginLogCreate(l.ctx, &sys.LoginLogCreateReq{
			UserID:        0,
			AppCode:       uc.AppCode,
			UserName:      req.Account, //todo 这里也要调整
			IpAddr:        ctxs.GetUserCtx(l.ctx).IP,
			LoginLocation: GetCityByIp(ctxs.GetUserCtx(l.ctx).IP),
			Browser:       browser,
			Os:            os,
			Msg:           er.Error(),
			Code:          400,
		})
		return nil, er
	}
	//登录成功记录
	l.svcCtx.LogRpc.LoginLogCreate(l.ctx, &sys.LoginLogCreateReq{
		AppCode:       uc.AppCode,
		UserID:        uResp.Info.UserID,
		UserName:      uResp.Info.UserName,
		IpAddr:        ctxs.GetUserCtx(l.ctx).IP,
		LoginLocation: GetCityByIp(ctxs.GetUserCtx(l.ctx).IP),
		Browser:       browser,
		Os:            os,
		Msg:           "登录成功",
		Code:          200,
	})
	info, err := l.svcCtx.UserRpc.UserRoleIndex(l.ctx, &sys.UserRoleIndexReq{
		UserID: uc.UserID,
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
			Email:       uResp.Info.Email,
			Phone:       uResp.Info.Phone,
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
