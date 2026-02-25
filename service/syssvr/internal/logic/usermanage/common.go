package usermanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/golang-jwt/jwt/v5"
	"regexp"
)

func checkUser(ctx context.Context, userID int64) (*relationDB.SysUserInfo, error) {
	po, err := relationDB.NewUserInfoRepo(ctx).FindOne(ctx, userID)
	if err == nil {
		return po, nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, nil
	}
	return nil, err
}

func CheckPwd(svcCtx *svc.ServiceContext, pwd string) error {
	if svcCtx.Config.UserOpt.NeedPassWord &&
		utils.CheckPasswordLever(pwd) < svcCtx.Config.UserOpt.PassLevel {
		return errors.PasswordLevel
	}
	return nil
}
func CheckUserName(userName string) error {
	if ret, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_-]{6,19}$", userName); !ret {
		return errors.UsernameFormatErr.AddDetail("账号必须以字母开头，且只能包含大小写字母和数字下划线和减号。 长度为6到20位之间,或等于邮箱手机号")
	}
	return nil
}

// 第三方jwt加密登录的claims
type ThirdJwtClaims struct {
	Account string `json:"account"`
	jwt.RegisteredClaims
}

// 解析第三方jwt,返回account
func ParseThirdJwt(tokenString string, secret string) (string, error) {
	if secret == "" {
		return "", errors.Parameter.AddMsg("未配置第三方jwt密钥")
	}
	var claims ThirdJwtClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", errors.Parameter.AddMsg("jwt校验失败").AddDetail(err)
	}
	if !token.Valid {
		return "", errors.Parameter.AddMsg("jwt无效")
	}
	if claims.Account == "" {
		return "", errors.Parameter.AddMsg("jwt中缺少account字段")
	}
	return claims.Account, nil
}
