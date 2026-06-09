package notifymanagelogic

import (
	"context"
	"strings"

	"gitee.com/unitedrhino/core/service/syssvr/internal/defext"
	"gitee.com/unitedrhino/core/service/syssvr/internal/pkg/unipush"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

var defaultUniPushNotifyCodes = []string{
	"ruleDeviceAlarm",
	"ruleScene",
}

func shouldSendAppPush(cfg SendMsgConfig, notifyCode string, svcCtx *svc.ServiceContext) bool {
	if svcCtx == nil || svcCtx.UniPush == nil || !svcCtx.UniPush.Enabled() {
		return false
	}
	if notifyCode == "" {
		return false
	}
	switch cfg.Type {
	case defext.NotifyTypeSystemNotice:
		// systemNotice 专走 uni-push
	case def.NotifyTypeMessage:
		// 场景联动：站内信只进消息中心，系统栏推送由 systemNotice 承担，避免同场景配两种渠道时推两条
		if notifyCode == "ruleScene" {
			return false
		}
		if !svcCtx.Config.UniPush.PushOnMessage {
			return false
		}
	default:
		return false
	}
	whitelist := svcCtx.Config.UniPush.NotifyCodes
	if len(whitelist) == 0 {
		whitelist = defaultUniPushNotifyCodes
	}
	for _, code := range whitelist {
		if code == notifyCode {
			return true
		}
	}
	return false
}

func pushAppNotifyAsync(parent context.Context, svcCtx *svc.ServiceContext, userIDs []int64, subject, body, notifyCode, notifyType string, messageID int64, params map[string]any) {
	if len(userIDs) == 0 || svcCtx == nil || svcCtx.UniPush == nil {
		return
	}
	// gRPC/HTTP 返回后 parent 会被取消，异步推送到 uniCloud 须脱离取消链
	parent = context.WithoutCancel(parent)
	ctxs.GoNewCtx(parent, func(ctx context.Context) {
		if err := pushAppNotify(ctx, svcCtx, userIDs, subject, body, notifyCode, notifyType, messageID, params); err != nil {
			logx.WithContext(ctx).Errorf("unipush pushAppNotify err:%v notifyCode=%s", err, notifyCode)
		}
	})
}

func pushAppNotify(ctx context.Context, svcCtx *svc.ServiceContext, userIDs []int64, subject, body, notifyCode, notifyType string, messageID int64, params map[string]any) error {
	cids, err := relationDB.NewUserPushClientRepo(ctx).ListActiveClientIDs(ctx, userIDs)
	if err != nil {
		return err
	}
	if len(cids) == 0 {
		logx.WithContext(ctx).Slowf("unipush skip: no active push client ids userIDs=%v notifyCode=%s", userIDs, notifyCode)
		return nil
	}
	title := strings.TrimSpace(subject)
	content := strings.TrimSpace(body)
	tu := triggerUserFieldsFromParams(params)
	// ruleScene 的正文已在 SendNotifyMsg 中拼接触发人，避免重复
	if notifyCode != "ruleScene" {
		if line := FormatTriggerUserLine(tu); line != "" {
			content = PrependTriggerUserLine(content, tu)
		}
	}
	if title == "" {
		title = "YK AIoT"
	}
	if content == "" {
		content = title
	}
	route := "warn"
	warnTab := 0
	switch notifyCode {
	case "ruleScene":
		if notifyType == string(defext.NotifyTypeSystemNotice) {
			warnTab = 1
		}
	case "ruleDeviceAlarm":
		warnTab = 0
	default:
		route = ""
	}
	payload := map[string]any{
		"notifyCode": notifyCode,
		"route":      route,
		"notifyType": notifyType,
		"warnTab":    warnTab,
	}
	if notifyType == string(defext.NotifyTypeSystemNotice) {
		// Android 在线/后台须 force_notification 才弹系统通知栏；鸿蒙不识别该字段，靠 payload.showLocalNotice + 客户端 createPushMessage。
		payload["showLocalNotice"] = true
		payload["title"] = title
		payload["content"] = content
	}
	if messageID > 0 {
		payload["messageId"] = messageID
	}
	for k, v := range params {
		payload[k] = v
	}
	return svcCtx.UniPush.Send(ctx, unipush.SendReq{
		PushClientIDs: cids,
		Title:         title,
		Content:       content,
		Payload:       payload,
		RequestID:     uuid.NewString(),
	})
}
