package cli

import (
	"github.com/lollipopai/cli/internal/output"
	"github.com/spf13/cobra"
)

var productsCmd = &cobra.Command{
	Use:   "products",
	Short: "Product commands",
}

var productsSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search products",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call("lollipop.proto.product.v2.ProductV2", "Search", map[string]any{
			"keyword": args[0],
		})
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

var productsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a product by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call("lollipop.proto.product.v2.ProductV2", "Get", map[string]any{
			"id": args[0],
		})
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

func init() {
	productsCmd.AddCommand(productsSearchCmd)
	productsCmd.AddCommand(productsGetCmd)
}
