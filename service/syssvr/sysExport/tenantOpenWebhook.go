package sysExport

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/client/tenantmanage"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/eventBus"
	"github.com/maypok86/otter"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

const (
	CodeDmDeviceConn             = "dmDeviceConn"
	CodeDmDeviceDisConn          = "dmDeviceDisConn"
	CodeDmDevicePropertyReport   = "devicePropertyReport"
	CodeDmDevicePropertyReportV2 = "devicePropertyReportV2"
	CodeDmDeviceEventReport      = "deviceEventReport"
)

type Webhook struct {
	*caches.Cache[sys.TenantOpenWebHook, string]
	cache otter.Cache[string, struct{}]
}

func NewTenantOpenWebhook(pm tenantmanage.TenantManage, fastEvent *eventBus.FastEvent) (*Webhook, error) {
	c, err := NewTenantOpenWebhookCache(pm, fastEvent)
	if err != nil {
		return nil, err
	}
	cc, err := otter.MustBuilder[string, struct{}](10_000).CollectStats().
		Cost(func(key string, value struct{}) uint32 {
			return 1
		}).
		WithTTL(15 * time.Minute).
		Build()
	return &Webhook{Cache: c, cache: cc}, nil
}

func (i *Webhook) Publish(ctx context.Context, code string, in any) error {
	key := GenWebhookCacheKey(ctx, code)
	if _, ok := i.cache.Get(key); ok {
		return nil
	}
	hook, err := i.GetData(ctx, key)
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			i.cache.Set(key, struct{}{}) //避免反复查询
			return nil
		}
		return err
	}
	req := gorequest.New().Retry(3, time.Second*2)
	url := hook.Hosts[0] + hook.Uri
	req.Post(url)
	for k, v := range hook.Handler {
		req.Set(k, v)
	}
	resp, body, errs := req.Type("json").Send(in).End()
	if errs != nil {
		return errors.System.AddDetail(errs)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.System.AddDetail(body)
	}
	return nil
}
