package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestCreds(t *testing.T) (cleanup func()) {
	t.Helper()
	origDir := ConfigDir
	origFile := CredentialsFile

	tmp := t.TempDir()
	ConfigDir = tmp
	CredentialsFile = filepath.Join(tmp, "credentials.json")

	return func() {
		ConfigDir = origDir
		CredentialsFile = origFile
	}
}

func TestLoadCredentials_Missing(t *testing.T) {
	cleanup := setupTestCreds(t)
	defer cleanup()

	creds := LoadCredentials()
	assert.Equal(t, &Credentials{}, creds)
}

func TestLoadCredentials_Corrupt(t *testing.T) {
	cleanup := setupTestCreds(t)
	defer cleanup()

	os.WriteFile(CredentialsFile, []byte("not json{{{"), 0600)
	creds := LoadCredentials()
	assert.Equal(t, &Credentials{}, creds)
}

func TestSaveCredentials_FilePermissions(t *testing.T) {
	cleanup := setupTestCreds(t)
	defer cleanup()

	creds := &Credentials{BaseURL: "https://example.com"}
	require.NoError(t, SaveCredentials(creds))

	// Check file permissions are 0600
	info, err := os.Stat(CredentialsFile)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	cleanup := setupTestCreds(t)
	defer cleanup()

	original := &Credentials{
		BaseURL:           "https://api.example.com",
		OAuthAccessToken:  "access-tok",
		OAuthRefreshToken: "refresh-tok",
		OAuthExpiresAt:    1700000000,
		OAuthClientID:     "client-123",
	}

	require.NoError(t, SaveCredentials(original))
	loaded := LoadCredentials()
	assert.Equal(t, original, loaded)
}

func TestSaveCredentials_JSONFormat(t *testing.T) {
	cleanup := setupTestCreds(t)
	defer cleanup()

	creds := &Credentials{
		BaseURL:          "https://api.example.com",
		JWT:              "jwt-token",
		OAuthAccessToken: "oauth-token",
	}
	require.NoError(t, SaveCredentials(creds))

	data, err := os.ReadFile(CredentialsFile)
	require.NoError(t, err)

	// Should be valid JSON with expected keys
	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))
	assert.Equal(t, "https://api.example.com", raw["base_url"])
	assert.Equal(t, "jwt-token", raw["jwt"])
	assert.Equal(t, "oauth-token", raw["oauth_access_token"])
}

func TestCredentials_GetBaseURL(t *testing.T) {
	assert.Equal(t, DefaultBaseURL, (&Credentials{}).GetBaseURL())
	assert.Equal(t, "https://custom.com", (&Credentials{BaseURL: "https://custom.com"}).GetBaseURL())
}

func TestCredentials_GetToken_OAuth(t *testing.T) {
	creds := &Credentials{
		OAuthAccessToken: "oauth-tok",
		JWT:              "jwt-tok",
	}
	tok, err := creds.GetToken()
	assert.NoError(t, err)
	assert.Equal(t, "oauth-tok", tok, "OAuth should take priority over JWT")
}

func TestCredentials_GetToken_JWT(t *testing.T) {
	creds := &Credentials{JWT: "jwt-tok"}
	tok, err := creds.GetToken()
	assert.NoError(t, err)
	assert.Equal(t, "jwt-tok", tok)
}

func TestCredentials_GetToken_None(t *testing.T) {
	_, err := (&Credentials{}).GetToken()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not logged in")
}

func TestCredentials_IsOAuthTokenExpiring(t *testing.T) {
	assert.False(t, (&Credentials{}).IsOAuthTokenExpiring(), "no token")
	assert.False(t, (&Credentials{OAuthAccessToken: "tok"}).IsOAuthTokenExpiring(), "no expiry set")
	assert.True(t, (&Credentials{
		OAuthAccessToken: "tok",
		OAuthExpiresAt:   1, // way in the past
	}).IsOAuthTokenExpiring())
}
