package svc

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/config"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/shared/clients"
	"gitee.com/i-Things/core/shared/conf"
	"gitee.com/i-Things/core/shared/ctxs"
	"sync"
)

type Clients struct {
	MiniProgram *clients.MiniProgram
	DingTalk    *clients.DingTalk
}
type ClientsManage struct {
	Config config.Config
}

var (
	tc = sync.Map{}
)

func NewClients(c config.Config) *ClientsManage {
	return &ClientsManage{Config: c}
}

func (c *ClientsManage) GetClients(ctx context.Context, tenantCode string) (Clients, error) {
	if tenantCode == "" {
		tenantCode = ctxs.GetUserCtx(ctx).TenantCode
	}
	val, ok := tc.Load(tenantCode)
	if ok {
		return val.(Clients), nil
	}
	//如果缓存里没有,需要查库
	cfg, err := relationDB.NewTenantConfigRepo(ctx).FindOneByFilter(ctx, relationDB.TenantConfigFilter{TenantCode: tenantCode})
	if err != nil {
		return Clients{}, err
	}
	var cli Clients
	if cfg.DingTalk != nil && cfg.DingTalk.AppSecret != "" {
		cli.DingTalk, err = clients.NewDingTalkClient(&conf.ThirdConf{
			AppID:     cfg.DingTalk.AppID,
			AppKey:    cfg.DingTalk.AppKey,
			AppSecret: cfg.DingTalk.AppSecret,
		})
		if err != nil {
			return Clients{}, err
		}
	}
	if cfg.WxMini != nil && cfg.WxMini.AppSecret != "" {
		cli.MiniProgram, _ = clients.NewWxMiniProgram(ctx, &conf.ThirdConf{
			AppID:     cfg.WxMini.AppID,
			AppKey:    cfg.WxMini.AppKey,
			AppSecret: cfg.WxMini.AppSecret,
		}, c.Config.CacheRedis)
	}
	tc.Store(tenantCode, cli)
	return cli, nil
}
