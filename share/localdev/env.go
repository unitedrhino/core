// env.go 提供本地开发环境保护开关。
package localdev

import (
	"os"
	"strings"
)

// SkipStartupSideEffects 返回本地开发模式是否跳过启动初始化副作用。
func SkipStartupSideEffects() bool {
	return envEnabled("CORE_LOCAL_SKIP_STARTUP_SIDE_EFFECTS")
}

// SkipAutoMigrate 返回本地开发模式是否跳过数据库自动迁移。
func SkipAutoMigrate() bool {
	return envEnabled("CORE_LOCAL_SKIP_AUTO_MIGRATE")
}

// envEnabled 判断环境变量是否显式启用。
func envEnabled(name string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(name))) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}
