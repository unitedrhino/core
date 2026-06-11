package notifymanagelogic

import (
	"testing"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/def"
)

func TestEmailConfigRequiresConfiguredChannel(t *testing.T) {
	tests := []struct {
		name    string
		channel *relationDB.SysNotifyChannel
	}{
		{
			name:    "missing channel",
			channel: nil,
		},
		{
			name:    "missing email config",
			channel: &relationDB.SysNotifyChannel{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := emailConfigFromChannel(tt.channel); err == nil {
				t.Fatal("expected an error for unconfigured email channel")
			}
		})
	}
}

func TestEmailConfigFromChannelMapsEmailFields(t *testing.T) {
	cfg, err := emailConfigFromChannel(&relationDB.SysNotifyChannel{
		Email: &relationDB.SysTenantEmail{
			From:     "noreply@example.com",
			Host:     "smtp.example.com",
			Secret:   "secret",
			Nickname: "YK",
			Port:     465,
			IsSSL:    def.True,
		},
	})
	if err != nil {
		t.Fatalf("expected configured email channel, got error: %v", err)
	}

	if cfg.From != "noreply@example.com" ||
		cfg.Host != "smtp.example.com" ||
		cfg.Secret != "secret" ||
		cfg.Nickname != "YK" ||
		cfg.Port != 465 ||
		!cfg.IsSSL {
		t.Fatalf("unexpected email config: %+v", cfg)
	}
}
