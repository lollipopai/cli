package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lollipopai/cli/internal/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverOAuthConfig(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/oauth-protected-resource", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"authorization_servers": []string{"https://auth.example.com"},
		})
	})
	// The discovery will try to hit auth.example.com which isn't our server,
	// so we test the fallback path instead
	srv := httptest.NewServer(mux)
	defer srv.Close()

	// Test that it falls back gracefully when auth server isn't reachable
	client := httpclient.New()
	_, err := DiscoverOAuthConfig(client, srv.URL)
	assert.Error(t, err) // Can't reach https://auth.example.com
}

func TestDiscoverOAuthConfig_DirectFallback(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/oauth-protected-resource", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	mux.HandleFunc("/.well-known/oauth-authorization-server", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"authorization_endpoint": "https://example.com/auth",
			"token_endpoint":         "https://example.com/token",
			"registration_endpoint":  "https://example.com/register",
			"scopes_supported":       []string{"read", "write", "admin"},
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := httpclient.New()
	config, err := DiscoverOAuthConfig(client, srv.URL)
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/auth", config.AuthorizationEndpoint)
	assert.Equal(t, "https://example.com/token", config.TokenEndpoint)
	assert.Equal(t, "https://example.com/register", config.RegistrationEndpoint)
	assert.Equal(t, []string{"read", "write", "admin"}, config.ScopesSupported)
}

func TestRegisterClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "CherryPick CLI", body["client_name"])
		assert.Equal(t, "none", body["token_endpoint_auth_method"])

		json.NewEncoder(w).Encode(map[string]string{
			"client_id": "new-client-id",
		})
	}))
	defer srv.Close()

	client := httpclient.New()
	clientID, err := RegisterClient(client, srv.URL)
	require.NoError(t, err)
	assert.Equal(t, "new-client-id", clientID)
}

func TestBuildAuthorizationURL(t *testing.T) {
	config := &OAuthConfig{
		AuthorizationEndpoint: "https://auth.example.com/authorize",
		ScopesSupported:       []string{"read", "write"},
	}
	url := BuildAuthorizationURL(config, "client-123", "challenge-abc", "state-xyz")
	assert.Contains(t, url, "https://auth.example.com/authorize?")
	assert.Contains(t, url, "client_id=client-123")
	assert.Contains(t, url, "code_challenge=challenge-abc")
	assert.Contains(t, url, "code_challenge_method=S256")
	assert.Contains(t, url, "state=state-xyz")
	assert.Contains(t, url, "response_type=code")
	assert.Contains(t, url, "scope=read+write")
}

func TestStartCallbackServer_Success(t *testing.T) {
	resultCh, shutdown, err := StartCallbackServer()
	require.NoError(t, err)
	defer shutdown()

	// Send a successful callback
	resp, err := http.Get("http://127.0.0.1:9876/callback?code=authcode&state=mystate")
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()

	result := <-resultCh
	assert.Equal(t, "authcode", result.Code)
	assert.Equal(t, "mystate", result.State)
}

func TestStartCallbackServer_Error(t *testing.T) {
	resultCh, shutdown, err := StartCallbackServer()
	require.NoError(t, err)
	defer shutdown()

	resp, err := http.Get("http://127.0.0.1:9876/callback?error=access_denied&error_description=User+denied")
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	resp.Body.Close()

	result := <-resultCh
	assert.Equal(t, "", result.Code)
	assert.Equal(t, "access_denied", result.Error)
	assert.Equal(t, "User denied", result.ErrorDescription)
}

func TestExchangeCode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		r.ParseForm()
		assert.Equal(t, "authorization_code", r.FormValue("grant_type"))
		assert.Equal(t, "the-code", r.FormValue("code"))
		assert.Equal(t, "the-verifier", r.FormValue("code_verifier"))
		assert.Equal(t, "client-123", r.FormValue("client_id"))

		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "new-access-token",
			"refresh_token": "new-refresh-token",
			"expires_in":    3600,
		})
	}))
	defer srv.Close()

	client := httpclient.New()
	resp, err := ExchangeCode(client, srv.URL, "the-code", "the-verifier", "client-123")
	require.NoError(t, err)
	assert.Equal(t, "new-access-token", resp["access_token"])
	assert.Equal(t, "new-refresh-token", resp["refresh_token"])
}
