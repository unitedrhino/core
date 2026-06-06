package defext

import "strings"

// PushClientPlatformUnbind 客户端上报 platform=unbind 时表示解绑当前用户与本机 cid，非真实 OS 平台。
const PushClientPlatformUnbind = "unbind"

// NormalizePushClientPlatform 统一客户端平台上报（鸿蒙 NEXT 常报 harmonyos，个推通道按 android/harmony 配置）。
func NormalizePushClientPlatform(platform string) string {
	p := strings.ToLower(strings.TrimSpace(platform))
	switch p {
	case PushClientPlatformUnbind:
		return PushClientPlatformUnbind
	case "harmonyos", "harmony", "ohos":
		return "harmony"
	case "android", "ios":
		return p
	default:
		if strings.Contains(p, "harmony") || strings.Contains(p, "ohos") {
			return "android"
		}
		if p == "" {
			return ""
		}
		return p
	}
}
