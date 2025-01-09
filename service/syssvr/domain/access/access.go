package access

import "gitee.com/unitedrhino/core/service/syssvr/domain/log"

type AuthType = int64

const (
	AuthTypeAll    = 1
	AuthTypeAdmin  = 2
	AuthTypeSupper = 3
)

type Access struct {
	Access map[string]*AccessInfo //授权组
}

type AccessInfo struct {
	Name       string    `json:"name"`       // 请求名称
	Code       string    `json:"code"`       // 请求名称
	Group      string    `json:"group"`      // 接口组
	IsNeedAuth int64     `json:"isNeedAuth"` // 是否需要认证（1是 2否）
	AuthType   string    `json:"authType"`   // 1(all) 全部人可以操作 2(admin) 默认授予租户管理员权限 3(superAdmin,supper) default租户才可以操作(超管是跨租户的)
	Desc       string    `json:"desc"`       // 备注
	Apis       []ApiInfo `json:"apis"`       //授权组下的接口
}

type ApiInfo struct {
	AccessCode    string `json:"accessCode"`                         // 范围编码
	Method        string `json:"method"`                             // 请求方式（1 GET 2 POST 3 HEAD 4 OPTIONS 5 PUT 6 DELETE 7 TRACE 8 CONNECT 9 其它）
	Route         string `json:"route"`                              // 路由
	Name          string `json:"name"`                               // 请求名称
	BusinessType  string `json:"businessType"`                       // 业务类型（1(add)新增 2修改(modify) 3删除(delete) 4查询(find) 5其它(other)
	RecordLogMode int64  `json:"recordLogMode,optional,range=[0:3]"` //   1为自动模式(读取类型忽略,其他类型记录日志) 2全部记录 3不记录
	Desc          string `json:"desc"`                               // 备注
}

func GetAuthType(authType string) AuthType {
	switch authType {
	case "admin":
		return 2
	case "supper", "supperAdmin":
		return 3
	default:
		return 1
	}
}

func (a ApiInfo) GetBusinessType() int64 {
	switch a.BusinessType {
	case "add":
		return log.OptAdd
	case "modify":
		return log.OptModify
	case "delete":
		return log.OptDel
	case "find":
		return log.OptQuery
	default:
		return log.OptOther
	}
}
