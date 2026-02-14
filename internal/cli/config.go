package cli

import (
	"fmt"
	"net/url"
	"time"

	"github.com/lollipopai/cli/internal/auth"
	"github.com/lollipopai/cli/internal/output"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration commands",
}

var configSetURLCmd = &cobra.Command{
	Use:   "set-url <url>",
	Short: "Set the base API URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rawURL := args[0]
		// Strip trailing slash
		for len(rawURL) > 0 && rawURL[len(rawURL)-1] == '/' {
			rawURL = rawURL[:len(rawURL)-1]
		}

		parsed, err := url.Parse(rawURL)
		if err != nil {
			output.Fatal(fmt.Sprintf("Invalid URL: %v", err))
		}
		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			output.Fatal("Base URL must use http or https scheme.")
		}
		if parsed.Hostname() == "" {
			output.Fatal("Base URL must include a hostname.")
		}

		creds := auth.LoadCredentials()
		creds.BaseURL = rawURL
		if err := auth.SaveCredentials(creds); err != nil {
			output.Fatal(err.Error())
		}
		output.Success(fmt.Sprintf("Base URL set to %s", output.Bold(rawURL)))
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		creds := auth.LoadCredentials()

		config := map[string]any{
			"base_url":         creds.GetBaseURL(),
			"credentials_file": auth.CredentialsFile,
			"has_jwt":          creds.JWT != "",
			"has_oauth_token":  creds.OAuthAccessToken != "",
			"oauth_client_id":  nilIfEmpty(creds.OAuthClientID),
		}

		if creds.OAuthExpiresAt > 0 {
			remaining := creds.OAuthExpiresAt - time.Now().Unix()
			if remaining > 0 {
				config["oauth_token_expires_in"] = fmt.Sprintf("%ds", remaining)
			} else {
				config["oauth_token_expired"] = true
			}
		}

		output.PrintJSON(config)
	},
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func init() {
	configCmd.AddCommand(configSetURLCmd)
	configCmd.AddCommand(configShowCmd)
}
