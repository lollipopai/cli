package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/lollipopai/cli/internal/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshOAuthToken_Success(t *testing.T) {
	cleanup := setupTestCreds(t)
	defer cleanup()

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, "refresh_token", r.FormValue("grant_type"))
		assert.Equal(t, "old-refresh", r.FormValue("refresh_token"))
		assert.Equal(t, "client-123", r.FormValue("client_id"))

		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "new-access",
			"refresh_token": "new-refresh",
			"expires_in":    7200,
		})
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	// No .well-known → falls back to {baseURL}/oauth/token
	creds := &Credentials{
		BaseURL:           srv.URL,
		OAuthAccessToken:  "old-access",
		OAuthRefreshToken: "old-refresh",
		OAuthClientID:     "client-123",
	}

	client := httpclient.New()
	err := RefreshOAuthToken(client, creds)
	require.NoError(t, err)

	assert.Equal(t, "new-access", creds.OAuthAccessToken)
	assert.Equal(t, "new-refresh", creds.OAuthRefreshToken)
	assert.Greater(t, creds.OAuthExpiresAt, int64(0))

	// Verify it was saved to disk
	loaded := LoadCredentials()
	assert.Equal(t, "new-access", loaded.OAuthAccessToken)
}

func TestRefreshOAuthToken_NoRefreshToken(t *testing.T) {
	creds := &Credentials{OAuthClientID: "client-123"}
	client := httpclient.New()
	err := RefreshOAuthToken(client, creds)
	assert.Error(t, err)
}

func TestRefreshOAuthToken_NoClientID(t *testing.T) {
	creds := &Credentials{OAuthRefreshToken: "refresh-tok"}
	client := httpclient.New()
	err := RefreshOAuthToken(client, creds)
	assert.Error(t, err)
}

func TestRefreshOAuthToken_ServerError(t *testing.T) {
	cleanup := setupTestCreds(t)
	defer cleanup()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"error":"invalid_grant"}`))
	}))
	defer srv.Close()

	creds := &Credentials{
		BaseURL:           srv.URL,
		OAuthRefreshToken: "bad-refresh",
		OAuthClientID:     "client-123",
	}

	client := httpclient.New()
	err := RefreshOAuthToken(client, creds)
	assert.Error(t, err, "refresh should return error, not fatal")
}

func TestRefreshOAuthToken_SavesCredentials(t *testing.T) {
	tmp := t.TempDir()
	origDir := ConfigDir
	origFile := CredentialsFile
	ConfigDir = tmp
	CredentialsFile = filepath.Join(tmp, "credentials.json")
	defer func() {
		ConfigDir = origDir
		CredentialsFile = origFile
	}()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/.well-known/oauth-authorization-server" {
			w.WriteHeader(404)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"access_token": "saved-tok",
			"expires_in":   3600,
		})
	}))
	defer srv.Close()

	creds := &Credentials{
		BaseURL:           srv.URL,
		OAuthRefreshToken: "refresh-tok",
		OAuthClientID:     "client-id",
	}

	client := httpclient.New()
	require.NoError(t, RefreshOAuthToken(client, creds))

	loaded := LoadCredentials()
	assert.Equal(t, "saved-tok", loaded.OAuthAccessToken)
}
