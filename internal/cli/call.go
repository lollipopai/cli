package cli

import (
	"encoding/json"
	"fmt"

	"github.com/lollipopai/cli/internal/output"
	"github.com/spf13/cobra"
)

var callCmd = &cobra.Command{
	Use:   "call <service> <method> [payload]",
	Short: "Make a raw Twirp RPC call",
	Long: `Make a raw Twirp RPC call to any service endpoint.

Examples:
  cpk call lollipop.proto.recipe.v1.RecipeV1 Search '{"query":"curry"}'
  cpk call lollipop.proto.user.v1.UserV1 Current
  cpk call lollipop.proto.basket.v1.BasketV1 Show`,
	Args: cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
		method := args[1]

		var payload any
		if len(args) == 3 {
			if err := json.Unmarshal([]byte(args[2]), &payload); err != nil {
				output.Fatal(fmt.Sprintf("Invalid JSON payload: %v", err))
			}
		}

		caller := newTwirpCaller()
		result, err := caller.Call(service, method, payload)
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}
