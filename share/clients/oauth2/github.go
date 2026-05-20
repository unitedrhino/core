// GitHub OAuth2 客户端封装
package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"
)

// GithubUser 表示 GitHub 用户信息响应
type GithubUser struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
}

// GithubEmail 表示 GitHub 邮箱列表项
type GithubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// GithubClient GitHub OAuth2 客户端
type GithubClient struct {
	config *oauth2.Config
}

// NewGithubClient 创建 GitHub OAuth2 客户端
func NewGithubClient(appID, appSecret, redirectURL string) *GithubClient {
	return &GithubClient{
		config: &oauth2.Config{
			ClientID:     appID,
			ClientSecret: appSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     oauth2github.Endpoint,
		},
	}
}

// GetAuthCodeURL 生成授权链接，支持 PKCE
func (c *GithubClient) GetAuthCodeURL(state, codeChallenge string) string {
	opts := []oauth2.AuthCodeOption{}
	if codeChallenge != "" {
		opts = append(opts, oauth2.SetAuthURLParam("code_challenge", codeChallenge),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	}
	return c.config.AuthCodeURL(state, opts...)
}

// ExchangeCode 用授权码换取 Token，支持 PKCE verifier
func (c *GithubClient) ExchangeCode(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error) {
	opts := []oauth2.AuthCodeOption{}
	if codeVerifier != "" {
		opts = append(opts, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	}
	return c.config.Exchange(ctx, code, opts...)
}

// GetUserInfo 用 access_token 获取用户信息，并补充 primary+verified 邮箱
func (c *GithubClient) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GithubUser, error) {
	client := c.config.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github user status: %d", resp.StatusCode)
	}
	var user GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	// 若主邮箱为空，拉取邮箱列表补充
	if user.Email == "" {
		emails, err := c.getEmails(ctx, client)
		if err == nil {
			for _, e := range emails {
				if e.Primary && e.Verified {
					user.Email = e.Email
					break
				}
			}
			// 若无 primary+verified，取第一个 verified
			if user.Email == "" {
				for _, e := range emails {
					if e.Verified {
						user.Email = e.Email
						break
					}
				}
			}
		}
	}
	return &user, nil
}

func (c *GithubClient) getEmails(ctx context.Context, client *http.Client) ([]GithubEmail, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github emails status: %d", resp.StatusCode)
	}
	var emails []GithubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, err
	}
	return emails, nil
}
