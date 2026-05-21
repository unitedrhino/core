// env_test.go 提供本地开发保护开关的单元测试。
package localdev

import "testing"

func TestSkipStartupSideEffectsReadsTruthyEnv(t *testing.T) {
	t.Setenv("CORE_LOCAL_SKIP_STARTUP_SIDE_EFFECTS", "1")

	if !SkipStartupSideEffects() {
		t.Fatalf("SkipStartupSideEffects() = false, want true")
	}
}

func TestSkipAutoMigrateReadsTruthyEnv(t *testing.T) {
	t.Setenv("CORE_LOCAL_SKIP_AUTO_MIGRATE", "true")

	if !SkipAutoMigrate() {
		t.Fatalf("SkipAutoMigrate() = false, want true")
	}
}

func TestProtectionFlagsDefaultToFalse(t *testing.T) {
	t.Setenv("CORE_LOCAL_SKIP_STARTUP_SIDE_EFFECTS", "")
	t.Setenv("CORE_LOCAL_SKIP_AUTO_MIGRATE", "")

	if SkipStartupSideEffects() {
		t.Fatalf("SkipStartupSideEffects() = true, want false")
	}
	if SkipAutoMigrate() {
		t.Fatalf("SkipAutoMigrate() = true, want false")
	}
}

func TestProtectionFlagsRejectFalseLikeValues(t *testing.T) {
	t.Setenv("CORE_LOCAL_SKIP_STARTUP_SIDE_EFFECTS", "false")
	t.Setenv("CORE_LOCAL_SKIP_AUTO_MIGRATE", "0")

	if SkipStartupSideEffects() {
		t.Fatalf("SkipStartupSideEffects() = true, want false")
	}
	if SkipAutoMigrate() {
		t.Fatalf("SkipAutoMigrate() = true, want false")
	}
}
