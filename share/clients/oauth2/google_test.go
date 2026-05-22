// Google OAuth2 客户端测试
package oauth2

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func testGoogleIDToken(t *testing.T, key *rsa.PrivateKey, kid string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":            "https://accounts.google.com",
		"aud":            "test-client-id",
		"exp":            time.Now().Add(time.Hour).Unix(),
		"iat":            time.Now().Unix(),
		"sub":            "google-sub-1",
		"email":          "user@example.com",
		"email_verified": true,
		"name":           "Test User",
		"picture":        "https://example.com/avatar.png",
	})
	token.Header["kid"] = kid
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return signed
}

func testGoogleJWKS(key *rsa.PrivateKey, kid string) []byte {
	jwks := map[string]any{
		"keys": []map[string]string{
			{
				"kty": "RSA",
				"use": "sig",
				"kid": kid,
				"alg": "RS256",
				"n":   base64.RawURLEncoding.EncodeToString(key.PublicKey.N.Bytes()),
				"e":   base64.RawURLEncoding.EncodeToString([]byte{0x01, 0x00, 0x01}),
			},
		},
	}
	body, _ := json.Marshal(jwks)
	return body
}

func TestGoogleClientGetUserInfoByIDTokenUsesJWKSCache(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	const kid = "test-kid"
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Header.Get("X-YKHL-JWKS-TOKEN") != "secret" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=3600")
		_, _ = w.Write(testGoogleJWKS(key, kid))
	}))
	defer srv.Close()

	cli := NewGoogleClient("test-client-id", "secret", "")
	cli.SetJWKSFetchConfig(srv.URL, "X-YKHL-JWKS-TOKEN", "secret")
	idToken := testGoogleIDToken(t, key, kid)

	user, err := cli.GetUserInfoByIDToken(context.Background(), idToken)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID != "google-sub-1" || user.Email != "user@example.com" || !user.VerifiedEmail {
		t.Fatalf("unexpected user: %#v", user)
	}
	if _, err := cli.GetUserInfoByIDToken(context.Background(), idToken); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("jwks cache miss, calls=%d", calls)
	}
}

func TestGoogleClientGetUserInfoByIDTokenRejectsWrongAudience(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	const kid = "test-kid"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600")
		_, _ = w.Write(testGoogleJWKS(key, kid))
	}))
	defer srv.Close()

	cli := NewGoogleClient("other-client-id", "secret", "")
	cli.SetJWKSFetchConfig(srv.URL, "", "")
	if _, err := cli.GetUserInfoByIDToken(context.Background(), testGoogleIDToken(t, key, kid)); err == nil {
		t.Fatal("expected audience mismatch")
	}
}
