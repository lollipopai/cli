package cli

import (
	"github.com/lollipopai/cli/internal/output"
	"github.com/spf13/cobra"
)

var ordersCmd = &cobra.Command{
	Use:   "orders",
	Short: "Order commands",
	Run: func(cmd *cobra.Command, args []string) {
		runOrdersList()
	},
}

var ordersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List order summaries",
	Run: func(cmd *cobra.Command, args []string) {
		runOrdersList()
	},
}

var ordersGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an order by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call("lollipop.proto.order.v1.OrderV1", "Get", map[string]any{
			"id": args[0],
		})
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

func runOrdersList() {
	caller := newTwirpCaller()
	result, err := caller.Call("lollipop.proto.order.v1.OrderV1", "SummaryList", nil)
	if err != nil {
		output.Fatal(err.Error())
	}
	output.PrintJSON(result)
}

func init() {
	ordersCmd.AddCommand(ordersListCmd)
	ordersCmd.AddCommand(ordersGetCmd)
}
