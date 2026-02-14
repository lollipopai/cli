package cli

import (
	"os"

	"github.com/lollipopai/cli/internal/auth"
	"github.com/lollipopai/cli/internal/output"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear saved credentials",
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.Remove(auth.CredentialsFile); err != nil {
			if os.IsNotExist(err) {
				output.Info("No credentials found. Already logged out.")
				return
			}
			output.Fatal(err.Error())
		}
		output.Success("Logged out. Credentials removed.")
	},
}
