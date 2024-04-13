package tenantOpenWebhook

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/client/tenantmanage"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/core/service/syssvr/sysExport"
	"gitee.com/i-Things/share/caches"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/eventBus"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

const (
	CodeDmDeviceConn           = "dmDeviceConn"
	CodeDmDeviceDisConn        = "dmDeviceDisConn"
	CodeDmDevicePropertyUpdate = "dmDevicePropertyUpdate"
)

type Info struct {
	*caches.Cache[sys.TenantOpenWebHook]
}

func NewTenantOpenWebhook(pm tenantmanage.TenantManage, fastEvent *eventBus.FastEvent) (*Info, error) {
	c, err := sysExport.NewTenantOpenWebhookCache(pm, fastEvent)
	if err != nil {
		return nil, err
	}
	return &Info{Cache: c}, nil
}

func (i *Info) Publish(ctx context.Context, code string, in any) error {
	hook, err := i.GetData(ctx, sysExport.GenWebhookCacheKey(ctx, code))
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
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
