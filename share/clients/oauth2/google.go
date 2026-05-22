// Google OAuth2 客户端封装
package oauth2

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	config          *oauth2.Config
	jwksURL         string
	jwksHeaderName  string
	jwksAccessToken string
	jwksMu          sync.RWMutex
	jwksKeys        map[string]*rsa.PublicKey
	jwksExpiresAt   time.Time
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

// SetJWKSFetchConfig 设置 Google ID Token 离线验签所需的 JWKS 获取参数
func (c *GoogleClient) SetJWKSFetchConfig(url, headerName, accessToken string) {
	c.jwksURL = strings.TrimSpace(url)
	c.jwksHeaderName = strings.TrimSpace(headerName)
	c.jwksAccessToken = strings.TrimSpace(accessToken)
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

// ResolveUser 用授权码或 ID Token 获取 Google 用户信息
func (c *GoogleClient) ResolveUser(ctx context.Context, codeOrIDToken string) (*GoogleUser, error) {
	if isJWTLike(codeOrIDToken) {
		return c.GetUserInfoByIDToken(ctx, codeOrIDToken)
	}
	token, err := c.ExchangeCode(ctx, codeOrIDToken, "")
	if err != nil {
		return nil, err
	}
	return c.GetUserInfo(ctx, token)
}

// GetUserInfoByIDToken 离线校验 Google ID Token 并提取用户信息
func (c *GoogleClient) GetUserInfoByIDToken(ctx context.Context, idToken string) (*GoogleUser, error) {
	if c.jwksURL == "" {
		return nil, fmt.Errorf("google jwks url is empty")
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(idToken, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, fmt.Errorf("unexpected google jwt alg: %s", token.Method.Alg())
		}
		kid, _ := token.Header["kid"].(string)
		if kid == "" {
			return nil, fmt.Errorf("google jwt missing kid")
		}
		return c.getJWKSPublicKey(ctx, kid)
	}, jwt.WithAudience(c.config.ClientID), jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}
	if token == nil || !token.Valid {
		return nil, fmt.Errorf("invalid google id token")
	}
	iss := getStringClaim(claims, "iss")
	if iss != "accounts.google.com" && iss != "https://accounts.google.com" {
		return nil, fmt.Errorf("invalid google issuer: %s", iss)
	}
	sub := getStringClaim(claims, "sub")
	if sub == "" {
		return nil, fmt.Errorf("google id token missing sub")
	}
	return &GoogleUser{
		ID:            sub,
		Email:         getStringClaim(claims, "email"),
		VerifiedEmail: getBoolClaim(claims, "email_verified"),
		Name:          getStringClaim(claims, "name"),
		Picture:       getStringClaim(claims, "picture"),
	}, nil
}

func isJWTLike(value string) bool {
	return strings.Count(value, ".") == 2
}

func (c *GoogleClient) getJWKSPublicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	c.jwksMu.RLock()
	if key := c.jwksKeys[kid]; key != nil && time.Now().Before(c.jwksExpiresAt) {
		c.jwksMu.RUnlock()
		return key, nil
	}
	c.jwksMu.RUnlock()
	if err := c.refreshJWKS(ctx); err != nil {
		return nil, err
	}
	c.jwksMu.RLock()
	defer c.jwksMu.RUnlock()
	if key := c.jwksKeys[kid]; key != nil {
		return key, nil
	}
	return nil, fmt.Errorf("google jwks missing kid: %s", kid)
}

func (c *GoogleClient) refreshJWKS(ctx context.Context) error {
	c.jwksMu.Lock()
	defer c.jwksMu.Unlock()
	if time.Now().Before(c.jwksExpiresAt) && len(c.jwksKeys) > 0 {
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.jwksURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if c.jwksHeaderName != "" && c.jwksAccessToken != "" {
		req.Header.Set(c.jwksHeaderName, c.jwksAccessToken)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("google jwks status: %d body: %s", resp.StatusCode, string(body))
	}
	var jwks struct {
		Keys []struct {
			Kty string `json:"kty"`
			Use string `json:"use"`
			Kid string `json:"kid"`
			Alg string `json:"alg"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.Unmarshal(body, &jwks); err != nil {
		return err
	}
	keys := make(map[string]*rsa.PublicKey, len(jwks.Keys))
	for _, key := range jwks.Keys {
		pub, err := rsaPublicKeyFromJWK(key.Kty, key.Kid, key.N, key.E)
		if err != nil {
			return err
		}
		keys[key.Kid] = pub
	}
	if len(keys) == 0 {
		return fmt.Errorf("google jwks has no keys")
	}
	c.jwksKeys = keys
	c.jwksExpiresAt = time.Now().Add(parseCacheMaxAge(resp.Header.Get("Cache-Control")))
	return nil
}

func rsaPublicKeyFromJWK(kty, kid, n, e string) (*rsa.PublicKey, error) {
	if kty != "RSA" || kid == "" || n == "" || e == "" {
		return nil, fmt.Errorf("invalid google jwk kid=%s", kid)
	}
	nb, err := base64.RawURLEncoding.DecodeString(n)
	if err != nil {
		return nil, err
	}
	eb, err := base64.RawURLEncoding.DecodeString(e)
	if err != nil {
		return nil, err
	}
	exponent := 0
	for _, b := range eb {
		exponent = exponent<<8 + int(b)
	}
	if exponent == 0 {
		return nil, fmt.Errorf("invalid google jwk exponent kid=%s", kid)
	}
	return &rsa.PublicKey{N: new(big.Int).SetBytes(nb), E: exponent}, nil
}

func parseCacheMaxAge(cacheControl string) time.Duration {
	const fallback = 6 * time.Hour
	for _, part := range strings.Split(cacheControl, ",") {
		item := strings.TrimSpace(strings.ToLower(part))
		if strings.HasPrefix(item, "max-age=") {
			seconds, err := time.ParseDuration(strings.TrimPrefix(item, "max-age=") + "s")
			if err == nil && seconds > 0 {
				return seconds
			}
		}
	}
	return fallback
}

func getBoolClaim(claims jwt.MapClaims, key string) bool {
	switch v := claims[key].(type) {
	case bool:
		return v
	case string:
		return v == "true"
	default:
		return false
	}
}
