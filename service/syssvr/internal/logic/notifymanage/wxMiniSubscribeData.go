package notifymanagelogic

import (
	"strings"
	"time"

	"github.com/silenceper/wechat/v2/miniprogram/subscribe"
	"github.com/spf13/cast"
)

// wxSubscribeValueMaxLen 微信订阅消息 thing 类字段常见长度上限
const wxSubscribeValueMaxLen = 20

func truncateWxSubscribeValue(s string) string {
	runes := []rune(s)
	if len(runes) <= wxSubscribeValueMaxLen {
		return s
	}
	return string(runes[:wxSubscribeValueMaxLen])
}

func firstNonEmpty(params map[string]any, keys ...string) string {
	for _, k := range keys {
		v, ok := params[k]
		if !ok {
			continue
		}
		s := strings.TrimSpace(cast.ToString(v))
		if s != "" {
			return s
		}
	}
	return ""
}

// formatWxSubscribeTime 格式化为微信 time 类字段常用展示
func formatWxSubscribeTime(t time.Time) string {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("CST", 8*3600)
	}
	return t.In(loc).Format("2006年01月02日 15:04")
}

// mapWxMiniTemplateFields 映射到公众平台模板字段（示例：thing1 设备名称、time2 操作时间、thing3 操作动作）
func mapWxMiniTemplateFields(params map[string]any) map[string]string {
	out := map[string]string{}
	if params == nil {
		params = map[string]any{}
	}
	if v := firstNonEmpty(params, "thing1", "deviceAlias", "deviceName", "title"); v != "" {
		out["thing1"] = v
	}
	if v := firstNonEmpty(params, "time2", "time", "operateTime", "operationTime", "notifyTime"); v != "" {
		out["time2"] = v
	} else {
		out["time2"] = formatWxSubscribeTime(time.Now())
	}
	if v := firstNonEmpty(params, "thing3", "body", "action", "sceneAction", "sceneName"); v != "" {
		out["thing3"] = v
	}
	return out
}

// buildWxMiniSubscribeData 将 params 与渲染后的正文映射为微信订阅消息 Data。
func buildWxMiniSubscribeData(params map[string]any, renderedBody string) map[string]*subscribe.DataItem {
	if params == nil {
		params = map[string]any{}
	}
	// 模板正文变量可合并进 params 供映射
	body := strings.TrimSpace(renderedBody)
	if body != "" {
		if _, ok := params["body"]; !ok {
			params["body"] = body
		}
		if _, ok := params["thing3"]; !ok && firstNonEmpty(params, "body") == "" {
			params["thing3"] = body
		}
	}
	mapped := mapWxMiniTemplateFields(params)
	// 显式传入的 thing/time 优先
	for _, k := range []string{"thing1", "time2", "thing3"} {
		if v := firstNonEmpty(params, k); v != "" {
			mapped[k] = v
		}
	}
	data := map[string]*subscribe.DataItem{}
	for k, v := range mapped {
		if v == "" {
			continue
		}
		data[k] = &subscribe.DataItem{Value: truncateWxSubscribeValue(v)}
	}
	return data
}
