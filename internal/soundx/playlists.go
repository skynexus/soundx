package soundx

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/skynexus/soundx/api"
	"github.com/skynexus/soundx/internal/db"
)

func (s *server) CreatePlaylists(ec echo.Context) error {
	// Process new playlists
	if ec.Request().Body == nil {
		err := fmt.Errorf("bad request: nil body")
		ec.JSON(400, err)
		return nil
	}

	var newPlaylistRequest api.NewPlaylistRequest
	if err := json.NewDecoder(ec.Request().Body).Decode(&newPlaylistRequest); err != nil {
		err := fmt.Errorf("invalid json body: %w", err)
		ec.JSON(400, err)
		return nil
	}

	var playlists []db.Playlist
	for _, p := range newPlaylistRequest.Data {
		playlist, err := fromNewPlaylistDto(p)
		if err != nil {
			ec.JSON(400, err)
			return nil
		}
		playlists = append(playlists, playlist)
	}

	if err := s.r.AddPlaylists(playlists); err != nil {
		log.Printf("error: %s", err)
		ec.JSON(500, err)
		return nil
	}

	// Send response
	var dto []api.Playlist
	for _, p := range playlists {
		dto = append(dto, toPlaylistDto(p))
	}

	ec.JSON(201, api.PlaylistResponse{Data: dto})
	return nil
}

func (s *server) GetPlaylists(ec echo.Context) error {
	playlists, getErr := s.r.GetPlaylists()
	if getErr != nil {
		log.Printf("error: %s", getErr)
		ec.JSON(500, getErr)
		return nil
	}

	// Send response
	var dto []api.Playlist
	for _, p := range playlists {
		dto = append(dto, toPlaylistDto(p))
	}

	ec.JSON(200, api.PlaylistResponse{Data: dto})
	return nil
}
