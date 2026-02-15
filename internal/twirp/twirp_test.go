package twirp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lollipopai/cli/internal/auth"
	"github.com/lollipopai/cli/internal/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCall_URLConstruction(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/twirp/lollipop.proto.recipe.v1.RecipeV1/Search", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		json.NewEncoder(w).Encode(map[string]any{"recipes": []any{}})
	}))
	defer srv.Close()

	creds := &auth.Credentials{
		BaseURL:          srv.URL,
		OAuthAccessToken: "test-token",
	}
	caller := NewCaller(httpclient.New(), creds)
	_, err := caller.Call("lollipop.proto.recipe.v1.RecipeV1", "Search", map[string]any{"query": "curry"})
	require.NoError(t, err)
}

func TestCall_AuthHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer my-token", r.Header.Get("Authorization"))
		json.NewEncoder(w).Encode(map[string]string{})
	}))
	defer srv.Close()

	creds := &auth.Credentials{BaseURL: srv.URL, OAuthAccessToken: "my-token"}
	caller := NewCaller(httpclient.New(), creds)
	_, err := caller.Call("svc", "Method", nil)
	require.NoError(t, err)
}

func TestCall_401Hint(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer srv.Close()

	creds := &auth.Credentials{BaseURL: srv.URL, OAuthAccessToken: "expired-token"}
	caller := NewCaller(httpclient.New(), creds)
	_, err := caller.Call("svc", "Method", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Try: chp login")
}

func TestCall_NilPayload(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		// nil payload should be sent as empty object
		assert.NotNil(t, body)
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}))
	defer srv.Close()

	creds := &auth.Credentials{BaseURL: srv.URL, OAuthAccessToken: "tok"}
	caller := NewCaller(httpclient.New(), creds)
	result, err := caller.Call("svc", "Method", nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCall_NoToken(t *testing.T) {
	creds := &auth.Credentials{BaseURL: "http://localhost:9999"}
	caller := NewCaller(httpclient.New(), creds)
	_, err := caller.Call("svc", "Method", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not logged in")
}

func TestCall_ParsesResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"name": "Chicken Tikka",
			"id":   42,
		})
	}))
	defer srv.Close()

	creds := &auth.Credentials{BaseURL: srv.URL, OAuthAccessToken: "tok"}
	caller := NewCaller(httpclient.New(), creds)
	result, err := caller.Call("svc", "Get", nil)
	require.NoError(t, err)

	m, ok := result.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Chicken Tikka", m["name"])
}
