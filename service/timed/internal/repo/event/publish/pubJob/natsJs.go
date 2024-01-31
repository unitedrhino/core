package pubJob

import (
	"context"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/events"
	"gitee.com/i-Things/share/utils"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	natsJsClient struct {
		client nats.JetStreamContext
	}
)

func newNatsJsClient(conf conf.NatsConf) (*natsJsClient, error) {
	js, err := clients.NewNatsJetStreamClient(conf)
	if err != nil {
		return nil, err
	}
	return &natsJsClient{client: js}, nil
}

func (n *natsJsClient) Publish(ctx context.Context, topic string, payload []byte) error {
	ret, err := n.client.Publish(topic, events.NewEventMsg(ctx, payload))
	if err != nil {
		logx.WithContext(ctx).Errorf("%s info:%v,ret:%v,err:%v", utils.FuncName(),
			string(payload), ret, err)
		return err
	}
	return nil
}
