package pubJob

import (
	"context"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/events"
	"gitee.com/unitedrhino/share/utils"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	natsClient struct {
		client *nats.Conn
	}
)

func newNatsClient(conf conf.NatsConf) (*natsClient, error) {
	js, err := clients.NewNatsClient(conf)
	if err != nil {
		return nil, err
	}
	return &natsClient{client: js}, nil
}

func (n *natsClient) Publish(ctx context.Context, topic string, payload []byte) error {
	err := n.client.Publish(topic, events.NewEventMsg(ctx, payload))
	if err != nil {
		logx.WithContext(ctx).Errorf("%s info:%v,err:%v", utils.FuncName(),
			string(payload), err)
		return err
	}
	return nil
}
