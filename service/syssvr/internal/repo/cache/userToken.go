package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/core/share/users"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
)

type UserToken struct {
}

func NewUserToken() *UserToken {
	return &UserToken{}
}

func (u *UserToken) GenKey(claims users.LoginClaims) string {
	return fmt.Sprintf("userToken:%v", claims.UserID)
}

func (u *UserToken) GenKey2(claims users.LoginClaims) string {
	return fmt.Sprintf("userToken:%v", claims.AppCode)
}

func (u *UserToken) Login(ctx context.Context, claims users.LoginClaims) error {
	err := caches.GetStore().Hset(u.GenKey(claims), claims.AppCode, utils.MarshalNoErr(claims))
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}

func (u *UserToken) CheckToken(ctx context.Context, claims users.LoginClaims) error {
	tk, err := caches.GetStore().Hget(u.GenKey(claims), claims.AppCode)
	if err != nil {
		return errors.NotLogin
	}
	var tkStu users.LoginClaims
	err = json.Unmarshal([]byte(tk), &tkStu)
	if err != nil {
		return errors.NotLogin
	}
	if tkStu.ID != claims.ID { //其他账号登录了,该账号被踢出
		return errors.AccountKickedOut
	}
	return nil
}
