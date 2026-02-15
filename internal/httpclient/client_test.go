package httpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetJSON_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.Header.Get("User-Agent"), "chp-cli/")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"hello": "world"})
	}))
	defer srv.Close()

	c := New()
	body, err := c.GetJSON(srv.URL+"/test", nil)
	require.NoError(t, err)

	var result map[string]string
	require.NoError(t, json.Unmarshal(body, &result))
	assert.Equal(t, "world", result["hello"])
}

func TestPostJSON_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload map[string]string
		json.NewDecoder(r.Body).Decode(&payload)
		assert.Equal(t, "bar", payload["foo"])

		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer srv.Close()

	c := New()
	body, err := c.PostJSON(srv.URL, map[string]string{"foo": "bar"}, nil)
	require.NoError(t, err)

	var result map[string]string
	require.NoError(t, json.Unmarshal(body, &result))
	assert.Equal(t, "ok", result["status"])
}

func TestAPIError_OnNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"error":"not found"}`))
	}))
	defer srv.Close()

	c := New()
	_, err := c.GetJSON(srv.URL, nil)
	require.Error(t, err)

	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 404, apiErr.StatusCode)
	assert.Contains(t, apiErr.Message, "not found")
}

func TestAPIError_MsgField(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte(`{"msg":"server exploded"}`))
	}))
	defer srv.Close()

	c := New()
	_, err := c.GetJSON(srv.URL, nil)
	apiErr := err.(*APIError)
	assert.Contains(t, apiErr.Message, "server exploded")
}

func TestAPIError_PlainText(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
		w.Write([]byte("Service Unavailable"))
	}))
	defer srv.Close()

	c := New()
	_, err := c.GetJSON(srv.URL, nil)
	apiErr := err.(*APIError)
	assert.Contains(t, apiErr.Message, "Service Unavailable")
}

func TestCustomHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer tok123", r.Header.Get("Authorization"))
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	c := New()
	_, err := c.PostJSON(srv.URL, nil, map[string]string{
		"Authorization": "Bearer tok123",
	})
	require.NoError(t, err)
}

func TestUserAgent(t *testing.T) {
	SetUserAgent("chp-cli/test")
	defer SetUserAgent("chp-cli/dev")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "chp-cli/test", r.Header.Get("User-Agent"))
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	c := New()
	_, err := c.GetJSON(srv.URL, nil)
	require.NoError(t, err)
}
