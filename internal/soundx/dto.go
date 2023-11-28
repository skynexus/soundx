package soundx

import (
	"strconv"

	"github.com/skynexus/soundx/api"
	"github.com/skynexus/soundx/internal/db"
)

func fromNewSoundDto(dto api.NewSound) db.Sound {
	var creds db.Credits

	for _, c := range dto.Credits {
		creds = append(creds, db.Credit{
			Name: c.Name,
			Role: c.Role,
		})
	}

	return db.Sound{
		Bpm:      dto.Bpm,
		Credits:  creds,
		Duration: dto.DurationInSeconds,
		Genres:   dto.Genres,
		Title:    dto.Title,
	}
}

func toSoundDto(sound db.Sound) api.Sound {
	var creds []api.Credit

	for _, c := range sound.Credits {
		creds = append(creds, api.Credit{
			Name: c.Name,
			Role: c.Role,
		})
	}

	return api.Sound{
		Id:                strconv.FormatInt(sound.Id, 10),
		Bpm:               sound.Bpm,
		Credits:           creds,
		DurationInSeconds: sound.Duration,
		Genres:            sound.Genres,
		Title:             sound.Title,
		CreatedAt:         sound.CreatedAt,
		UpdatedAt:         sound.UpdatedAt,
	}
}

func fromNewPlaylistDto(dto api.NewPlaylist) (db.Playlist, error) {
	var soundIds []int64

	for _, s := range dto.Sounds {
		soundId, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return db.Playlist{}, err
		}
		soundIds = append(soundIds, soundId)
	}

	return db.Playlist{
		Title:    dto.Title,
		SoundIds: soundIds,
	}, nil
}

func toPlaylistDto(playlist db.Playlist) api.Playlist {
	var soundIds []string

	for _, soundId := range playlist.SoundIds {
		soundIds = append(soundIds, strconv.FormatInt(soundId, 10))
	}

	return api.Playlist{
		Id:     strconv.FormatInt(playlist.Id, 10),
		Title:  playlist.Title,
		Sounds: soundIds,
	}
}
