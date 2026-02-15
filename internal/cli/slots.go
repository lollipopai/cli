package cli

import (
	"github.com/lollipopai/cli/internal/output"
	"github.com/spf13/cobra"
)

const slotService = "lollipop.proto.slot.v1.SlotV1"

var slotsCmd = &cobra.Command{
	Use:     "slots",
	Aliases: []string{"delivery"},
	Short:   "Delivery slot commands",
	Run: func(cmd *cobra.Command, args []string) {
		runSlotsList()
	},
}

var slotsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available delivery slots",
	Run: func(cmd *cobra.Command, args []string) {
		runSlotsList()
	},
}

var slotsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get delivery slot details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call(slotService, "Get", map[string]any{
			"id": args[0],
		})
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

var slotsBookCmd = &cobra.Command{
	Use:   "book <id>",
	Short: "Book a delivery slot",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call(slotService, "Book", map[string]any{
			"id": args[0],
		})
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

func runSlotsList() {
	caller := newTwirpCaller()
	result, err := caller.Call(slotService, "List", nil)
	if err != nil {
		output.Fatal(err.Error())
	}
	output.PrintJSON(result)
}

func init() {
	slotsCmd.AddCommand(slotsListCmd)
	slotsCmd.AddCommand(slotsGetCmd)
	slotsCmd.AddCommand(slotsBookCmd)
}
