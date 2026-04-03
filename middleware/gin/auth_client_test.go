package gin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthClient_Defaults(t *testing.T) {
	client := NewAuthClient(AuthClientConfig{})

	require.NotNil(t, client)
	assert.Equal(t, "edugo-central", client.config.JWTIssuer)
	assert.Equal(t, 5*time.Second, client.httpClient.Timeout)
	assert.Equal(t, 60*time.Second, client.cache.ttl)
	assert.Nil(t, client.jwtManager)
}

func TestTokenCache_SetAndGet(t *testing.T) {
	cache := newTokenCache(2 * time.Second)
	info := &TokenInfo{Valid: true, UserID: "u-1"}

	cache.Set("k1", info)
	got, ok := cache.Get("k1")
	require.True(t, ok)
	assert.Equal(t, "u-1", got.UserID)

	cache.entries["expired"] = &cacheEntry{info: info, expiresAt: time.Now().Add(-time.Second)}
	got, ok = cache.Get("expired")
	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestValidateToken_LocalSuccess(t *testing.T) {
	//nolint:gosec // Secreto de prueba para tokens en tests unitarios.
	cfg := AuthClientConfig{
		JWTSecret:       "test-secret-32-chars-minimum-123456",
		JWTIssuer:       "edugo-central",
		CacheEnabled:    true,
		FallbackEnabled: false,
	}
	client := NewAuthClient(cfg)

	ctx := &auth.UserContext{RoleID: "r1", RoleName: "admin", Permissions: []string{"users:read"}}
	token, _, err := client.jwtManager.GenerateTokenWithContext("user-1", "user@test.com", ctx, 10*time.Minute)
	require.NoError(t, err)

	info, err := client.ValidateToken(context.Background(), token)
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.True(t, info.Valid)
	assert.Equal(t, "user-1", info.UserID)
	assert.Equal(t, "user@test.com", info.Email)

	cacheKey := hashToken(token)
	cached, found := client.cache.Get(cacheKey)
	require.True(t, found)
	assert.True(t, cached.Valid)
}

func TestValidateToken_LocalFailNoFallback(t *testing.T) {
	//nolint:gosec // Secreto de prueba para tokens en tests unitarios.
	cfg := AuthClientConfig{
		JWTSecret:       "test-secret-32-chars-minimum-123456",
		JWTIssuer:       "edugo-central",
		FallbackEnabled: false,
		RemoteEnabled:   false,
	}
	client := NewAuthClient(cfg)

	info, err := client.ValidateToken(context.Background(), "invalid-token")
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.False(t, info.Valid)
	assert.NotEmpty(t, info.Error)
}

func TestValidateToken_NoValidationMethod(t *testing.T) {
	client := NewAuthClient(AuthClientConfig{})

	info, err := client.ValidateToken(context.Background(), "any-token")
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.False(t, info.Valid)
	assert.Equal(t, "no validation method available", info.Error)
}

func TestValidateToken_RemoteSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/auth/verify", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(TokenInfo{Valid: true, UserID: "remote-user", Email: "remote@test.com"}))
	}))
	defer server.Close()

	client := NewAuthClient(AuthClientConfig{
		BaseURL:       server.URL,
		RemoteEnabled: true,
		CacheEnabled:  true,
	})

	info, err := client.ValidateToken(context.Background(), "remote-token")
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.True(t, info.Valid)
	assert.Equal(t, "remote-user", info.UserID)

	cacheKey := hashToken("remote-token")
	cached, found := client.cache.Get(cacheKey)
	require.True(t, found)
	assert.Equal(t, "remote-user", cached.UserID)
}

func TestValidateToken_RemoteServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error":"boom"}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewAuthClient(AuthClientConfig{
		BaseURL:       server.URL,
		RemoteEnabled: true,
	})

	info, err := client.ValidateToken(context.Background(), "remote-token")
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.False(t, info.Valid)
	assert.True(t, strings.Contains(info.Error, "remote validation failed"))
}

func TestValidateToken_LocalFallbackToRemote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"valid":true,"user_id":"fallback-user","email":"fb@test.com"}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	//nolint:gosec // Secreto de prueba para tokens en tests unitarios.
	client := NewAuthClient(AuthClientConfig{
		JWTSecret:       "test-secret-32-chars-minimum-123456",
		JWTIssuer:       "edugo-central",
		BaseURL:         server.URL,
		RemoteEnabled:   true,
		FallbackEnabled: true,
	})

	info, err := client.ValidateToken(context.Background(), "invalid-local-token")
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.True(t, info.Valid)
	assert.Equal(t, "fallback-user", info.UserID)
}

func TestValidateToken_RemoteInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{invalid-json`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewAuthClient(AuthClientConfig{
		BaseURL:       server.URL,
		RemoteEnabled: true,
	})

	info, err := client.ValidateToken(context.Background(), "remote-token")
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.False(t, info.Valid)
	assert.True(t, strings.Contains(info.Error, "remote validation failed"))
}
