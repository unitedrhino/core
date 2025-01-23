package userDataAuth

import "gitee.com/unitedrhino/share/def"

type Area struct {
	AreaID         int64
	AreaIDPath     string
	AuthType       def.AuthType
	IsAuthChildren int64 //是否同时授权子节点,默认为2
}
