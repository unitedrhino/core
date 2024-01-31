package oss

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"sync"

	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/oss/common"
)

type Client struct {
	Handle
}

var (
	client  *Client
	newOnce sync.Once
)

func NewOssClient(c conf.OssConf) (cli *Client, err error) {
	newOnce.Do(func() {
		ossManager, er := newOssManager(c)
		if er != nil {
			err = errors.Parameter.AddMsgf("oss 初始化失败 err:%v", err)
			return
		}
		client = &Client{
			ossManager,
		}
	})

	return client, err
}

type OpOption func(*common.OptionKv)

func (c *Client) getDefaultOption(ctx context.Context) OpOption {
	return func(option *common.OptionKv) {
		option.SetHttpParams("x-process", "xxxxx")
	}
}
