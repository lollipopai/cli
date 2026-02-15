package cli

import (
	"github.com/lollipopai/cli/internal/output"
	"github.com/spf13/cobra"
)

const planService = "lollipop.proto.plan.v1.PlanV1"

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Meal plan commands",
	Run: func(cmd *cobra.Command, args []string) {
		runPlanShow()
	},
}

var planShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current meal plan",
	Run: func(cmd *cobra.Command, args []string) {
		runPlanShow()
	},
}

var planListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available plans",
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call(planService, "List", nil)
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

var planGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a specific plan",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call(planService, "Get", map[string]any{
			"id": args[0],
		})
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

var planAddRecipeCmd = &cobra.Command{
	Use:   "add-recipe <plan-id> <recipe-id>...",
	Short: "Add one or more recipes to a plan",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		planID := args[0]
		recipeIDs := args[1:]
		caller := newTwirpCaller()
		var result any
		for _, recipeID := range recipeIDs {
			var err error
			result, err = caller.Call(planService, "AddRecipe", map[string]any{
				"plan_id":   planID,
				"recipe_id": recipeID,
			})
			if err != nil {
				output.Fatal(err.Error())
			}
		}
		output.PrintJSON(result)
	},
}

var planRemoveRecipeCmd = &cobra.Command{
	Use:   "remove-recipe <plan-id> <recipe-id>",
	Short: "Remove a recipe from a plan",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call(planService, "RemoveRecipe", map[string]any{
			"plan_id":   args[0],
			"recipe_id": args[1],
		})
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

func runPlanShow() {
	caller := newTwirpCaller()
	result, err := caller.Call(planService, "Show", nil)
	if err != nil {
		output.Fatal(err.Error())
	}
	output.PrintJSON(result)
}

func init() {
	planCmd.AddCommand(planShowCmd)
	planCmd.AddCommand(planListCmd)
	planCmd.AddCommand(planGetCmd)
	planCmd.AddCommand(planAddRecipeCmd)
	planCmd.AddCommand(planRemoveRecipeCmd)
}
