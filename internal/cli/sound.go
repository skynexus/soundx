package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/skynexus/soundx/api"
	"github.com/spf13/cobra"
)

func printSound(sound api.Sound) {
	fmt.Printf("%16s: %s\n", "Id", color.MagentaString(sound.Id))
	fmt.Printf("%16s: %s\n", "Title", sound.Title)
	fmt.Printf("%16s: %d\n", "BPM", sound.Bpm)
	fmt.Printf("%16s: %d\n", "Duration", sound.DurationInSeconds)
	fmt.Printf("%16s: %s\n", "Genres", sound.Genres)
	fmt.Printf("%16s: %s\n", "Credit", sound.Credits)
	fmt.Printf("%16s: %v\n", "CreatedAt", sound.CreatedAt.Format("2006-01-02 15:04 -07:00"))
	fmt.Printf("%16s: %v\n", "UpdatedAt", sound.UpdatedAt.Format("2006-01-02 15:04 -07:00"))
}

func setupSoundCmd() *cobra.Command {
	var soundId int64
	var playlistId int64
	var title string
	var bpm int32
	var duration int32
	var credits []string
	var genres []string

	var soundCmd = &cobra.Command{
		Use:   "sound [subcommand]",
		Short: "Manage sounds",
		Long:  `A sound is a recorded material that can be listened to`,
		Args:  cobra.MinimumNArgs(1),
	}

	var soundListCmd = &cobra.Command{
		Use:   "list",
		Short: "List sounds",
		Long:  "List all availble sounds sorted by creation date.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			magenta := color.New(color.FgMagenta).SprintFunc()

			resp, err := soundx.GetSounds(ctx)
			if err != nil {
				log.Fatalf("error: soundx request failed: %s", err)
			} else if resp.StatusCode == 200 {
				var soundResponse api.SoundResponse
				if err := json.NewDecoder(resp.Body).Decode(&soundResponse); err != nil {
					log.Fatalf("error: could not parse soundx response: %s", err)
				} else {
					fmt.Printf("%4s %s\n", "Id", "Title")
					for _, sound := range soundResponse.Data {
						fmt.Printf("%4s %s\n", magenta(fmt.Sprintf("%4v", sound.Id)), sound.Title)
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

	var soundRecommendCmd = &cobra.Command{
		Use:   "recommend playlistid",
		Short: "Recommend sounds",
		Long:  "Recommend sounds based on playlist.",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}
			if n, err := strconv.ParseInt(args[0], 10, 64); err != nil {
				return fmt.Errorf("invalid playlist id: %w", err)
			} else {
				playlistId = n
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			magenta := color.New(color.FgMagenta).SprintFunc()

			sPlaylistId := strconv.FormatInt(playlistId, 10)
			input := api.GetRecommendedSoundsParams{
				PlaylistId: &sPlaylistId,
			}
			resp, err := soundx.GetRecommendedSounds(ctx, &input)
			if err != nil {
				log.Fatalf("error: soundx request failed: %s", err)
			} else if resp.StatusCode == 200 {
				var soundResponse api.SoundResponse
				if err := json.NewDecoder(resp.Body).Decode(&soundResponse); err != nil {
					log.Fatalf("error: could not parse soundx response: %s", err)
				} else {
					fmt.Printf("%4s %s\n", "Id", "Title")
					for _, sound := range soundResponse.Data {
						fmt.Printf("%4s %s\n", magenta(fmt.Sprintf("%4v", sound.Id)), sound.Title)
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

	var soundGetCmd = &cobra.Command{
		Use:   "get id",
		Short: "Get sound",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}
			if n, err := strconv.ParseInt(args[0], 10, 64); err != nil {
				return fmt.Errorf("invalid sound id: %w", err)
			} else {
				soundId = n
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			resp, err := soundx.GetSound(ctx, strconv.FormatInt(soundId, 10))
			if err != nil {
				log.Fatalf("error: soundx request failed: %s", err)
			} else if resp.StatusCode == 200 {
				var soundResponse api.SoundResponse
				if err := json.NewDecoder(resp.Body).Decode(&soundResponse); err != nil {
					log.Fatalf("error: could not parse soundx response: %s", err)
				} else if len(soundResponse.Data) < 1 {
					log.Fatalf("error: expected one sound in response, got none")
				} else {
					printSound(soundResponse.Data[0])
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

	var soundAddCmd = &cobra.Command{
		Use:   "add title",
		Short: "Add new sound",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			} else if len(args[0]) == 0 {
				return fmt.Errorf("invalid title: cannot be empty")
			} else {
				title = args[0]
			}
			return nil
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			if bpm <= 0 {
				log.Fatalf("invalid bpm: must be positive")
			}
			if duration <= 0 {
				log.Fatalf("invalid duration: must be positive")
			}
			for _, credit := range credits {
				xs := strings.Split(credit, ":")
				if len(xs) != 2 {
					log.Fatalf("invalid credit: must have the form role:name -- %q", credit)
				} else if len(xs[0]) == 0 {
					log.Fatalf("invalid role in credit: cannot be empty -- %q", credit)
				} else if len(xs[0]) == 0 {
					log.Fatalf("invalid name in credit: cannot be empty -- %q", credit)
				}
			}
			for _, genre := range genres {
				if len(genre) == 0 {
					log.Fatalf("invalid genre: cannot be empty -- %q", genre)
				}
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			creds := []api.Credit{}
			for _, credit := range credits {
				xs := strings.Split(credit, ":")
				creds = append(creds, api.Credit{
					Role: xs[0],
					Name: xs[1],
				})
			}

			input := api.NewSoundRequest{
				Data: []api.NewSound{
					{
						Title:             title,
						Bpm:               bpm,
						DurationInSeconds: duration,
						Credits:           creds,
						Genres:            genres,
					},
				},
			}

			resp, err := soundx.CreateSounds(ctx, input)
			if err != nil {
				log.Fatalf("error: soundx request failed: %s", err)
			} else if resp.StatusCode == 201 {
				var soundResponse api.SoundResponse
				if err := json.NewDecoder(resp.Body).Decode(&soundResponse); err != nil {
					log.Fatalf("error: could not parse soundx response: %s", err)
				} else if len(soundResponse.Data) < 1 {
					log.Fatalf("error: expected one sound in response, got none")
				} else {
					printSound(soundResponse.Data[0])
				}
			} else {
				log.Printf("error: http response code: %d", resp.StatusCode)
				buf := new(strings.Builder)
				if _, err := io.Copy(buf, resp.Body); err != nil {
					log.Fatalf("unexpected response code %d: not able to read response body: %s", resp.StatusCode, err)
				}
				log.Fatalf("unexpected response code %d: response body: %s", resp.StatusCode, buf.String())
			}
		},
	}
	soundAddCmd.Flags().StringSliceVarP(&credits, "credit", "c", []string{}, "sound credit in the form role:name")
	soundAddCmd.Flags().StringSliceVarP(&genres, "genre", "g", []string{}, "sound genre")
	soundAddCmd.MarkFlagRequired("genre")
	soundAddCmd.Flags().Int32VarP(&bpm, "bpm", "b", 0, "beats per minute of sound")
	soundAddCmd.MarkFlagRequired("bpm")
	soundAddCmd.Flags().Int32VarP(&duration, "duration", "d", 0, "duration of sound in seconds")
	soundAddCmd.MarkFlagRequired("duration")

	soundCmd.AddCommand(soundListCmd)
	soundCmd.AddCommand(soundGetCmd)
	soundCmd.AddCommand(soundAddCmd)
	soundCmd.AddCommand(soundRecommendCmd)
	return soundCmd
}
