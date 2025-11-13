package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"gitee.com/unitedrhino/core/share/users"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/syncx"
)

type UserToken struct {
	sf syncx.SingleFlight
}

func NewUserToken() *UserToken {
	return &UserToken{sf: syncx.NewSingleFlight()}
}

func (u *UserToken) GenKey(userID int64) string {
	return fmt.Sprintf("userToken:%v", userID)
}

func (u *UserToken) GenField(appCode, id string) string {
	return fmt.Sprintf("%s:%v", appCode, id)
}
func (u *UserToken) GenKey3(userID int64) string {
	return fmt.Sprintf("userToken:%v", userID)
}

func (u *UserToken) ParseField(field string) (appCode string, id string) {
	appCode, id, _ = strings.Cut(field, ":")
	return
}

func (u *UserToken) Login(ctx context.Context, claims users.LoginClaims, accessExpire int64, isSsl bool) error {
	ret, err := caches.GetStore().Hgetall(u.GenKey(claims.UserID))
	if err != nil {
		return errors.System.AddDetail(err)
	}
	if isSsl {
		for k, v := range ret {
			appCode, _ := u.ParseField(k)
			if appCode == claims.AppCode {
				caches.GetStore().Hdel(u.GenKey(claims.UserID), k)
				continue
			}
			//如果有过期的清除过期的
			var tkStu users.LoginClaims
			err = json.Unmarshal([]byte(v), &tkStu)
			if tkStu.ExpiresAt.Before(time.Now()) {
				caches.GetStore().Hdel(u.GenKey(claims.UserID), k)
			}
		}
	}
	if len(ret) > 5 { //一个人最多只能存在5个token
		var cs []users.LoginClaims
		for k, v := range ret {
			var c users.LoginClaims
			err = json.Unmarshal([]byte(v), &c)
			if err != nil {
				caches.GetStore().Hdel(u.GenKey(claims.UserID), k)
				continue
			}
			cs = append(cs, c)
		}
		sort.Slice(cs, func(i, j int) bool {
			a := cs[i].ExpiresAt
			b := cs[j].ExpiresAt
			// 规则2：两者都有过期时间 → 比较时间大小（早的在前）
			return a.Time.After(b.Time)
		})
		if len(cs) > 5 {
			for _, v := range cs[5:] {
				caches.GetStore().Hdel(u.GenKey(claims.UserID), u.GenField(v.AppCode, v.ID))
			}
		}
	}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(accessExpire) * time.Second))
	err = caches.GetStore().Hset(u.GenKey(claims.UserID), u.GenField(claims.AppCode, claims.ID), utils.MarshalNoErr(claims))
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}

func (u *UserToken) CheckToken(ctx context.Context, claims users.LoginClaims, accessExpire int64) error {
	_, err := u.sf.Do(fmt.Sprintf("%s:%s:%s", claims.UserID, claims.AppCode, claims.ID), func() (any, error) {
		tk, err := caches.GetStore().Hget(u.GenKey(claims.UserID), u.GenField(claims.AppCode, claims.ID))
		if err != nil || tk == "" {
			return nil, errors.NotLogin
		}
		var tkStu users.LoginClaims
		err = json.Unmarshal([]byte(tk), &tkStu)
		if err != nil {
			return nil, errors.NotLogin
		}
		if tkStu.ExpiresAt.Before(time.Now()) { //已经过期了
			return nil, errors.TokenExpired.WithMsg("登录过期,请退出重新登录")
		}
		if (tkStu.ExpiresAt.Unix()-time.Now().Unix())*2 < accessExpire { //需要刷新token
			tkStu.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(accessExpire) * time.Second))
			tkStu.IssuedAt = jwt.NewNumericDate(time.Now())
			tkStu.Issuer = cast.ToString(cast.ToInt64(tkStu.Issuer) + 1)
			caches.GetStore().Hset(u.GenKey(claims.UserID), u.GenField(claims.AppCode, claims.ID), utils.MarshalNoErr(tkStu))
		}
		return nil, nil
	})
	return err
}

func (u *UserToken) KickedOut(ctx context.Context, userID int64) error {
	_, err := caches.GetStore().Del(u.GenKey3(userID))
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}
