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
	Weather  Weather      `json:",optional"`
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
	GoogleJWKS struct {
		URL         string `json:",optional,env=GOOGLE_JWKS_URL"`                          //Google ID Token 验签用 JWKS 地址
		HeaderName  string `json:",default=X-YKHL-JWKS-TOKEN,env=GOOGLE_JWKS_HEADER_NAME"` //请求 JWKS 地址时携带的密钥请求头
		AccessToken string `json:",optional,env=GOOGLE_JWKS_ACCESS_TOKEN"`                 //请求 JWKS 地址时携带的共享密钥
	} `json:",optional"`
	ThirdJwtSecret string `json:",default=Xk9mQ2vR7pLw4nBz8jYf,env=THIRD_JWT_SECRET"` //第三方jwt加密登录的密钥
	Sms            conf.Sms
	//WrongPasswordCounter conf.WrongPasswordCounter `json:",optional"`

	CaptchaPhoneIpLimit      []conf.Limit `json:",optional"`
	CaptchaPhoneAccountLimit []conf.Limit `json:",optional"`
	CaptchaEmailIpLimit      []conf.Limit `json:",optional"`
	CaptchaEmailAccountLimit []conf.Limit `json:",optional"`

	CaptchaEmailGetLimit []conf.Limit `json:",optional"` //如果达到限制,需要输入验证码
	CaptchaPhoneGetLimit []conf.Limit `json:",optional"` //如果达到限制,需要输入验证码

	LoginPwdIpLimit      []conf.Limit `json:",optional"` //密码错误限制
	LoginPwdAccountLimit []conf.Limit `json:",optional"` //密码错误限制
	LoginPwdCaptchaLimit []conf.Limit `json:",optional"` //密码输入几次就要求输入验证码

	UniPush UniPushConf `json:",optional"`
}

// UniPushConf uni-push 2.0（方案 A：HTTP 调 uniCloud yk-push-send）
// 配置方式同 ThirdJwtSecret：优先 sys.yaml 的 Secret，可选 env=YK_PUSH_HTTP_SECRET 覆盖。
type UniPushConf struct {
	Enabled           bool     `json:",default=false,env=UNIPUSH_ENABLED"`
	AppId             string   `json:",default=__UNI__F82AD01"`
	HttpUrl           string   `json:",optional"`
	Secret            string   `json:",optional,env=YK_PUSH_HTTP_SECRET"`
	TimeoutMs         int      `json:",default=5000"`
	ForceNotification bool     `json:",default=true"`
	NotifyCodes       []string `json:",optional"`     // 白名单 notifyCode，空则用 ruleDeviceAlarm、ruleScene
	PushOnMessage     bool     `json:",default=true"` // 兼容：type=message 时是否旁路 uni-push（新配置请用 systemNotice）
}

type Weather struct {
	ApiKey  string `json:",optional"` //参考: https://dev.qweather.com/
	ApiHost string `json:",optional"` //参考: https://console.qweather.com/setting
}

var DefaultIpLimit = []conf.Limit{
	{Timeout: 5 * 60, TriggerTime: 30, ForbiddenTime: 5 * 60},                       //5分钟内错误30次,封禁5分钟
	{Timeout: 60 * 60, TriggerTime: 100, ForbiddenTime: 60 * 60 * 24},               //1个小时内错误100次,封禁一天
	{Timeout: 60 * 60 * 24 * 5, TriggerTime: 200, ForbiddenTime: 60 * 60 * 24 * 30}, //5天内错误200次,封禁30天
}

var DefaultAccountLimit = []conf.Limit{
	{Timeout: 5 * 60, TriggerTime: 8, ForbiddenTime: 5 * 60},                       //5分钟内错误3次,封禁5分钟
	{Timeout: 60 * 60, TriggerTime: 12, ForbiddenTime: 60 * 60 * 24},               //1个小时内错误10次,封禁一天
	{Timeout: 60 * 60 * 24 * 5, TriggerTime: 25, ForbiddenTime: 60 * 60 * 24 * 30}, //5天内错误20次,封禁30天
}

var DefaultCaptchaLimit = []conf.Limit{
	{Timeout: 5 * 60, TriggerTime: 3, ForbiddenTime: 5 * 60},                       //5分钟内错误3次,封禁5分钟
	{Timeout: 60 * 60, TriggerTime: 10, ForbiddenTime: 60 * 60 * 24},               //1个小时内错误10次,封禁一天
	{Timeout: 60 * 60 * 24 * 5, TriggerTime: 20, ForbiddenTime: 60 * 60 * 24 * 30}, //5天内错误20次,封禁30天
}
