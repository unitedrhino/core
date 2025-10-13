package cache

import (
	"context"
	"sync"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/clients/dingClient"
	"gitee.com/unitedrhino/share/clients/wxClient"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
)

type Clients struct {
	WxOfficial  *wxClient.WxOfficialAccount
	MiniProgram *wxClient.MiniProgram
	DingMini    *dingClient.DingTalk
	Config      *relationDB.SysTenantApp
}
type ThirdClientsManage struct {
	Config cache.ClusterConf
	sf     syncx.SingleFlight
}

type appConfig struct {
	WxOpen    *wxClient.WxOfficialAccount
	WxMini    *wxClient.MiniProgram
	DingMini  *dingClient.DingTalk
	TenantApp map[string]map[string]*relationDB.SysTenantApp //第一层是appCode,第二层是appID
}

func (a *appConfig) Validate(tenantCode string, appCode string) error {
	l1, ok := a.TenantApp[appCode]
	if !ok {
		return errors.Parameter.AddMsg("未配置应用")
	}
	if l1[tenantCode] == nil || l1[def.TenantCodeCommon] == nil { //兼容通用应用
		return errors.Parameter.AddMsg("未配置应用")
	}
	return nil
}

var (
	tc       = sync.Map{}
	wxOpen   = sync.Map{}
	wxMini   = sync.Map{}
	dingMini = sync.Map{}

	wxOpenMutex   = sync.RWMutex{}
	wxMiniMutex   = sync.RWMutex{}
	dingMiniMutex = sync.RWMutex{}
)

func NewThirdClientsManage(c cache.ClusterConf) *ThirdClientsManage {
	return &ThirdClientsManage{Config: c, sf: syncx.NewSingleFlight()}
}

func (c *ThirdClientsManage) Init(ctx context.Context) error {
	tas, err := relationDB.NewTenantAppRepo(ctx).FindByFilter(ctx, relationDB.TenantAppFilter{WithApp: true}, nil)
	if err != nil {
		return err
	}
	for _, ta := range tas {
		err = c.SetOne(ctx, ta)
		if err != nil {
			logx.WithContext(ctx).Error(utils.Fmt(ta), err)
		}
	}
	return nil
}

func (c *ThirdClientsManage) DelOne(ctx context.Context, tenantCode string, appCode string) error {
	return c.SetOne(ctxs.BindTenantCode(ctx, tenantCode, 0), &relationDB.SysTenantApp{TenantCode: dataType.TenantCode(tenantCode), AppCode: appCode})
}
func (c *ThirdClientsManage) SetOne(ctx context.Context, ta *relationDB.SysTenantApp) error {
	setAppConfig := func(i *appConfig) {
		if i.TenantApp == nil {
			i.TenantApp = map[string]map[string]*relationDB.SysTenantApp{ta.AppCode: {string(ta.TenantCode): ta}}
		} else {
			if i.TenantApp[ta.AppCode] == nil {
				i.TenantApp[ta.AppCode] = map[string]*relationDB.SysTenantApp{string(ta.TenantCode): ta}
			} else {
				i.TenantApp[ta.AppCode][string(ta.TenantCode)] = ta
			}
		}
		if ta.App.IsCommon == def.True && ta.TenantCode == def.TenantCodeDefault { //如果是公共的,需要兼容
			i.TenantApp[ta.AppCode][def.TenantCodeCommon] = ta
		}
	}
	if ta.WxOpen != nil && ta.WxOpen.AppID != "" {
		cli, err := wxClient.NewWxOfficialAccount(ctx, &conf.ThirdConf{
			AppID:     ta.WxOpen.AppID,
			AppKey:    ta.WxOpen.AppKey,
			AppSecret: ta.WxOpen.AppSecret,
		}, c.Config)
		if err != nil {
			return errors.Parameter.AddMsgf("微信开放平台配置有误,错误:%v", err.Error())
		}
		info, ok := wxOpen.Load(ta.WxOpen.AppID)
		if !ok {
			wxOpen.Store(ta.WxOpen.AppID, &appConfig{
				WxOpen:    cli,
				TenantApp: map[string]map[string]*relationDB.SysTenantApp{ta.AppCode: {string(ta.TenantCode): ta}},
			})
		} else {
			i := info.(*appConfig)
			i.WxOpen = cli
			setAppConfig(i)
		}
	}

	if ta.WxMini != nil && ta.WxMini.AppID != "" {
		cli, err := wxClient.NewWxMiniProgram(ctx, &conf.ThirdConf{
			AppID:     ta.WxMini.AppID,
			AppKey:    ta.WxMini.AppKey,
			AppSecret: ta.WxMini.AppSecret,
		}, c.Config)
		if err != nil {
			return errors.Parameter.AddMsgf("微信小程序配置有误,错误:%v", err.Error())
		}
		info, ok := wxOpen.Load(ta.WxOpen.AppID)
		if !ok {
			wxMini.Store(ta.WxMini.AppID, &appConfig{
				WxMini:    cli,
				TenantApp: map[string]map[string]*relationDB.SysTenantApp{ta.AppCode: {string(ta.TenantCode): ta}},
			})
		} else {
			i := info.(*appConfig)
			i.WxMini = cli
			setAppConfig(i)
		}
	}
	if ta.DingMini != nil && ta.DingMini.AppID != "" {
		cli, err := dingClient.NewDingTalkClient(&conf.ThirdConf{
			AppID:     ta.DingMini.AppID,
			AppKey:    ta.DingMini.AppKey,
			AppSecret: ta.DingMini.AppSecret,
		})
		if err != nil {
			return errors.Parameter.AddMsgf("钉钉配置有误,错误:%v", err.Error())
		}
		info, ok := dingMini.Load(ta.DingMini.AppID)
		if !ok {
			dingMini.Store(ta.DingMini.AppID, &appConfig{
				DingMini:  cli,
				TenantApp: map[string]map[string]*relationDB.SysTenantApp{ta.AppCode: {string(ta.TenantCode): ta}},
			})
		} else {
			i := info.(*appConfig)
			i.DingMini = cli
			setAppConfig(i)
		}
	}
	return nil
}

