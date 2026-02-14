package cli

import (
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
	Use:   "add-recipe <recipe-id>",
	Short: "Add a recipe to the basket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call(basketService, "AddRecipe", map[string]any{
			"recipe_id": args[0],
		})
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

var basketRemoveRecipeCmd = &cobra.Command{
	Use:   "remove-recipe <recipe-id>",
	Short: "Remove a recipe from the basket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call(basketService, "RemoveRecipe", map[string]any{
			"recipe_id": args[0],
		})
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

var basketAddProductCmd = &cobra.Command{
	Use:   "add-product <product-id>",
	Short: "Add a product to the basket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		payload := map[string]any{
			"product_id": args[0],
		}
		if basketQuantity > 0 {
			payload["quantity"] = basketQuantity
		}
		caller := newTwirpCaller()
		result, err := caller.Call(basketService, "AddProduct", payload)
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

var basketRemoveProductCmd = &cobra.Command{
	Use:   "remove-product <product-id>",
	Short: "Remove a product from the basket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call(basketService, "RemoveProduct", map[string]any{
			"product_id": args[0],
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

func runBasketShow() {
	caller := newTwirpCaller()
	result, err := caller.Call(basketService, "Show", nil)
	if err != nil {
		output.Fatal(err.Error())
	}
	output.PrintJSON(result)
}

func init() {
	basketAddProductCmd.Flags().IntVarP(&basketQuantity, "quantity", "q", 0, "Quantity to add (default: server decides)")
	basketCmd.AddCommand(basketShowCmd)
	basketCmd.AddCommand(basketAddRecipeCmd)
	basketCmd.AddCommand(basketRemoveRecipeCmd)
	basketCmd.AddCommand(basketAddProductCmd)
	basketCmd.AddCommand(basketRemoveProductCmd)
	basketCmd.AddCommand(basketClearCmd)
}
