// 用户管理逻辑测试。
package usermanagelogic

import (
	"context"
	"database/sql"
	"testing"

	"gitee.com/unitedrhino/core/share/clients/oauth2"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"

	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
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

func newAppleBindTestRepo(t *testing.T) *relationDB.UserInfoRepo {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&relationDB.SysUserInfo{}); err != nil {
		t.Fatalf("migrate sys_user_info: %v", err)
	}
	return relationDB.NewUserInfoRepo(db)
}

func insertAppleBindTestUser(t *testing.T, repo *relationDB.UserInfoRepo, ui *relationDB.SysUserInfo) {
	t.Helper()
	if ui.TenantCode == "" {
		ui.TenantCode = dataType.TenantCode(def.TenantCodeDefault)
	}
	if err := repo.Insert(context.Background(), ui); err != nil {
		t.Fatalf("insert user: %v", err)
	}
}

// TestFindOrBindAppleUserReturnsExistingAppleID 验证已有 Apple 用户 ID 时直接返回原用户
func TestFindOrBindAppleUserReturnsExistingAppleID(t *testing.T) {
	repo := newAppleBindTestRepo(t)
	insertAppleBindTestUser(t, repo, &relationDB.SysUserInfo{
		UserID:      1000,
		AppleUserID: sql.NullString{Valid: true, String: "apple-sub"},
		Email:       sql.NullString{Valid: true, String: "old@example.com"},
	})

	got, err := findOrBindAppleUser(context.Background(), repo, def.TenantCodeDefault, &oauth2.AppleUser{
		Sub:           "apple-sub",
		Email:         "new@example.com",
		EmailVerified: "true",
	})
	if err != nil {
		t.Fatalf("findOrBindAppleUser error: %v", err)
	}
	if got.UserID != 1000 {
		t.Fatalf("UserID = %d, want existing apple user", got.UserID)
	}
}

// TestFindOrBindAppleUserBindsVerifiedEmail 验证已验证 Apple 邮箱会绑定到同邮箱老用户
func TestFindOrBindAppleUserBindsVerifiedEmail(t *testing.T) {
	repo := newAppleBindTestRepo(t)
	insertAppleBindTestUser(t, repo, &relationDB.SysUserInfo{
		UserID: 1001,
		Email:  sql.NullString{Valid: true, String: "821465404@qq.com"},
	})

	got, err := findOrBindAppleUser(context.Background(), repo, def.TenantCodeDefault, &oauth2.AppleUser{
		Sub:           "000507.84556ee0529148f6ba9ae8ddd3678e69.0724",
		Email:         "821465404@qq.com",
		EmailVerified: "true",
	})
	if err != nil {
		t.Fatalf("findOrBindAppleUser error: %v", err)
	}
	if got.UserID != 1001 {
		t.Fatalf("UserID = %d, want existing user", got.UserID)
	}
	if !got.AppleUserID.Valid || got.AppleUserID.String != "000507.84556ee0529148f6ba9ae8ddd3678e69.0724" {
		t.Fatalf("AppleUserID = %#v, want bound apple sub", got.AppleUserID)
	}
}

// TestFindOrBindAppleUserSkipsUnverifiedEmail 验证未验证邮箱不会自动绑定旧账号
func TestFindOrBindAppleUserSkipsUnverifiedEmail(t *testing.T) {
	repo := newAppleBindTestRepo(t)
	insertAppleBindTestUser(t, repo, &relationDB.SysUserInfo{
		UserID: 1002,
		Email:  sql.NullString{Valid: true, String: "user@example.com"},
	})

	_, err := findOrBindAppleUser(context.Background(), repo, def.TenantCodeDefault, &oauth2.AppleUser{
		Sub:           "apple-sub",
		Email:         "user@example.com",
		EmailVerified: "false",
	})
	if !errors.Cmp(err, errors.NotFind) {
		t.Fatalf("err = %v, want NotFind", err)
	}
}

// TestFindOrBindAppleUserReturnsNotFoundForNewEmail 验证邮箱不存在时交给上层自动注册
func TestFindOrBindAppleUserReturnsNotFoundForNewEmail(t *testing.T) {
	repo := newAppleBindTestRepo(t)

	_, err := findOrBindAppleUser(context.Background(), repo, def.TenantCodeDefault, &oauth2.AppleUser{
		Sub:           "new-apple-sub",
		Email:         "new@example.com",
		EmailVerified: "true",
	})
	if !errors.Cmp(err, errors.NotFind) {
		t.Fatalf("err = %v, want NotFind", err)
	}
}

