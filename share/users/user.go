package users

import "gitee.com/unitedrhino/share/utils"

type RegType = string

// phone 手机号 wxOpen 微信开放平台 wxIn 微信内 wxMiniP 微信小程序 pwd 账号密码
const (
	RegEmail      RegType = "email"      //邮箱
	RegPhone RegType = "phone"      //手机号
	RegWxOpen RegType = "wxOpen"     //微信开放平台登录
	RegWxIn RegType = "wxIn"       //微信内登录
	RegWxMiniP RegType = "wxMiniP"    //微信小程序
	RegWxOfficial RegType = "wxOfficial" //微信公众号登录
	RegDingApp RegType = "dingApp"    //钉钉应用(包含小程序,h5等方式)
	RegPwd RegType = "pwd"        //账号密码注册
	RegGoogle RegType = "google"     //google
	RegGithub RegType = "github"     //github
	RegApple  RegType = "apple"      //苹果
	RegHuawei RegType = "huawei"     //华为
	RegJwt    RegType = "jwt"        //第三方jwt加密登录
)

type UserInfoType uint8

const (
	Uid        UserInfoType = iota //用户UID
	InviterUid                     //邀请人用户id
	UserName                       //用户登录名
	GroupId                        //用户组id
	Email                          //邮箱
	Phone                          //手机号
	Wechat                         //微信
	InfoMax                        //结束
	AuthId                         //权限id
)

type UserStatus = int64

const (
	NotRegisterStatus UserStatus = iota //未注册完成状态只注册了第一步
	NormalStatus                        //正常状态
)

func GetLoginNameType(userName string) UserInfoType {
	if utils.IsPhone(userName) == true {
		return Phone
	}
	return UserName
}
