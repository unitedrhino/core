package relationDB

import (
	"testing"

	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestNotifyTemplatePreloadsChannelAndConfig(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&SysNotifyConfig{}, &SysNotifyChannel{}, &SysNotifyTemplate{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	config := SysNotifyConfig{
		ID:           1,
		Group:        "captcha",
		Code:         "sysUserRegisterCaptcha",
		Name:         "用户注册验证码",
		SupportTypes: []def.NotifyType{def.NotifyTypeEmail},
		EnableTypes:  []def.NotifyType{def.NotifyTypeEmail},
		Desc:         "captcha",
		IsRecord:     def.False,
		Params:       map[string]string{"code": "验证码"},
	}
	if err := db.Create(&config).Error; err != nil {
		t.Fatalf("create config: %v", err)
	}

	channel := SysNotifyChannel{
		ID:         2,
		TenantCode: dataType.TenantCode("default"),
		Type:       def.NotifyTypeEmail,
		Email: &SysTenantEmail{
			From: "noreply@example.com",
			Host: "smtp.example.com",
		},
		Name: "email channel",
		Desc: "email",
	}
	if err := db.Create(&channel).Error; err != nil {
		t.Fatalf("create channel: %v", err)
	}

	template := SysNotifyTemplate{
		ID:         8,
		TenantCode: dataType.TenantCode("default"),
		Name:       "邮箱注册",
		NotifyCode: config.Code,
		Type:       def.NotifyTypeEmail,
		Subject:    "验证码",
		Body:       "{{.code}}",
		ChannelID:  channel.ID,
	}
	if err := db.Create(&template).Error; err != nil {
		t.Fatalf("create template: %v", err)
	}

	var got SysNotifyTemplate
	if err := db.Preload("Channel").Preload("Config").Where("id = ?", template.ID).First(&got).Error; err != nil {
		t.Fatalf("find template: %v", err)
	}
	if got.Channel == nil {
		t.Fatal("expected channel to be preloaded from channel_id")
	}
	if got.Channel.Email == nil || got.Channel.Email.From != "noreply@example.com" {
		t.Fatalf("unexpected channel email: %+v", got.Channel.Email)
	}
	if got.Config == nil {
		t.Fatal("expected config to be preloaded from notify_code")
	}
	if got.Config.Code != config.Code {
		t.Fatalf("unexpected config: %+v", got.Config)
	}
}
