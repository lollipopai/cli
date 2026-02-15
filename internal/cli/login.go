package cli

import (
	"fmt"
	"time"

	"github.com/lollipopai/cli/internal/auth"
	"github.com/lollipopai/cli/internal/httpclient"
	"github.com/lollipopai/cli/internal/output"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Sign in to Cherrypick via OAuth",
	Run: func(cmd *cobra.Command, args []string) {
		runLogin()
	},
}

func runLogin() {
	creds := auth.LoadCredentials()
	baseURL := creds.GetBaseURL()
	client := httpclient.New()

	output.Info(fmt.Sprintf("Starting OAuth login for %s...", baseURL))

	// Step 1: Discover OAuth configuration
	config, err := auth.DiscoverOAuthConfig(client, baseURL)
	if err != nil {
		output.Fatal(err.Error())
	}

	// Step 2: Register client if needed
	clientID := creds.OAuthClientID
	if clientID == "" {
		output.Info("Registering CLI client...")
		clientID, err = auth.RegisterClient(client, config.RegistrationEndpoint)
		if err != nil {
			output.Fatal(err.Error())
		}
		creds.OAuthClientID = clientID
		auth.SaveCredentials(creds)
		output.Success(fmt.Sprintf("Client registered: %s", clientID))
	}

	// Step 3: Generate PKCE values
	verifier, err := auth.GenerateCodeVerifier()
	if err != nil {
		output.Fatal(fmt.Sprintf("Failed to generate PKCE verifier: %v", err))
	}
	challenge := auth.GenerateCodeChallenge(verifier)
	state, err := auth.GenerateState()
	if err != nil {
		output.Fatal(fmt.Sprintf("Failed to generate state: %v", err))
	}

	// Step 4: Build authorization URL
	authURL := auth.BuildAuthorizationURL(config, clientID, challenge, state)

	// Step 5: Start callback server
	resultCh, shutdown, err := auth.StartCallbackServer()
	if err != nil {
		output.Fatal(err.Error())
	}
	defer shutdown()

	// Step 6: Open browser
	output.Info("Opening browser for authorization...")
	fmt.Printf("\n  %s\n  %s\n\n", output.Dim("If the browser doesn't open, visit:"), authURL)
	browser.OpenURL(authURL)

	// Step 7: Wait for callback (120s timeout)
	output.Info("Waiting for authorization callback...")
	select {
	case result := <-resultCh:
		if result.Code == "" {
			errMsg := result.ErrorDescription
			if errMsg == "" {
				errMsg = result.Error
			}
			if errMsg == "" {
				errMsg = "Unknown error"
			}
			output.Fatal(fmt.Sprintf("Authorization failed: %s", errMsg))
		}

		// Verify state
		if result.State != state {
			output.Fatal("OAuth state mismatch - possible CSRF attack.")
		}

		// Step 8: Exchange code for tokens
		output.Info("Exchanging authorization code for tokens...")
		tokenResp, err := auth.ExchangeCode(client, config.TokenEndpoint, result.Code, verifier, clientID)
		if err != nil {
			output.Fatal(err.Error())
		}

		// Step 9: Save tokens
		if at, ok := tokenResp["access_token"].(string); ok {
			creds.OAuthAccessToken = at
		}
		if rt, ok := tokenResp["refresh_token"].(string); ok {
			creds.OAuthRefreshToken = rt
		}
		if ei, ok := tokenResp["expires_in"].(float64); ok {
			creds.OAuthExpiresAt = time.Now().Unix() + int64(ei)
		}
		creds.BaseURL = baseURL
		auth.SaveCredentials(creds)

		output.Success("OAuth login successful!")
		output.Info(fmt.Sprintf("Credentials saved to %s", auth.CredentialsFile))

	case <-time.After(120 * time.Second):
		output.Fatal("Timed out waiting for OAuth callback.")
	}
}
