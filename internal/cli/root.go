package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/skynexus/soundx/api"
	"github.com/spf13/cobra"
)

var soundx *api.Client

var rootCmd = &cobra.Command{
	Use:   "soundx-cli",
	Short: "The SoundX Command Line Interface is a tool to work with SoundX",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		apiUrl := "http://localhost:8080"
		if s, ok := os.LookupEnv("SOUNDX_URL"); ok {
			apiUrl = s
		}

		if client, clientErr := api.NewClient(apiUrl); clientErr != nil {
			log.Fatalf("error: %s", clientErr)
		} else {
			soundx = client
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(setupSoundCmd())
	rootCmd.AddCommand(setupPlaylistCmd())
}
