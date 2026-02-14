package twirp

import (
	"encoding/json"
	"fmt"

	"github.com/lollipopai/cli/internal/auth"
	"github.com/lollipopai/cli/internal/httpclient"
	"github.com/lollipopai/cli/internal/output"
)

// Caller makes authenticated Twirp RPC calls.
type Caller struct {
	Client *httpclient.Client
	Creds  *auth.Credentials
}

// NewCaller creates a Caller with the given HTTP client and credentials.
func NewCaller(client *httpclient.Client, creds *auth.Credentials) *Caller {
	return &Caller{Client: client, Creds: creds}
}

// Call invokes a Twirp RPC method. Auto-refreshes OAuth tokens when expiring.
// servicePath is e.g. "lollipop.proto.recipe.v1.RecipeV1", method is e.g. "Search".
func (c *Caller) Call(servicePath, method string, payload any) (any, error) {
	// Auto-refresh if token is expiring
	if c.Creds.IsOAuthTokenExpiring() {
		if err := auth.RefreshOAuthToken(c.Client, c.Creds); err != nil {
			output.Warn("OAuth token expired and refresh failed. Try: cpk login")
		}
	}

	token, err := c.Creds.GetToken()
	if err != nil {
		return nil, err
	}

	baseURL := c.Creds.GetBaseURL()
	url := fmt.Sprintf("%s/api/twirp/%s/%s", baseURL, servicePath, method)

	if payload == nil {
		payload = map[string]any{}
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	body, err := c.Client.PostJSON(url, payload, headers)
	if err != nil {
		if apiErr, ok := err.(*httpclient.APIError); ok && apiErr.StatusCode == 401 {
			return nil, fmt.Errorf("%s\nTry: cpk login", apiErr.Message)
		}
		return nil, err
	}

	var result any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON response: %w", err)
	}
	return result, nil
}
