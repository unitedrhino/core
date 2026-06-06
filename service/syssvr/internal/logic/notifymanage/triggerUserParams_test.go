package notifymanagelogic

import "testing"

func TestFormatTriggerUserLine(t *testing.T) {
	auto := FormatTriggerUserLine(triggerUserFields{Type: triggerTypeAuto})
	if auto != "系统自动触发" {
		t.Fatalf("auto=%q", auto)
	}
	autoUser := FormatTriggerUserLine(triggerUserFields{Type: triggerTypeAuto, Nick: "李四"})
	if autoUser != "由 李四 触发" {
		t.Fatalf("autoUser=%q", autoUser)
	}
	manual := FormatTriggerUserLine(triggerUserFields{Type: triggerTypeManual, Nick: "张三"})
	if manual != "由 张三 手动触发" {
		t.Fatalf("manual=%q", manual)
	}
}

func TestFormatSceneButtonTriggerLine(t *testing.T) {
	got := FormatSceneButtonTriggerLine(triggerUserFields{
		Type:    triggerTypeDeviceButton,
		Account: "13800138000",
	}, "按钮1")
	want := "13800138000触发了按钮1"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	gotNick := FormatSceneButtonTriggerLine(triggerUserFields{
		Type:    triggerTypeDeviceButton,
		Nick:    "超级",
		Account: "13800138000",
	}, "遥控器1 开启")
	wantNick := "13800138000触发了遥控器1 开启"
	if gotNick != wantNick {
		t.Fatalf("prefer account: got %q want %q", gotNick, wantNick)
	}
}

func TestCompactRuleSceneActionBody(t *testing.T) {
	raw := "1 遥控器1 开启\nWIFI数字遥控1 遥控器1 开启"
	got := compactRuleSceneActionBody(raw)
	want := "WIFI数字遥控1 遥控器1 开启"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestPrependTriggerUserLine(t *testing.T) {
	got := PrependTriggerUserLine("设备已联动", triggerUserFields{Type: triggerTypeManual, Nick: "张三"})
	want := "由 张三 手动触发\n设备已联动"
	if got != want {
		t.Fatalf("prepend=%q", got)
	}
	legacy := PrependTriggerUserLine("设备已联动\n由 张三 手动触发", triggerUserFields{Type: triggerTypeManual, Nick: "张三"})
	if legacy != want {
		t.Fatalf("legacy=%q", legacy)
	}
}
