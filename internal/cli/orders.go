package cli

import (
	"fmt"
	"strings"

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

		uids := extractProductUIDs(result)
		if len(uids) > 0 {
			fmt.Println()
			output.Info("Product UIDs (for re-ordering):")
			fmt.Printf("  %s\n", strings.Join(uids, " "))
			fmt.Println()
			output.Info("Re-add to basket:")
			fmt.Printf("  chp basket add-product %s\n", strings.Join(uids, " "))
		}
	},
}

// extractProductUIDs walks a JSON-decoded response tree and collects product
// identifier values from known field names.
func extractProductUIDs(v any) []string {
	var uids []string
	seen := map[string]bool{}

	var walk func(v any)
	walk = func(v any) {
		switch val := v.(type) {
		case map[string]any:
			for key, child := range val {
				switch key {
				case "sainsburys_uid", "product_uid", "product_id", "uid":
					if s, ok := child.(string); ok && s != "" && !seen[s] {
						seen[s] = true
						uids = append(uids, s)
					}
				default:
					walk(child)
				}
			}
		case []any:
			for _, item := range val {
				walk(item)
			}
		}
	}
	walk(v)
	return uids
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
