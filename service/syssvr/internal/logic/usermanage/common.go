package usermanagelogic

import (
	"context"
	"database/sql"
	"strings"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/share/clients/oauth2"
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

// applyOAuthLoginAccount 为 OAuth 自动注册用户填充账号和展示名，避免 Apple 等平台不返回昵称时用户资料为空
func applyOAuthLoginAccount(ui *relationDB.SysUserInfo) {
	if ui == nil {
		return
	}
	if !ui.Email.Valid || ui.Email.String == "" {
		return
	}
	if !ui.UserName.Valid || ui.UserName.String == "" {
		ui.UserName = ui.Email
	}
	if ui.NickName == "" {
		ui.NickName = ui.Email.String
	}
}

// findOrBindAppleUser 按 Apple 用户 ID 查找用户；未找到时用已验证邮箱绑定同租户老用户。
func findOrBindAppleUser(ctx context.Context, repo *relationDB.UserInfoRepo, tenantCode string, aUser *oauth2.AppleUser) (*relationDB.SysUserInfo, error) {
	if aUser == nil || strings.TrimSpace(aUser.Sub) == "" {
		return nil, errors.Parameter.AddMsg("Apple用户标识为空")
	}
	appleUserID := strings.TrimSpace(aUser.Sub)
	uc, err := repo.FindOneByFilter(ctx, relationDB.UserInfoFilter{
		TenantCode:  tenantCode,
		AppleUserID: appleUserID,
	})
	if err == nil {
		return uc, nil
	}
	if !errors.Cmp(err, errors.NotFind) {
		return nil, err
	}
	email := strings.TrimSpace(aUser.Email)
	if email == "" || !strings.EqualFold(strings.TrimSpace(aUser.EmailVerified), "true") {
		return nil, errors.NotFind
	}
	uc, err = repo.FindOneByFilter(ctx, relationDB.UserInfoFilter{
		TenantCode: tenantCode,
		Emails:     []string{email},
	})
	if err != nil {
		return nil, err
	}
	if uc.AppleUserID.Valid && strings.TrimSpace(uc.AppleUserID.String) != "" {
		if uc.AppleUserID.String == appleUserID {
			return uc, nil
		}
		return nil, errors.BindAccount
	}
	err = repo.UpdateWithField(ctx, relationDB.UserInfoFilter{
		TenantCode: tenantCode,
		UserIDs:    []int64{uc.UserID},
	}, map[string]any{
		"apple_user_id": sql.NullString{Valid: true, String: appleUserID},
	})
	if err != nil {
		if errors.Cmp(err, errors.Duplicate) {
			return findAppleUserAfterDuplicate(ctx, repo, tenantCode, appleUserID, err)
		}
		return nil, err
	}
	uc.AppleUserID = sql.NullString{Valid: true, String: appleUserID}
	return uc, nil
}

// findAppleUserAfterDuplicate 在并发绑定撞唯一索引时重新读取已绑定的 Apple 用户。
func findAppleUserAfterDuplicate(ctx context.Context, repo *relationDB.UserInfoRepo, tenantCode string, appleUserID string, duplicateErr error) (*relationDB.SysUserInfo, error) {
	uc, err := repo.FindOneByFilter(ctx, relationDB.UserInfoFilter{
		TenantCode:  tenantCode,
		AppleUserID: appleUserID,
	})
	if err == nil {
		return uc, nil
	}
	return nil, duplicateErr
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
