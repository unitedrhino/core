package notifymanagelogic

import (
	"testing"
)

func TestBuildWxMiniSubscribeData_wxTemplateFields(t *testing.T) {
	data := buildWxMiniSubscribeData(map[string]any{
		"deviceAlias": "客厅灯",
		"body":        "已打开",
		"time2":       "2026年05月30日 15:00",
	}, "")
	if data["thing1"] == nil || data["thing1"].Value != "客厅灯" {
		t.Fatalf("expected thing1=deviceAlias, got %+v", data["thing1"])
	}
	if data["time2"] == nil || data["time2"].Value != "2026年05月30日 15:00" {
		t.Fatalf("expected time2")
	}
	if data["thing3"] == nil || data["thing3"].Value != "已打开" {
		t.Fatalf("expected thing3=body")
	}
}

func TestBuildWxMiniSubscribeData_customKeys(t *testing.T) {
	data := buildWxMiniSubscribeData(map[string]any{"thing3": "场景执行"}, "")
	if data["thing3"] == nil || data["thing3"].Value != "场景执行" {
		t.Fatalf("expected thing3")
	}
	if data["time2"] == nil || data["time2"].Value == "" {
		t.Fatalf("expected default time2")
	}
}
