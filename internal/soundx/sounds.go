package soundx

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/skynexus/soundx/api"
	"github.com/skynexus/soundx/internal/db"
)

func (s *server) CreateSounds(ec echo.Context) error {
	// Process new sounds
	if ec.Request().Body == nil {
		err := fmt.Errorf("bad request: nil body")
		ec.JSON(400, err)
		return nil
	}

	var newSoundRequest api.NewSoundRequest
	if err := json.NewDecoder(ec.Request().Body).Decode(&newSoundRequest); err != nil {
		err := fmt.Errorf("invalid json body: %w", err)
		ec.JSON(400, err)
		return nil
	}

	var sounds []db.Sound
	for _, s := range newSoundRequest.Data {
		sounds = append(sounds, fromNewSoundDto(s))
	}

	if err := s.r.AddSounds(sounds); err != nil {
		log.Printf("error: %s", err)
		ec.JSON(500, err)
		return nil
	}

	// Send response
	var dto []api.Sound
	for _, s := range sounds {
		dto = append(dto, toSoundDto(s))
	}

	ec.JSON(201, api.SoundResponse{Data: dto})
	return nil
}

func (s *server) GetSounds(ec echo.Context) error {
	sounds, getErr := s.r.GetSounds()
	if getErr != nil {
		log.Printf("error: %s", getErr)
		ec.JSON(500, getErr)
		return nil
	}

	// Send response
	var dto []api.Sound
	for _, s := range sounds {
		dto = append(dto, toSoundDto(s))
	}

	ec.JSON(200, api.SoundResponse{Data: dto})
	return nil
}

func (s *server) GetSound(ec echo.Context, id string) error {
	soundId, parseErr := parseId("sound", id)
	if parseErr != nil {
		return parseErr
	}

	sounds, getErr := s.r.GetSoundsByIds([]int64{soundId})
	if getErr != nil {
		log.Printf("error: %s", getErr)
		ec.JSON(500, getErr)
		return nil
	} else if len(sounds) < 1 {
		ec.JSON(404, "sound not found")
		return nil
	}

	// Send response
	var dto []api.Sound
	for _, s := range sounds {
		dto = append(dto, toSoundDto(s))
	}

	ec.JSON(200, api.SoundResponse{Data: dto})
	return nil
}

func (s *server) GetRecommendedSounds(ec echo.Context, params api.GetRecommendedSoundsParams) error {
	var playlistId *int64

	if params.PlaylistId != nil {
		if id, parseErr := parseId("playlist", *params.PlaylistId); parseErr != nil {
			return parseErr
		} else {
			playlistId = &id
		}
	}

	sounds, getErr := s.r.GetRecommendedSounds(playlistId)
	if getErr != nil {
		log.Printf("error: %s", getErr)
		ec.JSON(500, getErr)
		return nil
	}

	// Send response
	var dto []api.Sound
	for _, s := range sounds {
		dto = append(dto, toSoundDto(s))
	}

	ec.JSON(200, api.SoundResponse{Data: dto})
	return nil
}
