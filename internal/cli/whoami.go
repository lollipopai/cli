package cli

import (
	"github.com/lollipopai/cli/internal/output"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user profile",
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call("lollipop.proto.user.v1.UserV1", "Current", nil)
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}
