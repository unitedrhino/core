package notifymanagelogic

import (
	"context"
	"strings"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"github.com/spf13/cast"
)

const (
	triggerTypeManual       = "manual"
	triggerTypeAuto         = "auto"
	triggerTypeDeviceButton = "deviceButton" // 物理按键/属性上报，推断的分享用户触发
	triggerUserNickAuto     = "系统自动"
	paramTriggerUserID   = "triggerUserId"
	paramTriggerUserNick = "triggerUserNick"
	paramTriggerUserAcct = "triggerUserAccount"
	paramTriggerType     = "triggerType"
)

type triggerUserFields struct {
	UserID  int64
	Nick    string
	Account string
	Type    string
}

func enrichTriggerUserParams(ctx context.Context, params map[string]any) triggerUserFields {
	out := triggerUserFields{}
	if params == nil {
		return out
	}
	out.Type = cast.ToString(params[paramTriggerType])
	out.UserID = cast.ToInt64(params[paramTriggerUserID])
	out.Nick = cast.ToString(params[paramTriggerUserNick])
	out.Account = cast.ToString(params[paramTriggerUserAcct])

	if out.Type == "" && out.UserID <= 0 {
		return out
	}
	if out.Type == "" {
		out.Type = triggerTypeManual
		params[paramTriggerType] = out.Type
	}
	if out.UserID > 0 {
		fillTriggerUserFromDB(ctx, &out, params)
	} else if out.Type == triggerTypeAuto && strings.TrimSpace(out.Nick) == "" {
		out.Nick = triggerUserNickAuto
		params[paramTriggerUserNick] = out.Nick
	}
	return out
}

func fillTriggerUserFromDB(ctx context.Context, out *triggerUserFields, params map[string]any) {
	if out == nil || out.UserID <= 0 {
		return
	}
	u, err := relationDB.NewUserInfoRepo(ctx).FindOne(ctx, out.UserID)
	if err != nil || u == nil {
		return
	}
	if strings.TrimSpace(out.Nick) == "" {
		out.Nick = triggerUserNick(u)
	}
	if strings.TrimSpace(out.Account) == "" {
		out.Account = triggerUserAccount(u)
	}
	params[paramTriggerUserNick] = out.Nick
	params[paramTriggerUserAcct] = out.Account
}

func triggerUserNick(u *relationDB.SysUserInfo) string {
	if u == nil {
		return ""
	}
	if n := strings.TrimSpace(u.NickName); n != "" {
		return n
	}
	if u.UserName.Valid {
		return strings.TrimSpace(u.UserName.String)
	}
	return triggerUserAccount(u)
}

func triggerUserAccount(u *relationDB.SysUserInfo) string {
	if u == nil {
		return ""
	}
	if u.Phone.Valid && u.Phone.String != "" {
		return u.Phone.String
	}
	if u.Email.Valid && u.Email.String != "" {
		return u.Email.String
	}
	if u.UserName.Valid {
		return u.UserName.String
	}
	return ""
}

func triggerUserDisplayName(t triggerUserFields) string {
	name := strings.TrimSpace(t.Account)
	if name == "" {
		name = strings.TrimSpace(t.Nick)
	}
	if name == "" && t.UserID > 0 {
		name = cast.ToString(t.UserID)
	}
	return name
}

// compactRuleSceneActionBody 多行模板正文合并为一条动作描述（取信息最完整的一行）。
func compactRuleSceneActionBody(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return ""
	}
	lines := strings.Split(body, "\n")
	var parts []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			parts = append(parts, l)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	best := parts[0]
	for _, l := range parts[1:] {
		if len([]rune(l)) > len([]rune(best)) {
			best = l
		}
	}
	return best
}

// FormatSceneButtonTriggerLine 物理按键场景：{账号}触发了{动作}（不含「的账号」字样）。
func FormatSceneButtonTriggerLine(t triggerUserFields, actionName string) string {
	name := triggerUserDisplayName(t)
	action := strings.TrimSpace(actionName)
	if action == "" {
		return FormatTriggerUserLine(t)
	}
	if name == "" {
		return "触发了" + action
	}
	return name + "触发了" + action
}

// FormatRuleSceneNotifyLine ruleScene 通知正文前的触发人一行。
func FormatRuleSceneNotifyLine(t triggerUserFields, renderedBody string) string {
	if t.Type == triggerTypeDeviceButton && t.UserID > 0 {
		action := compactRuleSceneActionBody(renderedBody)
		if action == "" {
			action = strings.TrimSpace(renderedBody)
		}
		return FormatSceneButtonTriggerLine(t, action)
	}
	return FormatTriggerUserLine(t)
}

// FormatTriggerUserLine 触发人展示文案（与 App msg_trigger_* 文案一致，用于推送正文等）。
func FormatTriggerUserLine(t triggerUserFields) string {
	name := triggerUserDisplayName(t)
	switch t.Type {
	case triggerTypeManual:
		if name != "" {
			return "由 " + name + " 手动触发"
		}
		return "由用户手动触发"
	case triggerTypeAuto:
		if name != "" && name != triggerUserNickAuto {
			return "由 " + name + " 触发"
		}
		return triggerUserNickAuto + "触发"
	default:
		if name != "" {
			return "由 " + name + " 触发"
		}
		return ""
	}
}

func triggerUserFieldsFromParams(params map[string]any) triggerUserFields {
	if params == nil {
		return triggerUserFields{}
	}
	return triggerUserFields{
		Type:    cast.ToString(params[paramTriggerType]),
		UserID:  cast.ToInt64(params[paramTriggerUserID]),
		Nick:    cast.ToString(params[paramTriggerUserNick]),
		Account: cast.ToString(params[paramTriggerUserAcct]),
	}
}

// PrependRuleSceneTriggerLine 场景联动通知：在正文前展示触发人（含 deviceButton 专用格式）。
func PrependRuleSceneTriggerLine(body string, t triggerUserFields) string {
	return prependNotifyLine(body, FormatRuleSceneNotifyLine(t, body))
}

// PrependTriggerUserLine 在正文开头展示触发人（已包含则不再重复）。
func PrependTriggerUserLine(body string, t triggerUserFields) string {
	line := FormatTriggerUserLine(t)
	return prependNotifyLine(body, line)
}

func prependNotifyLine(body, line string) string {
	if line == "" {
		return body
	}
	body = strings.TrimSpace(body)
	if body == line {
		return line
	}
	if strings.HasPrefix(body, line+"\n") {
		return body
	}
	if strings.HasPrefix(body, line) && (len(body) == len(line) || body[len(line)] == '\n') {
		return body
	}
	// 兼容历史数据：触发人在末尾时移到逻辑上的「前面」展示
	if strings.HasSuffix(body, "\n"+line) {
		body = strings.TrimSpace(body[:len(body)-len(line)-1])
	} else if strings.HasSuffix(body, line) && (len(body) == len(line) || body[len(body)-len(line)-1] == '\n') {
		body = strings.TrimSpace(strings.TrimSuffix(body, line))
	}
	if body == "" {
		return line
	}
	return line + "\n" + body
}

// AppendTriggerUserLine 兼容旧名，实际在正文前展示触发人。
func AppendTriggerUserLine(body string, t triggerUserFields) string {
	return PrependTriggerUserLine(body, t)
}
