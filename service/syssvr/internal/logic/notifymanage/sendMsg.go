package notifymanagelogic

import (
	"bytes"
	"context"
	"strings"
	"text/template"

	"gitee.com/unitedrhino/core/service/syssvr/internal/defext"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/share/clients/dingClient"
	"gitee.com/unitedrhino/share/clients/smsClient"
	"gitee.com/unitedrhino/share/clients/wxClient"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zhaoyunxing92/dingtalk/v2/request"
	"gorm.io/gorm"
)

type SendMsgConfig struct {
	UserIDs       []int64 // 只有填写了这项才会记录
	Accounts      []string
	AccountType   string
	NotifyCode    string         // 通知的code
	BakNotifyCode string         // 备用通知code
	TemplateID    int64          // 指定模板
	Type          def.NotifyType // 通知类型
	Params        map[string]any
	Str1          string
	Str2          string
	Str3          string
	NoRecord      bool
}

func SendNotifyMsg(ctx context.Context, svcCtx *svc.ServiceContext, cfg SendMsgConfig) error {
	var (
		subject      string
		body         string
		signName     string
		templateCode string
		err          error
		temp         *relationDB.SysNotifyTemplate
		channel      *relationDB.SysNotifyChannel
		config       *relationDB.SysNotifyConfig
	)
	if cfg.TemplateID != 0 {
		t, err := relationDB.NewNotifyTemplateRepo(ctx).FindOne(ctx, cfg.TemplateID)
		if err != nil {
			return err
		}
		temp = t
		channel = t.Channel
		config = t.Config
	} else {
		// c, err := relationDB.NewNotifyConfigTemplateRepo(ctx).FindOneByFilter(ctx, relationDB.NotifyConfigTemplateFilter{
		// 	NotifyCode: cfg.NotifyCode,
		// 	Type:       cfg.Type,
		// })
		c, err := relationDB.NewNotifyTemplateRepo(ctx).FindOneByFilter(ctx, relationDB.NotifyTemplateFilter{
			NotifyCode: cfg.NotifyCode,
			Type:       cfg.Type,
		})
		if err != nil {
			if errors.Cmp(err, errors.NotFind) {
				if len(cfg.BakNotifyCode) > 0 {
					// c, err = relationDB.NewNotifyConfigTemplateRepo(ctx).FindOneByFilter(ctx, relationDB.NotifyConfigTemplateFilter{
					// 	NotifyCode: cfg.NotifyCode,
					// 	Type:       cfg.Type,
					// })
					c, err = relationDB.NewNotifyTemplateRepo(ctx).FindOneByFilter(ctx, relationDB.NotifyTemplateFilter{
						NotifyCode: cfg.NotifyCode,
						Type:       cfg.Type,
					})
					if errors.Cmp(err, errors.NotFind) {
						return errors.NotEnable
					}
					return err
				}
				return errors.NotEnable
			}
			return err
		}
		temp = c
		// temp = c.Template
		channel = c.Channel
		config = c.Config
	}
	if channel == nil && temp != nil && temp.ChannelID != 0 {
		channel, err = relationDB.NewNotifyChannelRepo(ctx).FindOne(ctx, temp.ChannelID)
		if err != nil {
			return err
		}
	}

	// 模板 Preload Config 关联偶发为空，按 notifyCode 兜底加载，避免 config 空指针。
	if config == nil && temp != nil && temp.NotifyCode != "" {
		config, err = relationDB.NewNotifyConfigRepo(ctx).FindOneByFilter(ctx, relationDB.NotifyConfigFilter{
			Code: temp.NotifyCode,
		})
		if err != nil {
			return err
		}
	}
	if config == nil {
		return errors.NotEnable.AddMsg("通知配置不存在")
	}
	if len(config.EnableTypes) > 0 && !utils.SliceIn(cfg.Type, config.EnableTypes...) {
		return errors.NotEnable.AddMsg("通知类型未启用:" + string(cfg.Type))
	}

	triggerUser := enrichTriggerUserParams(ctx, cfg.Params)

	if temp != nil {
		subject = temp.Subject
		body = temp.Body
		signName = temp.SignName
		templateCode = temp.TemplateCode
	}
	{
		tmpl, err := template.New(config.Code).Parse(body)
		if err != nil {
			return errors.System.AddMsg("模版解析失败").AddDetail(err)
		}
		buffer := &bytes.Buffer{}
		err = tmpl.Execute(buffer, cfg.Params)
		if err != nil {
			return errors.System.AddMsg("模版匹配失败").AddDetail(err)
		}
		body = buffer.String()
	}
	{
		tmpl, err := template.New(config.Code).Parse(subject)
		if err != nil {
			return errors.System.AddMsg("模版解析失败").AddDetail(err)
		}
		buffer := &bytes.Buffer{}
		err = tmpl.Execute(buffer, cfg.Params)
		if err != nil {
			return errors.System.AddMsg("模版匹配失败").AddDetail(err)
		}
		subject = buffer.String()
		if subject == "" {
			subject = config.Name
		}
	}
	if triggerUser.Type != "" && cfg.NotifyCode == "ruleScene" {
		if triggerUser.Type == triggerTypeDeviceButton && cfg.Type == defext.NotifyTypeSystemNotice {
			// 系统推送：仅一行「触发了{动作}」，避免「手机号+的账号」与正文重复多行
			body = FormatRuleSceneNotifyLine(triggerUser, body)
		} else {
			body = PrependRuleSceneTriggerLine(body, triggerUser)
		}
	}
	var users []*relationDB.SysUserInfo
	var messageID int64
	needRecord := config.IsRecord == def.True && !cfg.NoRecord
	needUsers := needRecord || shouldSendAppPush(cfg, string(cfg.NotifyCode), svcCtx)
	if needUsers && (len(cfg.UserIDs) != 0 || len(cfg.Accounts) != 0) {
		users, err = relationDB.NewUserInfoRepo(ctx).FindUserCore(ctx, relationDB.UserInfoFilter{UserIDs: cfg.UserIDs, Accounts: cfg.Accounts})
		if err != nil {
			return err
		}
	}
	if needRecord { // 需要记录到消息中心中
		if len(users) != 0 {
			miPo := relationDB.SysMessageInfo{
				Group:              config.Group,
				NotifyCode:         cfg.NotifyCode,
				NotifyType:         string(cfg.Type),
				Subject:            subject,
				Body:               body,
				Str1:               cfg.Str1,
				Str2:               cfg.Str2,
				Str3:               cfg.Str3,
				TriggerUserID:      triggerUser.UserID,
				TriggerUserNick:    triggerUser.Nick,
				TriggerUserAccount: triggerUser.Account,
				TriggerType:        triggerUser.Type,
				IsDirectNotify:     def.True,
			}
			err = stores.GetTenantConn(ctx).Transaction(func(tx *gorm.DB) error {
				mi := relationDB.NewMessageInfoRepo(tx)
				err := mi.Insert(ctx, &miPo)
				if err != nil {
					return err
				}
				var umPos []*relationDB.SysUserMessage
				for _, v := range users {
					umPos = append(umPos, &relationDB.SysUserMessage{
						UserID:    v.UserID,
						MessageID: miPo.ID,
						Group:     miPo.Group,
						IsRead:    def.False,
					})
				}
				return relationDB.NewUserMessageRepo(tx).MultiInsert(ctx, umPos)
			})
			if err != nil {
				return err
			}
			messageID = miPo.ID
		}
	}
	if shouldSendAppPush(cfg, string(cfg.NotifyCode), svcCtx) {
		// 合并 UserIDs 与 Accounts 解析出的用户（场景自动触发时常见：UserIDs=设备主人 + Accounts=创建者手机号）
		pushUserIDs := make([]int64, 0, len(cfg.UserIDs)+len(users))
		seenPushUID := map[int64]struct{}{}
		addPushUID := func(uid int64) {
			if uid <= 0 {
				return
			}
			if _, ok := seenPushUID[uid]; ok {
				return
			}
			seenPushUID[uid] = struct{}{}
			pushUserIDs = append(pushUserIDs, uid)
		}
		for _, uid := range cfg.UserIDs {
			addPushUID(uid)
		}
		for _, u := range users {
			addPushUID(u.UserID)
		}
		if len(pushUserIDs) > 0 {
			pushAppNotifyAsync(ctx, svcCtx, pushUserIDs, subject, body, string(cfg.NotifyCode), string(cfg.Type), messageID, cfg.Params)
		} else if len(cfg.Accounts) > 0 || len(cfg.UserIDs) > 0 {
			logx.WithContext(ctx).Slowf("unipush skip: no push target user resolved notifyCode=%s accounts=%v userIDs=%v",
				cfg.NotifyCode, cfg.Accounts, cfg.UserIDs)
		}
	}
	switch cfg.Type {
	case def.NotifyTypeSms:
		var accounts = cfg.Accounts
		if len(users) != 0 {
			for _, v := range users {
				if v.Phone.Valid {
					accounts = append(accounts, v.Phone.String)
				}
			}
		}
		if len(accounts) == 0 {
			logx.WithContext(ctx).Slowf("sms skip: no phone accounts notifyCode=%s userIDs=%v accounts=%v",
				cfg.NotifyCode, cfg.UserIDs, cfg.Accounts)
			return nil
		}
		err = svcCtx.Sms.SendSms(ctx, smsClient.SendSmsParam{
			PhoneNumbers:  accounts,
			SignName:      signName,
			TemplateCode:  templateCode,
			TemplateParam: cfg.Params,
		})
		if err != nil {
			return err
		}
	case def.NotifyTypeDingWebhook:
		if channel == nil || channel.WebHook == "" {
			return errors.NotEnable.AddMsg("通道没有配置")
		}
		cli := dingClient.NewDingRobotClient(channel.WebHook)
		_, err := cli.SendRobotMsg(dingClient.NewTextMessage(body))
		return err
	case def.NotifyTypeDingTalk:
		if channel == nil || channel.App == nil {
			return errors.NotEnable.AddMsg("通道没有配置")
		}
		cli, err := dingClient.NewDingTalkClient(&conf.ThirdConf{
			AppKey:    channel.App.AppKey,
			AppSecret: channel.App.AppSecret,
		})
		if err != nil {
			return err
		}
		var userIDs []string
		for _, v := range users {
			if v.DingTalkUserID.Valid {
				userIDs = append(userIDs, cast.ToString(v.DingTalkUserID.String))
			}
		}
		if len(userIDs) > 0 {
			_, err = cli.SendCorpConvMessage(&request.CorpConvMessage{
				AgentId:    cast.ToInt(channel.App.AppID),
				UserIdList: strings.Join(userIDs, ","),
				Msg:        dingClient.NewTextMessage(body),
			})
			if err != nil {
				return err
			}
		}
	case def.NotifyTypeWxEWebhook:
		if channel == nil || channel.WebHook == "" {
			return errors.NotEnable.AddMsg("通道没有配置")
		}
		err := wxClient.SendRobotMsg(ctx, channel.WebHook, body)
		return err
	case def.NotifyTypeEmail:
		var accounts = cfg.Accounts
		if len(users) != 0 {
			for _, v := range users {
				if v.Email.Valid {
					accounts = append(accounts, v.Email.String)
				}
			}
		}
		if len(accounts) == 0 {
			return nil
		}
		emailConf, err := emailConfigFromChannel(channel)
		if err != nil {
			return err
		}
		err = utils.SendEmail(emailConf, accounts, subject,
			body)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
		return err
	case defext.NotifyTypeSystemNotice:
		// 系统推送：模板渲染完成后由 uni-push 旁路发送（可选 IsRecord 写消息中心）
		return nil
	}
	return nil
}

func emailConfigFromChannel(channel *relationDB.SysNotifyChannel) (conf.Email, error) {
	if channel == nil || channel.Email == nil {
		return conf.Email{}, errors.NotEnable.AddMsg("通道没有配置")
	}
	return conf.Email{
		From:     channel.Email.From,
		Host:     channel.Email.Host,
		Secret:   channel.Email.Secret,
		Nickname: channel.Email.Nickname,
		Port:     channel.Email.Port,
		IsSSL:    channel.Email.IsSSL == def.True,
	}, nil
}
