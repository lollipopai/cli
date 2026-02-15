package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lollipopai/cli/internal/output"
	"github.com/spf13/cobra"
)

const basketService = "lollipop.proto.basket.v1.BasketV1"

var basketQuantity int

var basketCmd = &cobra.Command{
	Use:   "basket",
	Short: "Basket commands",
	Run: func(cmd *cobra.Command, args []string) {
		runBasketShow()
	},
}

var basketShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current basket",
	Run: func(cmd *cobra.Command, args []string) {
		runBasketShow()
	},
}

var basketAddRecipeCmd = &cobra.Command{
	Use:   "add-recipe <recipe-id>...",
	Short: "Add one or more recipes to the basket",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		var result any
		for _, id := range args {
			var err error
			result, err = caller.Call(basketService, "AddRecipe", map[string]any{
				"recipe_id": id,
			})
			if err != nil {
				output.Fatal(err.Error())
			}
		}
		output.PrintJSON(result)
	},
}

var basketRemoveRecipeCmd = &cobra.Command{
	Use:   "remove-recipe <recipe-id>...",
	Short: "Remove one or more recipes from the basket",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		var result any
		for _, id := range args {
			var err error
			result, err = caller.Call(basketService, "RemoveRecipe", map[string]any{
				"recipe_id": id,
			})
			if err != nil {
				output.Fatal(err.Error())
			}
		}
		output.PrintJSON(result)
	},
}

var basketAddProductCmd = &cobra.Command{
	Use:   "add-product <uid>...",
	Short: "Add one or more products to the basket by Sainsbury's product UID",
	Long: `Add one or more products to the basket by Sainsbury's product UID.

Per-item quantities can be specified with uid:qty syntax.
The -q flag sets the default quantity for items without an explicit quantity.

Examples:
  cpk basket add-product 7834128 7209381           # qty defaults to 1
  cpk basket add-product 7834128:2 7209381:3       # per-item quantities
  cpk basket add-product 7834128 7209381 -q 2      # -q sets default for all
  cpk basket add-product 7834128:3 7209381 -q 2    # 7834128→3, 7209381→2`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defaultQty := basketQuantity
		if defaultQty == 0 {
			defaultQty = 1
		}
		caller := newTwirpCaller()
		var result any
		for _, arg := range args {
			uid, qty := parseProductArg(arg, defaultQty)
			payload := map[string]any{
				"product_id": uid,
				"quantity":   qty,
			}
			var err error
			result, err = caller.Call(basketService, "AddProduct", payload)
			if err != nil {
				output.Fatal(err.Error())
			}
		}
		output.PrintJSON(result)
	},
}

var basketRemoveProductCmd = &cobra.Command{
	Use:   "remove-product <uid>...",
	Short: "Remove one or more products from the basket by Sainsbury's product UID",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		var result any
		for _, uid := range args {
			var err error
			result, err = caller.Call(basketService, "RemoveProduct", map[string]any{
				"product_id": uid,
			})
			if err != nil {
				output.Fatal(err.Error())
			}
		}
		output.PrintJSON(result)
	},
}

var basketSetQuantityCmd = &cobra.Command{
	Use:   "set-quantity <uid> <qty>",
	Short: "Set the quantity of a product in the basket by Sainsbury's product UID",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		qty, err := strconv.Atoi(args[1])
		if err != nil {
			output.Fatal(fmt.Sprintf("invalid quantity %q: must be a number", args[1]))
		}
		caller := newTwirpCaller()
		result, err := caller.Call(basketService, "SetQuantity", map[string]any{
			"product_id": args[0],
			"quantity":   qty,
		})
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

var basketClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the basket",
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call(basketService, "Clear", nil)
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

// parseProductArg splits "uid:qty" into uid and quantity.
// If no colon is present, defaultQty is used.
func parseProductArg(arg string, defaultQty int) (string, int) {
	parts := strings.SplitN(arg, ":", 2)
	if len(parts) == 2 {
		qty, err := strconv.Atoi(parts[1])
		if err != nil {
			output.Fatal(fmt.Sprintf("invalid quantity in %q: %s", arg, err))
		}
		return parts[0], qty
	}
	return arg, defaultQty
}

func runBasketShow() {
	caller := newTwirpCaller()
	result, err := caller.Call(basketService, "Show", nil)
	if err != nil {
		output.Fatal(err.Error())
	}
	output.PrintJSON(result)
}

func init() {
	basketAddProductCmd.Flags().IntVarP(&basketQuantity, "quantity", "q", 0, "Default quantity for items without :qty suffix (default: 1)")
	basketCmd.AddCommand(basketShowCmd)
	basketCmd.AddCommand(basketAddRecipeCmd)
	basketCmd.AddCommand(basketRemoveRecipeCmd)
	basketCmd.AddCommand(basketAddProductCmd)
	basketCmd.AddCommand(basketRemoveProductCmd)
	basketCmd.AddCommand(basketSetQuantityCmd)
	basketCmd.AddCommand(basketClearCmd)
}
