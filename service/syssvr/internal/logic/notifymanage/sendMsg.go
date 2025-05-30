package notifymanagelogic

import (
	"bytes"
	"context"
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
	"github.com/silenceper/wechat/v2/miniprogram/subscribe"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zhaoyunxing92/dingtalk/v2/request"
	"gorm.io/gorm"
	"strings"
	"text/template"
)

type SendMsgConfig struct {
	UserIDs       []int64 //只有填写了这项才会记录
	Accounts      []string
	AccountType   string
	NotifyCode    string         //通知的code
	BakNotifyCode string         //备用通知code
	TemplateID    int64          //指定模板
	Type          def.NotifyType //通知类型
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
		c, err := relationDB.NewNotifyConfigTemplateRepo(ctx).FindOneByFilter(ctx, relationDB.NotifyConfigTemplateFilter{
			NotifyCode: cfg.NotifyCode,
			Type:       cfg.Type,
		})
		if err != nil {
			if errors.Cmp(err, errors.NotFind) {
				if len(cfg.BakNotifyCode) > 0 {
					c, err = relationDB.NewNotifyConfigTemplateRepo(ctx).FindOneByFilter(ctx, relationDB.NotifyConfigTemplateFilter{
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
		temp = c.Template
		channel = c.Template.Channel
		config = c.Config
	}

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
	var users []*relationDB.SysUserInfo
	if config.IsRecord == def.True && !cfg.NoRecord { //需要记录到消息中心中
		users, err = relationDB.NewUserInfoRepo(ctx).FindUserCore(ctx, relationDB.UserInfoFilter{UserIDs: cfg.UserIDs, Accounts: cfg.Accounts})
		if err != nil {
			return err
		}
		if len(users) != 0 {
			err = stores.GetTenantConn(ctx).Transaction(func(tx *gorm.DB) error {
				mi := relationDB.NewMessageInfoRepo(tx)
				miPo := relationDB.SysMessageInfo{
					Group:          config.Group,
					NotifyCode:     cfg.NotifyCode,
					Subject:        subject,
					Body:           body,
					Str1:           cfg.Str1,
					Str2:           cfg.Str2,
					Str3:           cfg.Str3,
					IsDirectNotify: def.True,
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
				return err
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
	case def.NotifyTypeWxMini:
		if channel == nil {
			return errors.Parameter.AddMsg("没有配置微信小程序")
		}
		cli, err := svcCtx.Cm.GetClients(ctx, channel.AppCode)
		if err != nil {
			return err
		}
		if cli.MiniProgram == nil {
			return errors.Parameter.AddMsg("没有配置微信小程序")
		}
		var userIDs []string
		for _, v := range users {
			if v.WechatOpenID.Valid {
				userIDs = append(userIDs, cast.ToString(v.WechatOpenID.String))
			}
		}
		if len(userIDs) > 0 {
			for _, user := range userIDs {
				err = cli.MiniProgram.GetSubscribe().Send(&subscribe.Message{
					ToUser:     user,
					TemplateID: templateCode,
					Data:       map[string]*subscribe.DataItem{},
				})
				if err != nil {
					return err
				}
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
		err = utils.SendEmail(conf.Email{
			From:     temp.Channel.Email.From,
			Host:     temp.Channel.Email.Host,
			Secret:   temp.Channel.Email.Secret,
			Nickname: temp.Channel.Email.Nickname,
			Port:     temp.Channel.Email.Port,
			IsSSL:    temp.Channel.Email.IsSSL == def.True,
		}, accounts, subject,
			body)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
		return err
	}
	return nil
}
