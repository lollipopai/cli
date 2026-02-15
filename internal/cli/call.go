package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lollipopai/cli/internal/output"
	"github.com/spf13/cobra"
)

const twirpServicePrefix = "lollipop.proto."

var callCmd = &cobra.Command{
	Use:   "call <service> <method> [payload]",
	Short: "Make a raw Twirp RPC call",
	Long: `Make a raw Twirp RPC call to any service endpoint.

The lollipop.proto. prefix is added automatically.

Examples:
  cpk call recipe.v1.RecipeV1 Search '{"query":"curry"}'
  cpk call user.v1.UserV1 Current
  cpk call basket.v1.BasketV1 Show
  cpk call product.v2.ProductV2 Search '{"keyword":"eggs"}'
  cpk call slot.v1.SlotV1 List
  cpk call plan.v1.PlanV1 Show`,
	Args: cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
		if !strings.HasPrefix(service, twirpServicePrefix) {
			service = twirpServicePrefix + service
		}
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
