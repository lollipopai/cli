package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCodeVerifier_Length(t *testing.T) {
	v, err := GenerateCodeVerifier()
	require.NoError(t, err)
	// 64 random bytes → 86 base64url chars (matches Python secrets.token_urlsafe(64))
	assert.Len(t, v, 86)
	assert.GreaterOrEqual(t, len(v), 43, "PKCE minimum")
	assert.LessOrEqual(t, len(v), 128, "PKCE maximum")
}

func TestGenerateCodeVerifier_Unique(t *testing.T) {
	v1, _ := GenerateCodeVerifier()
	v2, _ := GenerateCodeVerifier()
	assert.NotEqual(t, v1, v2)
}

func TestGenerateCodeVerifier_Base64URL(t *testing.T) {
	v, _ := GenerateCodeVerifier()
	// Should only contain base64url chars (no padding)
	for _, c := range v {
		assert.True(t,
			(c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
				(c >= '0' && c <= '9') || c == '-' || c == '_',
			"unexpected character: %c", c)
	}
}

func TestGenerateCodeChallenge_S256(t *testing.T) {
	// RFC 7636 Appendix B test vector
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	expected := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"
	assert.Equal(t, expected, GenerateCodeChallenge(verifier))
}

func TestGenerateCodeChallenge_Manual(t *testing.T) {
	verifier := "test-verifier-value"
	h := sha256.Sum256([]byte(verifier))
	expected := base64.RawURLEncoding.EncodeToString(h[:])
	assert.Equal(t, expected, GenerateCodeChallenge(verifier))
}

func TestGenerateState_Length(t *testing.T) {
	s, err := GenerateState()
	require.NoError(t, err)
	// 32 bytes base64url = 43 chars
	assert.Len(t, s, 43)
}

func TestGenerateState_Unique(t *testing.T) {
	s1, _ := GenerateState()
	s2, _ := GenerateState()
	assert.NotEqual(t, s1, s2)
}
