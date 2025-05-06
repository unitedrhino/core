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
func (u *UserToken) GenKey3(userID int64) string {
	return fmt.Sprintf("userToken:%v", userID)
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

func (u *UserToken) CheckToken(ctx context.Context, claims users.LoginClaims, isSsl bool) error {
	tk, err := caches.GetStore().Hget(u.GenKey(claims), claims.AppCode)
	if err != nil {
		return errors.NotLogin
	}
	var tkStu users.LoginClaims
	err = json.Unmarshal([]byte(tk), &tkStu)
	if err != nil {
		return errors.NotLogin
	}
	if !isSsl {
		return nil
	}
	if tkStu.ID != claims.ID { //其他账号登录了,该账号被踢出
		return errors.AccountKickedOut
	}
	return nil
}
func (u *UserToken) KickedOut(ctx context.Context, userID int64) error {
	_, err := caches.GetStore().Del(u.GenKey3(userID))
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}
