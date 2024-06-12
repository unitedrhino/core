package notifymanagelogic

import (
	"bytes"
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"gorm.io/gorm"
	"text/template"
)

type SendMsgConfig struct {
	UserIDs     []int64 //只有填写了这项才会记录
	Accounts    []string
	AccountType string
	NotifyCode  string         //通知的code
	Type        def.NotifyType //通知类型
	Params      map[string]any
	Str1        string
	Str2        string
	Str3        string
}

func SendNotifyMsg(ctx context.Context, svcCtx *svc.ServiceContext, cfg SendMsgConfig) error {
	c, err := relationDB.NewTenantNotifyTemplateRepo(ctx).FindOneByFilter(ctxs.WithCommonTenant(ctx), relationDB.TenantNotifyTemplateFilter{
		NotifyCode: cfg.NotifyCode,
		Type:       cfg.Type,
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return errors.NotEnable
		}
		return err
	}
	var (
		subject      = c.Config.DefaultSubject
		body         = c.Config.DefaultBody
		signName     = c.Config.DefaultSignName
		templateCode = c.Config.DefaultTemplateCode
	)
	if c.Template != nil {
		subject = c.Template.Subject
		body = c.Template.Body
		signName = c.Template.SignName
		templateCode = c.Template.TemplateCode
	}
	{
		tmpl, err := template.New(c.Config.Code).Parse(body)
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
		tmpl, err := template.New(c.Config.Code).Parse(subject)
		if err != nil {
			return errors.System.AddMsg("模版解析失败").AddDetail(err)
		}
		buffer := &bytes.Buffer{}
		err = tmpl.Execute(buffer, cfg.Params)
		if err != nil {
			return errors.System.AddMsg("模版匹配失败").AddDetail(err)
		}
		subject = buffer.String()
	}
	var users []*relationDB.SysUserInfo
	if c.Config.IsRecord == def.True { //需要记录到消息中心中
		users, err = relationDB.NewUserInfoRepo(ctx).FindUserCore(ctx, relationDB.UserInfoFilter{UserIDs: cfg.UserIDs})
		if err != nil {
			return err
		}
		if len(users) != 0 {
			err = stores.GetTenantConn(ctx).Transaction(func(tx *gorm.DB) error {
				mi := relationDB.NewMessageInfoRepo(tx)
				miPo := relationDB.SysMessageInfo{
					Group:      c.Config.Group,
					NotifyCode: cfg.NotifyCode,
					Subject:    subject,
					Body:       body,
					Str1:       cfg.Str1,
					Str2:       cfg.Str2,
					Str3:       cfg.Str3,
				}
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
				return nil
			}
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
			return nil
		}
		err = svcCtx.Sms.SendSms(ctx, clients.SendSmsParam{
			PhoneNumbers:  accounts,
			SignName:      signName,
			TemplateCode:  templateCode,
			TemplateParam: cfg.Params,
		})
		if err != nil {
			return err
		}
	case def.NotifyTypeDingWebhook:
		channel := c.Template.Channel
		if channel == nil || channel.WebHook == "" {
			return errors.NotEnable.AddMsg("通道没有配置")
		}
		cli := clients.NewDingRobotClient(channel.WebHook)
		_, err := cli.SendRobotMsg(clients.NewTextMessage(body))
		return err
	case def.NotifyTypeDingTalk:
		channel := c.Template.Channel
		if channel == nil || channel.DingTalk == nil {
			return errors.NotEnable.AddMsg("通道没有配置")
		}
		//cli, err := clients.NewDingTalkClient(&conf.ThirdConf{
		//	AppKey:    channel.DingTalk.AppKey,
		//	AppSecret: channel.DingTalk.AppSecret,
		//})
		//if err != nil {
		//	return err
		//}
		//cli.SendCorpConvMessage()
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
		err = utils.SenEmail(conf.Email{
			From:     c.Template.Channel.Email.From,
			Host:     c.Template.Channel.Email.Host,
			Secret:   c.Template.Channel.Email.Secret,
			Nickname: c.Template.Channel.Email.Nickname,
			Port:     c.Template.Channel.Email.Port,
			IsSSL:    c.Template.Channel.Email.IsSSL == def.True,
		}, accounts, subject,
			body)
	}
	return nil
}
