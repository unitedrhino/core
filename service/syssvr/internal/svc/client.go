package svc

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/config"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type Clients struct {
	MiniProgram *clients.MiniProgram
	MiniDing    *clients.DingTalk
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
	uc := ctxs.GetUserCtx(ctx)
	if tenantCode == "" {
		tenantCode = uc.TenantCode
	}
	logx.WithContext(ctx).Error(utils.Fmt(uc))
	val, ok := tc.Load(tenantCode + uc.AppCode)
	if ok {
		return val.(Clients), nil
	}
	//如果缓存里没有,需要查库
	cfg, err := relationDB.NewTenantAppRepo(ctx).FindOneByFilter(ctx, relationDB.TenantAppFilter{TenantCode: tenantCode, AppCodes: []string{uc.AppCode}})
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return Clients{}, err
		}
		cfg2, err := relationDB.NewAppInfoRepo(ctx).FindOneByFilter(ctx, relationDB.AppInfoFilter{Code: uc.AppCode})
		if err != nil {
			return Clients{}, err
		}
		cfg = &relationDB.SysTenantApp{
			MiniWx: cfg2.MiniWx,
		}
	}
	var cli Clients
	if cfg.MiniDing != nil && cfg.MiniDing.AppSecret != "" {
		cli.MiniDing, err = clients.NewDingTalkClient(&conf.ThirdConf{
			AppID:     cfg.MiniDing.AppID,
			AppKey:    cfg.MiniDing.AppKey,
			AppSecret: cfg.MiniDing.AppSecret,
		})
		if err != nil {
			return Clients{}, err
		}
	}
	if cfg.MiniWx != nil && cfg.MiniWx.AppSecret != "" {
		cli.MiniProgram, _ = clients.NewWxMiniProgram(ctx, &conf.ThirdConf{
			AppID:     cfg.MiniWx.AppID,
			AppKey:    cfg.MiniWx.AppKey,
			AppSecret: cfg.MiniWx.AppSecret,
		}, c.Config.CacheRedis)
	}
	tc.Store(tenantCode, cli)
	return cli, nil
}
