package svc

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/config"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/clients/dingClient"
	"gitee.com/i-Things/share/clients/wxClient"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type Clients struct {
	WxOfficial  *wxClient.WxOfficialAccount
	MiniProgram *wxClient.MiniProgram
	DingMini    *dingClient.DingTalk
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

func (c *ClientsManage) GetClients(ctx context.Context, appCode string) (Clients, error) {
	uc := ctxs.GetUserCtx(ctx)
	if appCode == "" {
		appCode = uc.AppCode
	}
	var tenantCode = uc.TenantCode
	logx.WithContext(ctx).Error(utils.Fmt(uc))
	val, ok := tc.Load(tenantCode + appCode)
	if ok {
		return val.(Clients), nil
	}
	//如果缓存里没有,需要查库
	cfg, err := relationDB.NewTenantAppRepo(ctx).FindOneByFilter(ctx, relationDB.TenantAppFilter{TenantCode: tenantCode, AppCodes: []string{appCode}})
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return Clients{}, err
		}
	}
	var cli Clients
	if cfg.DingMini != nil && cfg.DingMini.AppSecret != "" {
		cli.DingMini, err = dingClient.NewDingTalkClient(&conf.ThirdConf{
			AppID:     cfg.DingMini.AppID,
			AppKey:    cfg.DingMini.AppKey,
			AppSecret: cfg.DingMini.AppSecret,
		})
		if err != nil {
			return Clients{}, err
		}
	}
	if cfg.WxMini != nil && cfg.WxMini.AppSecret != "" {
		cli.MiniProgram, _ = wxClient.NewWxMiniProgram(ctx, &conf.ThirdConf{
			AppID:     cfg.WxMini.AppID,
			AppKey:    cfg.WxMini.AppKey,
			AppSecret: cfg.WxMini.AppSecret,
		}, c.Config.CacheRedis)
	}
	if cfg.WxOpen != nil && cfg.WxOpen.AppSecret != "" {
		cli.WxOfficial, _ = wxClient.NewWxOfficialAccount(ctx, &conf.ThirdConf{
			AppID:     cfg.WxOpen.AppID,
			AppKey:    cfg.WxOpen.AppKey,
			AppSecret: cfg.WxOpen.AppSecret,
		}, c.Config.CacheRedis)
	}
	tc.Store(tenantCode, cli)
	return cli, nil
}
