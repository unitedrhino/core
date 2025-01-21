package config

import (
	"gitee.com/unitedrhino/share/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	CaptchaLen int `json:",default=6"` //验证码长度
	zrpc.RpcServerConf
	Database   conf.Database
	CacheRedis cache.ClusterConf
	UserToken  struct {
		AccessSecret string
		AccessExpire int64
	}
	OssConf  conf.OssConf `json:",optional"`
	Event    conf.EventConf
	Register struct {
		NeedDetail   bool   `json:",default=true"` //注册的时候是否需要填写用户信息,账号密码
		SecondSecret string //第二步需要的token秘钥
		SecondExpire int64  //token过期时间 单位:秒
	} `json:",optional"`
	TimedJobRpc conf.RpcClientConf `json:",optional"`
	UserOpt     struct {
		NeedUserName bool  `json:",default=true"` //注册是否必须填写账号密码
		NeedPassWord bool  `json:",default=true"` //注册是否必须填写账号密码
		PassLevel    int32 `json:",default=2"`    //用户密码强度级别
	} // 用户登录注册选项
	Map struct {
		Mode         string `json:",default=gaode"`
		AccessKey    string
		AccessSecret string
	}
	Sms conf.Sms
	//WrongPasswordCounter conf.WrongPasswordCounter `json:",optional"`

	CaptchaPhoneIpLimit      []conf.Limit `json:",optional"`
	CaptchaPhoneAccountLimit []conf.Limit `json:",optional"`
	CaptchaEmailIpLimit      []conf.Limit `json:",optional"`
	CaptchaEmailAccountLimit []conf.Limit `json:",optional"`
	LoginPwdIpLimit          []conf.Limit `json:",optional"` //密码错误限制
	LoginPwdAccountLimit     []conf.Limit `json:",optional"` //密码错误限制
}

var DefaultIpLimit = []conf.Limit{
	{Timeout: 5 * 60, TriggerTime: 30, ForbiddenTime: 5 * 60},                       //5分钟内错误30次,封禁5分钟
	{Timeout: 60 * 60, TriggerTime: 100, ForbiddenTime: 60 * 60 * 24},               //1个小时内错误100次,封禁一天
	{Timeout: 60 * 60 * 24 * 5, TriggerTime: 200, ForbiddenTime: 60 * 60 * 24 * 30}, //5天内错误200次,封禁30天
}

var DefaultAccountLimit = []conf.Limit{
	{Timeout: 5 * 60, TriggerTime: 3, ForbiddenTime: 5 * 60},                       //5分钟内错误3次,封禁5分钟
	{Timeout: 60 * 60, TriggerTime: 10, ForbiddenTime: 60 * 60 * 24},               //1个小时内错误10次,封禁一天
	{Timeout: 60 * 60 * 24 * 5, TriggerTime: 20, ForbiddenTime: 60 * 60 * 24 * 30}, //5天内错误20次,封禁30天
}
