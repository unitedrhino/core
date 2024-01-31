package subscribe

import (
	"context"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/errors"
)

type (
	SubApp interface {
		Subscribe(handle Handle) error
	}
	Handle      func(ctx context.Context) ServerEvent
	ServerEvent interface {
		DataClean() error
	}
)

func NewSubServer(c conf.EventConf) (SubApp, error) {
	switch c.Mode {
	case conf.EventModeNats, conf.EventModeNatsJs:
		return newNatsClient(c)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)

}
