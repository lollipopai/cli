package cli

import (
	"github.com/lollipopai/cli/internal/output"
	"github.com/spf13/cobra"
)

var recipesCmd = &cobra.Command{
	Use:   "recipes",
	Short: "Recipe commands",
}

var recipesSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search recipes",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call("lollipop.proto.recipe.v1.RecipeV1", "Search", map[string]any{
			"query": args[0],
		})
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

var recipesGetCmd = &cobra.Command{
	Use:   "get <slug-or-id>",
	Short: "Get a recipe by slug or ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		// If it contains letters, treat as slug; otherwise use id field
		isSlug := false
		for _, c := range identifier {
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
				isSlug = true
				break
			}
		}

		var payload map[string]any
		if isSlug {
			payload = map[string]any{"slug": identifier}
		} else {
			payload = map[string]any{"id": identifier}
		}

		caller := newTwirpCaller()
		result, err := caller.Call("lollipop.proto.recipe.v1.RecipeV1", "GetBySlug", payload)
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

func init() {
	recipesCmd.AddCommand(recipesSearchCmd)
	recipesCmd.AddCommand(recipesGetCmd)
}
