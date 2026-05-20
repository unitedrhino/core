// Google OAuth2 客户端封装
package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	oauth2google "golang.org/x/oauth2/google"
)

// GoogleUser 表示 Google 用户信息响应
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// GoogleClient Google OAuth2 客户端
type GoogleClient struct {
	config *oauth2.Config
}

// NewGoogleClient 创建 Google OAuth2 客户端
// appID 为 Google Client ID，appSecret 为 Google Client Secret，redirectURL 为回调地址
func NewGoogleClient(appID, appSecret, redirectURL string) *GoogleClient {
	return &GoogleClient{
		config: &oauth2.Config{
			ClientID:     appID,
			ClientSecret: appSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     oauth2google.Endpoint,
		},
	}
}

// GetAuthCodeURL 生成授权链接，支持 PKCE
func (c *GoogleClient) GetAuthCodeURL(state, codeChallenge string) string {
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("prompt", "select_account"),
	}
	if codeChallenge != "" {
		opts = append(opts, oauth2.SetAuthURLParam("code_challenge", codeChallenge),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	}
	return c.config.AuthCodeURL(state, opts...)
}

// ExchangeCode 用授权码换取 Token，支持 PKCE verifier
func (c *GoogleClient) ExchangeCode(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error) {
	opts := []oauth2.AuthCodeOption{}
	if codeVerifier != "" {
		opts = append(opts, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	}
	return c.config.Exchange(ctx, code, opts...)
}

// GetUserInfo 用 access_token 获取用户信息
func (c *GoogleClient) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
	client := c.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo status: %d", resp.StatusCode)
	}
	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