func (c *ThirdClientsManage) GetWxMiniClient(ctx context.Context, appCode string, appID string) (*wxClient.MiniProgram, error) {
	ac, ok := wxMini.Load(appID)
	if !ok {
		return nil, errors.Parameter.AddMsg("未配置应用")
	}
	cfg := ac.(*appConfig)
	err := cfg.Validate(ctxs.GetUserCtxNoNil(ctx).TenantCode, appCode)
	if err != nil {
		return nil, err
	}
	return cfg.WxMini, nil
}

func (c *ThirdClientsManage) GetWxOpenClient(ctx context.Context, appCode string, appID string) (*wxClient.WxOfficialAccount, error) {
	ac, ok := wxOpen.Load(appID)
	if !ok {
		return nil, errors.Parameter.AddMsg("未配置应用")
	}
	cfg := ac.(*appConfig)
	err := cfg.Validate(ctxs.GetUserCtxNoNil(ctx).TenantCode, appCode)
	if err != nil {
		return nil, err
	}
	return cfg.WxOpen, nil
}

func (c *ThirdClientsManage) GetDingAppClient(ctx context.Context, appCode string, appID string) (*dingClient.DingTalk, error) {
	ac, ok := dingMini.Load(appID)
	if !ok {
		return nil, errors.Parameter.AddMsg("未配置应用")
	}
	cfg := ac.(*appConfig)
	if appCode == "" {
		appCode = ctxs.GetUserCtxNoNil(ctx).AppCode
	}
	err := cfg.Validate(ctxs.GetUserCtxNoNil(ctx).TenantCode, appCode)
	if err != nil {
		return nil, err
	}
	return cfg.DingMini, nil
}
func (c *ThirdClientsManage) GetClients(ctx context.Context, appCode string) (Clients, error) {
	uc := ctxs.GetUserCtx(ctx)
	if appCode == "" {
		appCode = uc.AppCode
	}
	var tenantCode = uc.TenantCode
	var key = tenantCode + ":" + appCode
	val, ok := tc.Load(tenantCode + appCode)
	if ok {
		return val.(Clients), nil
	}
	cli, err := c.sf.Do(key, func() (any, error) {
		//如果缓存里没有,需要查库
		cfg, err := relationDB.NewTenantAppRepo(ctx).FindOneByFilter(ctx, relationDB.TenantAppFilter{TenantCode: tenantCode, AppCodes: []string{appCode}})
		if err != nil {
			if !errors.Cmp(err, errors.NotFind) {
				return Clients{}, err
			}
			return Clients{}, errors.Parameter.AddMsg("未配置应用")
		}
		var cli Clients
		cli.Config = cfg
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
			cli.MiniProgram, err = wxClient.NewWxMiniProgram(ctx, &conf.ThirdConf{
				AppID:     cfg.WxMini.AppID,
				AppKey:    cfg.WxMini.AppKey,
				AppSecret: cfg.WxMini.AppSecret,
			}, c.Config)
			if err != nil {
				return Clients{}, err
			}
		}
		if cfg.WxOpen != nil && cfg.WxOpen.AppSecret != "" {
			cli.WxOfficial, err = wxClient.NewWxOfficialAccount(ctx, &conf.ThirdConf{
				AppID:     cfg.WxOpen.AppID,
				AppKey:    cfg.WxOpen.AppKey,
				AppSecret: cfg.WxOpen.AppSecret,
			}, c.Config)
			if err != nil {
				return Clients{}, err
			}
		}
		tc.Store(tenantCode, cli)
		return cli, nil
	})

	return cli.(Clients), err
}
