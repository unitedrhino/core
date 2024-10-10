package subscribe

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/errors"
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

func NewSubServer(c conf.EventConf, nodeID int64) (SubApp, error) {
	switch c.Mode {
	case conf.EventModeNats, conf.EventModeNatsJs:
		return newNatsClient(c, nodeID)
	}
	return nil, errors.Parameter.AddMsgf("mode:%v not support", c.Mode)

}
