package auth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/lollipopai/cli/internal/httpclient"
	"github.com/lollipopai/cli/internal/output"
)

// RefreshOAuthToken attempts to refresh an OAuth token. Updates and saves creds on success.
// Returns error (not fatal) so callers can handle gracefully.
func RefreshOAuthToken(client *httpclient.Client, creds *Credentials) error {
	if creds.OAuthRefreshToken == "" || creds.OAuthClientID == "" {
		return fmt.Errorf("no refresh token or client ID")
	}

	baseURL := creds.GetBaseURL()

	// Discover token endpoint
	tokenEndpoint := baseURL + "/oauth/token"
	body, err := client.GetJSON(baseURL+"/.well-known/oauth-authorization-server", nil)
	if err == nil {
		var meta map[string]any
		if json.Unmarshal(body, &meta) == nil {
			if ep, ok := meta["token_endpoint"].(string); ok {
				tokenEndpoint = ep
			}
		}
	}

	params := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {creds.OAuthRefreshToken},
		"client_id":     {creds.OAuthClientID},
	}

	body, err = client.PostForm(tokenEndpoint, params, nil)
	if err != nil {
		return fmt.Errorf("token refresh failed: %w", err)
	}

	var resp map[string]any
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("invalid refresh response: %w", err)
	}

	if at, ok := resp["access_token"].(string); ok {
		creds.OAuthAccessToken = at
	}
	if rt, ok := resp["refresh_token"].(string); ok {
		creds.OAuthRefreshToken = rt
	}
	if ei, ok := resp["expires_in"].(float64); ok {
		creds.OAuthExpiresAt = time.Now().Unix() + int64(ei)
	}

	if err := SaveCredentials(creds); err != nil {
		return fmt.Errorf("failed to save refreshed credentials: %w", err)
	}

	output.Info("OAuth token refreshed.")
	return nil
}
