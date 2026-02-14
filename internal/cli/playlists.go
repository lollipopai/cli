package cli

import (
	"github.com/lollipopai/cli/internal/output"
	"github.com/spf13/cobra"
)

var playlistsCmd = &cobra.Command{
	Use:   "playlists",
	Short: "Playlist commands",
	Run: func(cmd *cobra.Command, args []string) {
		runPlaylistsList()
	},
}

var playlistsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List playlists",
	Run: func(cmd *cobra.Command, args []string) {
		runPlaylistsList()
	},
}

var playlistsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a playlist by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		caller := newTwirpCaller()
		result, err := caller.Call("lollipop.proto.playlist.v1.PlaylistV1", "Get", map[string]any{
			"id": args[0],
		})
		if err != nil {
			output.Fatal(err.Error())
		}
		output.PrintJSON(result)
	},
}

func runPlaylistsList() {
	caller := newTwirpCaller()
	result, err := caller.Call("lollipop.proto.playlist.v1.PlaylistV1", "List", nil)
	if err != nil {
		output.Fatal(err.Error())
	}
	output.PrintJSON(result)
}

func init() {
	playlistsCmd.AddCommand(playlistsListCmd)
	playlistsCmd.AddCommand(playlistsGetCmd)
}
