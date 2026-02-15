package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultBaseURL    = "https://alpha.lollipopai.com"
	OAuthRedirectPort = 9876
	OAuthRedirectURI  = "http://127.0.0.1:9876/callback"
)

var (
	ConfigDir       string
	CredentialsFile string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	ConfigDir = filepath.Join(home, ".cpk")
	CredentialsFile = filepath.Join(ConfigDir, "credentials.json")
}

// Credentials holds all stored auth state. JSON field names match the Python version
// so users don't need to re-authenticate after switching.
type Credentials struct {
	BaseURL           string `json:"base_url,omitempty"`
	JWT               string `json:"jwt,omitempty"`
	OAuthAccessToken  string `json:"oauth_access_token,omitempty"`
	OAuthRefreshToken string `json:"oauth_refresh_token,omitempty"`
	OAuthExpiresAt    int64  `json:"oauth_expires_at,omitempty"`
	OAuthClientID     string `json:"oauth_client_id,omitempty"`
}

// LoadCredentials reads credentials from disk. Returns zero-value on missing/corrupt file.
func LoadCredentials() *Credentials {
	data, err := os.ReadFile(CredentialsFile)
	if err != nil {
		return &Credentials{}
	}
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return &Credentials{}
	}
	return &creds
}

// SaveCredentials writes credentials to disk with 0600 permissions.
func SaveCredentials(creds *Credentials) error {
	if err := os.MkdirAll(ConfigDir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}
	fd, err := os.OpenFile(CredentialsFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = fd.Write(data)
	return err
}

// GetBaseURL returns the stored base URL or the default.
func (c *Credentials) GetBaseURL() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}
	return DefaultBaseURL
}

// GetToken returns the best available auth token.
// It prefers OAuth, then falls back to JWT.
func (c *Credentials) GetToken() (string, error) {
	if c.OAuthAccessToken != "" {
		return c.OAuthAccessToken, nil
	}
	if c.JWT != "" {
		return c.JWT, nil
	}
	return "", errors.New("not logged in. Run: cpk login")
}

// IsOAuthTokenExpiring returns true if the OAuth token is expired or within 60s of expiry.
func (c *Credentials) IsOAuthTokenExpiring() bool {
	return c.OAuthAccessToken != "" &&
		c.OAuthExpiresAt > 0 &&
		time.Now().Unix() > c.OAuthExpiresAt-60
}
