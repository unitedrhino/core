package users

import (
	"gitee.com/i-Things/share/errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Custom claims structure
type LoginClaims struct {
	UserID     int64  `json:",string"`
	Account    string //账号
	RoleIDs    []int64
	TenantCode string `json:",string"`
	IsAdmin    int64
	IsAllData  int64
	jwt.StandardClaims
}

func GetLoginJwtToken(secretKey string, iat, seconds, userID int64, account string, tenantCode string, roleIDs []int64, isAllData int64, isAdmin int64) (string, error) {
	claims := LoginClaims{
		UserID:     userID,
		RoleIDs:    roleIDs,
		TenantCode: tenantCode,
		IsAdmin:    isAdmin,
		Account:    account,
		IsAllData:  isAllData,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: iat + seconds,
			IssuedAt:  iat,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// 更新token
func RefreshLoginToken(tokenString string, secretKey string, AccessExpire int64) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &LoginClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*LoginClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = AccessExpire
		return CreateToken(secretKey, *claims)
	}
	return "", errors.TokenInvalid
}
