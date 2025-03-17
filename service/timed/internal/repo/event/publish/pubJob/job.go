package pubJob

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
)

type (
	PubJob struct {
		natsJs *natsJsClient
		nats   *natsClient
	}
)

func NewPubJob(c conf.EventConf) (*PubJob, error) {
	switch c.Mode {
	case conf.EventModeNats:
		nats, err := newNatsClient(c.Nats)
		if err != nil {
			return nil, err
		}
		pj := PubJob{nats: nats}
		return &pj, nil
	case conf.EventModeNatsJs:
		natsJs, err := newNatsJsClient(c.Nats)
		if err != nil {
			return nil, err
		}
		pj := PubJob{natsJs: natsJs}
		return &pj, nil
	}
	return nil, nil

}
func (p *PubJob) Publish(ctx context.Context, pubType string, topic string, payload []byte) error {
	return p.nats.Publish(ctx, topic, payload)
	//if pubType == conf.EventModeNatsJs {
	//	return p.natsJs.Publish(ctx, topic, payload)
	//} else {
	//
	//}
	//return nil
}
