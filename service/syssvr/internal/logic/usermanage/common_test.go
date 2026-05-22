// 用户管理逻辑测试。
package usermanagelogic

import (
	"context"
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

// TestUserInfoToPbKeepsExternalHeadImgURL 验证第三方头像 URL 不会被包装成本地 OSS 下载地址
func TestUserInfoToPbKeepsExternalHeadImgURL(t *testing.T) {
	const headImg = "https://lh3.googleusercontent.com/a/avatar=s96-c"
	ui := &relationDB.SysUserInfo{HeadImg: headImg}

	got := UserInfoToPb(context.Background(), ui, nil)

	if got.HeadImg != headImg {
		t.Fatalf("HeadImg = %q, want external URL unchanged", got.HeadImg)
	}
}

// TestLocalizedRegisterProjectName 验证注册自动项目默认名会按请求语言切换
func TestLocalizedRegisterProjectName(t *testing.T) {
	tests := []struct {
		name           string
		projectName    string
		acceptLanguage string
		want           string
	}{
		{
			name:           "english default project name",
			projectName:    "我的物联",
			acceptLanguage: "en-US,en;q=0.9",
			want:           "My Home",
		},
		{
			name:           "chinese default project name",
			projectName:    "我的物联",
			acceptLanguage: "zh-CN,zh;q=0.9",
			want:           "我的物联",
		},
		{
			name:           "empty language keeps current default",
			projectName:    "我的物联",
			acceptLanguage: "",
			want:           "我的物联",
		},
		{
			name:           "custom project name is preserved",
			projectName:    "客户自定义",
			acceptLanguage: "en-US,en;q=0.9",
			want:           "客户自定义",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := localizedRegisterProjectName(tt.projectName, tt.acceptLanguage)
			if got != tt.want {
				t.Fatalf("localizedRegisterProjectName(%q, %q) = %q, want %q",
					tt.projectName, tt.acceptLanguage, got, tt.want)
			}
		})
	}
}
