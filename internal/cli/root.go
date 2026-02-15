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
	Use:   "cpk",
	Short: "Cherrypick CLI - interact with the Cherrypick API",
	Long: `Cherrypick CLI - interact with the Cherrypick API

Examples:
  cpk login                              Sign in via OAuth in browser
  cpk whoami                             Show current user
  cpk recipes search curry               Search for recipes
  cpk recipes get chicken-tikka          Get recipe by slug
  cpk products search milk               Search products
  cpk products get 7834128               Get product by Sainsbury's UID
  cpk basket                             Show basket
  cpk basket add-recipe 1 2 3            Add recipes to basket
  cpk basket add-product 7834128 7209381 Add products to basket
  cpk basket add-product 7834128:2       Add product with quantity
  cpk basket set-quantity 7834128 4      Change product quantity
  cpk orders                             List orders
  cpk orders get 42                      Get order with product UIDs
  cpk slots                              List delivery slots
  cpk slots book 5                       Book a delivery slot
  cpk plan                               Show current meal plan
  cpk plan add-recipe 1 100 101          Add recipes to a plan
  cpk playlists                          List playlists
  cpk config set-url https://alpha.lollipopai.com
  cpk config show                        Show current config
  cpk call <service> <method> [payload]  Raw Twirp call
  cpk logout                             Clear credentials`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute(version string) {
	Version = version
	rootCmd.Version = version
	httpclient.SetUserAgent("cpk-cli/" + version)

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
