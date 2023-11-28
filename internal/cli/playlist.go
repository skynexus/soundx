package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/skynexus/soundx/api"
	"github.com/spf13/cobra"
)

func setupPlaylistCmd() *cobra.Command {
	var playlistCmd = &cobra.Command{
		Use:   "playlist [subcommand]",
		Short: "Manage playlists",
		Long:  `A playlist is a collection of sounds that can be listened to`,
		Args:  cobra.MinimumNArgs(1),
	}

	var playlistListCmd = &cobra.Command{
		Use:   "list",
		Short: "List playlists",
		Long:  "List all availble playlists sorted by creation date.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			magenta := color.New(color.FgMagenta).SprintFunc()

			resp, err := soundx.GetPlaylists(ctx)
			if err != nil {
				log.Fatalf("error: soundx request failed: %s", err)
			} else if resp.StatusCode == 200 {
				var playlistResponse api.PlaylistResponse
				if err := json.NewDecoder(resp.Body).Decode(&playlistResponse); err != nil {
					log.Fatalf("error: could not parse soundx response: %s", err)
				} else {
					fmt.Printf("%4s %s\n", "Id", "Title")
					for _, playlist := range playlistResponse.Data {
						fmt.Printf("%4s %s\n", magenta(fmt.Sprintf("%4v", playlist.Id)), playlist.Title)
					}
				}
			} else {
				log.Printf("soundx returned error (%d)", resp.StatusCode)
				buf := new(strings.Builder)
				if _, err := io.Copy(buf, resp.Body); err != nil {
					log.Fatalf("unexpected response code %d: not able to read response body: %s", resp.StatusCode, err)
				}
				log.Fatalf("unexpected response code %d: response body: %s", resp.StatusCode, buf.String())
			}
		},
	}

	playlistCmd.AddCommand(playlistListCmd)
	return playlistCmd
}