// TestFindOrBindAppleUserRejectsDifferentBoundApple 验证同邮箱账号已绑定其他 Apple 账号时不会覆盖
func TestFindOrBindAppleUserRejectsDifferentBoundApple(t *testing.T) {
	repo := newAppleBindTestRepo(t)
	insertAppleBindTestUser(t, repo, &relationDB.SysUserInfo{
		UserID:      1003,
		Email:       sql.NullString{Valid: true, String: "user@example.com"},
		AppleUserID: sql.NullString{Valid: true, String: "other-apple-sub"},
	})

	_, err := findOrBindAppleUser(context.Background(), repo, def.TenantCodeDefault, &oauth2.AppleUser{
		Sub:           "new-apple-sub",
		Email:         "user@example.com",
		EmailVerified: "true",
	})
	if !errors.Cmp(err, errors.BindAccount) {
		t.Fatalf("err = %v, want BindAccount", err)
	}
}

// TestFindOrBindGoogleUserReturnsExistingGoogleID 验证已有 Google 用户 ID 时直接返回原用户
func TestFindOrBindGoogleUserReturnsExistingGoogleID(t *testing.T) {
	repo := newAppleBindTestRepo(t)
	insertAppleBindTestUser(t, repo, &relationDB.SysUserInfo{
		UserID:       2000,
		GoogleUserID: sql.NullString{Valid: true, String: "google-sub"},
		Email:        sql.NullString{Valid: true, String: "old@example.com"},
	})

	got, err := findOrBindGoogleUser(context.Background(), repo, def.TenantCodeDefault, &oauth2.GoogleUser{
		ID:            "google-sub",
		Email:         "new@example.com",
		VerifiedEmail: true,
	})
	if err != nil {
		t.Fatalf("findOrBindGoogleUser error: %v", err)
	}
	if got.UserID != 2000 {
		t.Fatalf("UserID = %d, want existing google user", got.UserID)
	}
}

// TestFindOrBindGoogleUserBindsVerifiedEmail 验证已验证 Google 邮箱会绑定到同邮箱老用户
func TestFindOrBindGoogleUserBindsVerifiedEmail(t *testing.T) {
	repo := newAppleBindTestRepo(t)
	insertAppleBindTestUser(t, repo, &relationDB.SysUserInfo{
		UserID: 2001,
		Email:  sql.NullString{Valid: true, String: "821465404@qq.com"},
	})

	got, err := findOrBindGoogleUser(context.Background(), repo, def.TenantCodeDefault, &oauth2.GoogleUser{
		ID:            "google-sub-821465404",
		Email:         "821465404@qq.com",
		VerifiedEmail: true,
	})
	if err != nil {
		t.Fatalf("findOrBindGoogleUser error: %v", err)
	}
	if got.UserID != 2001 {
		t.Fatalf("UserID = %d, want existing user", got.UserID)
	}
	if !got.GoogleUserID.Valid || got.GoogleUserID.String != "google-sub-821465404" {
		t.Fatalf("GoogleUserID = %#v, want bound google sub", got.GoogleUserID)
	}
}

// TestFindOrBindGoogleUserSkipsUnverifiedEmail 验证未验证邮箱不会自动绑定旧账号
func TestFindOrBindGoogleUserSkipsUnverifiedEmail(t *testing.T) {
	repo := newAppleBindTestRepo(t)
	insertAppleBindTestUser(t, repo, &relationDB.SysUserInfo{
		UserID: 2002,
		Email:  sql.NullString{Valid: true, String: "user@example.com"},
	})

	_, err := findOrBindGoogleUser(context.Background(), repo, def.TenantCodeDefault, &oauth2.GoogleUser{
		ID:            "google-sub",
		Email:         "user@example.com",
		VerifiedEmail: false,
	})
	if !errors.Cmp(err, errors.NotFind) {
		t.Fatalf("err = %v, want NotFind", err)
	}
}

// TestFindOrBindGoogleUserReturnsNotFoundForNewEmail 验证邮箱不存在时交给上层自动注册
func TestFindOrBindGoogleUserReturnsNotFoundForNewEmail(t *testing.T) {
	repo := newAppleBindTestRepo(t)

	_, err := findOrBindGoogleUser(context.Background(), repo, def.TenantCodeDefault, &oauth2.GoogleUser{
		ID:            "new-google-sub",
		Email:         "new@example.com",
		VerifiedEmail: true,
	})
	if !errors.Cmp(err, errors.NotFind) {
		t.Fatalf("err = %v, want NotFind", err)
	}
}

// TestFindOrBindGoogleUserRejectsDifferentBoundGoogle 验证同邮箱账号已绑定其他 Google 账号时不会覆盖
func TestFindOrBindGoogleUserRejectsDifferentBoundGoogle(t *testing.T) {
	repo := newAppleBindTestRepo(t)
	insertAppleBindTestUser(t, repo, &relationDB.SysUserInfo{
		UserID:       2003,
		Email:        sql.NullString{Valid: true, String: "user@example.com"},
		GoogleUserID: sql.NullString{Valid: true, String: "other-google-sub"},
	})

	_, err := findOrBindGoogleUser(context.Background(), repo, def.TenantCodeDefault, &oauth2.GoogleUser{
		ID:            "new-google-sub",
		Email:         "user@example.com",
		VerifiedEmail: true,
	})
	if !errors.Cmp(err, errors.BindAccount) {
		t.Fatalf("err = %v, want BindAccount", err)
	}
}
