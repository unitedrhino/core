// Apple OAuth2 客户端封装
package oauth2

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AppleUser 表示 Apple ID Token 解析后的用户信息
type AppleUser struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"` // 可能是 "true"/"false" 或 bool
	IsPrivateEmail string `json:"is_private_email"`
}

// AppleClient Apple OAuth2 客户端
type AppleClient struct {
	clientID    string // Bundle ID / Services ID
	teamID      string
	keyID       string
	privateKey  *ecdsa.PrivateKey
	redirectURI string
}

// normalizeApplePrivateKeyPEM 将裸 Base64 或完整 PEM/.p8 内容规范为 PEM 格式
func normalizeApplePrivateKeyPEM(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.Contains(raw, "BEGIN") {
		return raw
	}
	var b64 strings.Builder
	for _, r := range raw {
		if r != '\n' && r != '\r' && r != ' ' {
			b64.WriteRune(r)
		}
	}
	s := b64.String()
	var lines []string
	for i := 0; i < len(s); i += 64 {
		end := i + 64
		if end > len(s) {
			end = len(s)
		}
		lines = append(lines, s[i:end])
	}
	return "-----BEGIN PRIVATE KEY-----\n" + strings.Join(lines, "\n") + "\n-----END PRIVATE KEY-----\n"
}

// NewAppleClient 创建 Apple OAuth2 客户端
// privateKeyPEM 为 Apple 提供的 .p8 文件内容（PEM 格式，也支持仅粘贴 Base64 主体）
func NewAppleClient(clientID, teamID, keyID, privateKeyPEM, redirectURI string) (*AppleClient, error) {
	privateKeyPEM = normalizeApplePrivateKeyPEM(privateKeyPEM)
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("apple private key pem decode failed")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("apple private key parse failed: %w", err)
	}
	ecKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("apple private key is not ECDSA")
	}
	return &AppleClient{
		clientID:    clientID,
		teamID:      teamID,
		keyID:       keyID,
		privateKey:  ecKey,
		redirectURI: redirectURI,
	}, nil
}

// GetAuthCodeURL 生成 Apple 授权链接
func (c *AppleClient) GetAuthCodeURL(state string) string {
	v := url.Values{}
	v.Set("client_id", c.clientID)
	v.Set("redirect_uri", c.redirectURI)
	v.Set("response_type", "code")
	v.Set("scope", "name email")
	v.Set("response_mode", "query")
	v.Set("state", state)
	return "https://appleid.apple.com/auth/authorize?" + v.Encode()
}

// exchangeCodePost 向 Apple 换取 token；includeRedirect 为 false 时不传 redirect_uri（原生 App 场景）
func (c *AppleClient) exchangeCodePost(ctx context.Context, code, clientSecret string, includeRedirect bool) (*AppleUser, string, error) {
	data := url.Values{}
	data.Set("client_id", c.clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	if includeRedirect && c.redirectURI != "" {
		data.Set("redirect_uri", c.redirectURI)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://appleid.apple.com/auth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("apple token status: %d body: %s", resp.StatusCode, string(body))
	}
	var result struct {
		IDToken string `json:"id_token"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, "", err
	}
	user, err := c.parseIDToken(result.IDToken)
	if err != nil {
		return nil, "", err
	}
	return user, result.IDToken, nil
}

// ExchangeCode 用授权码换取 Token，返回 ID Token 和用户信息
func (c *AppleClient) ExchangeCode(ctx context.Context, code string) (*AppleUser, string, error) {
	clientSecret, err := c.buildClientSecret()
	if err != nil {
		return nil, "", err
	}
	// 先按 Web 配置（含 redirect_uri）请求；400 时再尝试不传 redirect_uri（iOS/Android 原生）
	user, idToken, err := c.exchangeCodePost(ctx, code, clientSecret, true)
	if err == nil {
		return user, idToken, nil
	}
	if c.redirectURI != "" {
		user2, idToken2, err2 := c.exchangeCodePost(ctx, code, clientSecret, false)
		if err2 == nil {
			return user2, idToken2, nil
		}
	}
	return nil, "", err
}

// buildClientSecret 生成 Apple Client Secret（ES256 JWT，有效期 5 分钟）
func (c *AppleClient) buildClientSecret() (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": c.teamID,
		"iat": now.Unix(),
		"exp": now.Add(5 * time.Minute).Unix(),
		"aud": "https://appleid.apple.com",
		"sub": c.clientID,
	})
	token.Header["kid"] = c.keyID
	return token.SignedString(c.privateKey)
}

// parseIDToken 解析 Apple ID Token（不验证签名，依赖 HTTPS 安全）
func (c *AppleClient) parseIDToken(idToken string) (*AppleUser, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id_token format")
	}
	payload, err := jwt.NewParser().DecodeSegment(parts[1])
	if err != nil {
		return nil, err
	}
	var claims jwt.MapClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}
	user := &AppleUser{Sub: getStringClaim(claims, "sub")}
	if email, ok := claims["email"].(string); ok {
		user.Email = email
	}
	if ev, ok := claims["email_verified"]; ok {
		switch v := ev.(type) {
		case string:
			user.EmailVerified = v
		case bool:
			if v {
				user.EmailVerified = "true"
			} else {
				user.EmailVerified = "false"
			}
		}
	}
	if ip, ok := claims["is_private_email"]; ok {
		switch v := ip.(type) {
		case string:
			user.IsPrivateEmail = v
		case bool:
			if v {
				user.IsPrivateEmail = "true"
			} else {
				user.IsPrivateEmail = "false"
			}
		}
	}
	return user, nil
}

func getStringClaim(claims jwt.MapClaims, key string) string {
	if v, ok := claims[key].(string); ok {
		return v
	}
	return ""
}
