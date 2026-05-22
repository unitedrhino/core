package usermanagelogic

import (
	"database/sql"
	"testing"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
)

// TestApplyOAuthLoginAccountFillsAccountAndNicknameFromEmail 验证 Apple OAuth 仅返回邮箱时能补齐账号和昵称
func TestApplyOAuthLoginAccountFillsAccountAndNicknameFromEmail(t *testing.T) {
	ui := &relationDB.SysUserInfo{
		Email: sql.NullString{Valid: true, String: "18059688688@163.com"},
	}

	applyOAuthLoginAccount(ui)

	if !ui.UserName.Valid || ui.UserName.String != "18059688688@163.com" {
		t.Fatalf("UserName = %#v, want email", ui.UserName)
	}
	if ui.NickName != "18059688688@163.com" {
		t.Fatalf("NickName = %q, want email", ui.NickName)
	}
}
