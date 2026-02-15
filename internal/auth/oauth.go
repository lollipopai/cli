package auth

import (
	"encoding/json"
	"fmt"
	"html"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/lollipopai/cli/internal/httpclient"
	"github.com/lollipopai/cli/internal/output"
)

// OAuthConfig holds discovered authorization server endpoints.
type OAuthConfig struct {
	AuthorizationEndpoint string
	TokenEndpoint         string
	RegistrationEndpoint  string
	ScopesSupported       []string
}

// CallbackResult is returned by the local OAuth callback server.
type CallbackResult struct {
	Code             string
	State            string
	Error            string
	ErrorDescription string
}

// DiscoverOAuthConfig performs two-step .well-known discovery.
func DiscoverOAuthConfig(client *httpclient.Client, baseURL string) (*OAuthConfig, error) {
	// Step 1: Try protected resource metadata
	authServerURL := baseURL
	body, err := client.GetJSON(baseURL+"/.well-known/oauth-protected-resource", nil)
	if err == nil {
		var prMeta map[string]any
		if json.Unmarshal(body, &prMeta) == nil {
			if servers, ok := prMeta["authorization_servers"].([]any); ok && len(servers) > 0 {
				if s, ok := servers[0].(string); ok {
					authServerURL = s
				}
			}
		}
	} else {
		output.Warn("Could not fetch protected resource metadata, using base URL as auth server.")
	}

	// Step 2: Fetch authorization server metadata
	body, err = client.GetJSON(authServerURL+"/.well-known/oauth-authorization-server", nil)
	if err != nil {
		return nil, fmt.Errorf("could not fetch authorization server metadata: %w", err)
	}

	var asMeta map[string]any
	if err := json.Unmarshal(body, &asMeta); err != nil {
		return nil, fmt.Errorf("invalid authorization server metadata: %w", err)
	}

	config := &OAuthConfig{
		AuthorizationEndpoint: stringOrDefault(asMeta, "authorization_endpoint", authServerURL+"/oauth/authorize"),
		TokenEndpoint:         stringOrDefault(asMeta, "token_endpoint", authServerURL+"/oauth/token"),
		RegistrationEndpoint:  stringOrDefault(asMeta, "registration_endpoint", authServerURL+"/oauth/register"),
	}

	if scopes, ok := asMeta["scopes_supported"].([]any); ok {
		for _, s := range scopes {
			if str, ok := s.(string); ok {
				config.ScopesSupported = append(config.ScopesSupported, str)
			}
		}
	}
	if len(config.ScopesSupported) == 0 {
		config.ScopesSupported = []string{"read", "write"}
	}

	return config, nil
}

// RegisterClient performs dynamic client registration.
func RegisterClient(client *httpclient.Client, registrationEndpoint string) (string, error) {
	payload := map[string]any{
		"client_name":                "Cherrypick CLI",
		"redirect_uris":             []string{OAuthRedirectURI},
		"token_endpoint_auth_method": "none",
	}
	body, err := client.PostJSON(registrationEndpoint, payload, nil)
	if err != nil {
		return "", fmt.Errorf("dynamic client registration failed: %w", err)
	}
	var resp map[string]any
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("invalid registration response: %w", err)
	}
	clientID, ok := resp["client_id"].(string)
	if !ok {
		return "", fmt.Errorf("no client_id in registration response")
	}
	return clientID, nil
}

// BuildAuthorizationURL constructs the OAuth authorize URL with PKCE params.
func BuildAuthorizationURL(config *OAuthConfig, clientID, challenge, state string) string {
	scope := strings.Join(config.ScopesSupported, " ")
	if scope == "" {
		scope = "read write"
	}
	params := url.Values{
		"client_id":             {clientID},
		"redirect_uri":         {OAuthRedirectURI},
		"response_type":        {"code"},
		"scope":                {scope},
		"state":                {state},
		"code_challenge":       {challenge},
		"code_challenge_method": {"S256"},
	}
	return config.AuthorizationEndpoint + "?" + params.Encode()
}

// StartCallbackServer starts a temporary HTTP server on 127.0.0.1:9876 to receive
// the OAuth redirect. Returns a channel for the result and a shutdown function.
func StartCallbackServer() (chan *CallbackResult, func(), error) {
	resultCh := make(chan *CallbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		result := &CallbackResult{
			Code:             q.Get("code"),
			State:            q.Get("state"),
			Error:            q.Get("error"),
			ErrorDescription: q.Get("error_description"),
		}

		w.Header().Set("Content-Type", "text/html")
		if result.Code != "" {
			w.WriteHeader(200)
			fmt.Fprint(w, `<html><body style="font-family:system-ui;text-align:center;padding:60px">
<h1>Logged in!</h1>
<p>You can close this tab and return to the terminal.</p>
</body></html>`)
		} else {
			errMsg := result.ErrorDescription
			if errMsg == "" {
				errMsg = "Unknown error"
			}
			w.WriteHeader(400)
			fmt.Fprintf(w, `<html><body style="font-family:system-ui;text-align:center;padding:60px">
<h1>Login failed</h1>
<p>%s</p>
</body></html>`, html.EscapeString(errMsg))
		}

		resultCh <- result
	})

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", OAuthRedirectPort))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start callback server on port %d: %w", OAuthRedirectPort, err)
	}

	server := &http.Server{Handler: mux}
	go server.Serve(listener)

	shutdown := func() { server.Close() }
	return resultCh, shutdown, nil
}

// ExchangeCode exchanges an authorization code for tokens.
func ExchangeCode(client *httpclient.Client, tokenEndpoint, code, verifier, clientID string) (map[string]any, error) {
	params := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {OAuthRedirectURI},
		"client_id":     {clientID},
		"code_verifier": {verifier},
	}
	body, err := client.PostForm(tokenEndpoint, params, nil)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}
	var resp map[string]any
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("invalid token response: %w", err)
	}
	return resp, nil
}

func stringOrDefault(m map[string]any, key, def string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return def
}
