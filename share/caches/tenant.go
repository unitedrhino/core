package caches

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/core/share/domain/tenant"
	"gitee.com/unitedrhino/share/caches"
)

// 生产用户数据权限缓存key
func genTenantKey() string {
	return "tenant"
}

func InitTenant(ctx context.Context, tenants ...*tenant.Info) error {
	if len(tenants) == 0 {
		return nil
	}
	return caches.GetStore().HmsetCtx(ctx, genTenantKey(), DoToTenantMap(tenants...))
}

func SetTenant(ctx context.Context, t *tenant.Info) error {
	val, _ := json.Marshal(t)
	return caches.GetStore().HsetCtx(ctx, genTenantKey(), t.Code, string(val))
}

func DelTenant(ctx context.Context, code string) error {
	_, err := caches.GetStore().HdelCtx(ctx, genTenantKey(), code)
	return err
}

func GetTenant(ctx context.Context, code string) (*tenant.Info, error) {
	val, err := caches.GetStore().HgetCtx(ctx, genTenantKey(), code)
	if err != nil {
		return nil, err
	}
	var ret tenant.Info
	err = json.Unmarshal([]byte(val), &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetTenantCodes(ctx context.Context) ([]string, error) {
	return caches.GetStore().HkeysCtx(ctx, genTenantKey())
}

func DoToTenantMap(tenants ...*tenant.Info) map[string]string {
	var ret = map[string]string{}
	for _, v := range tenants {
		val, _ := json.Marshal(v)
		ret[v.Code] = string(val)
	}
	return ret
}
