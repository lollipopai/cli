package cli

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/lollipopai/cli/internal/auth"
	"github.com/lollipopai/cli/internal/httpclient"
	"github.com/lollipopai/cli/internal/twirp"
	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "chp",
	Short: "Cherrypick CLI - interact with the Cherrypick API",
	Long: `Cherrypick CLI - interact with the Cherrypick API

Examples:
  chp login                              Sign in via OAuth in browser
  chp whoami                             Show current user
  chp recipes search curry               Search for recipes
  chp recipes get chicken-tikka          Get recipe by slug
  chp products search milk               Search products
  chp products get 7834128               Get product by Sainsbury's UID
  chp basket                             Show basket
  chp basket add-recipe 1 2 3            Add recipes to basket
  chp basket add-product 7834128 7209381 Add products to basket
  chp basket add-product 7834128:2       Add product with quantity
  chp basket set-quantity 7834128 4      Change product quantity
  chp orders                             List orders
  chp orders get 42                      Get order with product UIDs
  chp slots                              List delivery slots
  chp slots book 5                       Book a delivery slot
  chp plan                               Show current meal plan
  chp plan add-recipe 1 100 101          Add recipes to a plan
  chp playlists                          List playlists
  chp config show                        Show current config
  chp call recipe.v1.RecipeV1 Search     Raw Twirp call
  chp logout                             Clear credentials`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute(version string) {
	Version = version
	rootCmd.Version = version
	httpclient.SetUserAgent("chp-cli/" + version)

	// SIGINT → exit 130
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println()
		os.Exit(130)
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// newTwirpCaller creates a Caller with loaded credentials. Used by command handlers.
func newTwirpCaller() *twirp.Caller {
	creds := auth.LoadCredentials()
	client := httpclient.New()
	return twirp.NewCaller(client, creds)
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(recipesCmd)
	rootCmd.AddCommand(productsCmd)
	rootCmd.AddCommand(ordersCmd)
	rootCmd.AddCommand(playlistsCmd)
	rootCmd.AddCommand(basketCmd)
	rootCmd.AddCommand(slotsCmd)
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(callCmd)
}
