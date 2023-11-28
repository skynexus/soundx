package db

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

type Playlist struct {
	Id        int64
	Title     string
	SoundIds  []int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func playlistSelectBuilder() sq.SelectBuilder {
	return sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select(
			"playlist.id",
			"playlist.title",
			"playlist.sound_ids",
			"playlist.created_at",
			"playlist.updated_at",
		).
		From("playlists AS playlist")
}

func playlistScan(row sq.RowScanner) (*Playlist, error) {
	var playlist Playlist
	err := row.Scan(
		&playlist.Id,
		&playlist.Title,
		pq.Array(&playlist.SoundIds),
		&playlist.CreatedAt,
		&playlist.UpdatedAt,
	)
	return &playlist, err
}

func (r *Repository) AddPlaylists(prototypes []Playlist) error {
	now := time.Now().UTC()

	insertQuery := pgStatementBuilder().
		Insert("playlists").
		Columns(
			"title",
			"sound_ids",
			"created_at",
			"updated_at")

	var ts []time.Time
	var soundIds []int64
	for _, prototype := range prototypes {
		ts = append(ts, now)
		insertQuery = insertQuery.Values(
			prototype.Title,
			pq.Array(prototype.SoundIds),
			now,
			now)
		soundIds = append(soundIds, prototype.SoundIds...)
		now = now.Add(time.Millisecond)
	}

	insertQuery = insertQuery.Suffix("RETURNING id")

	return r.Transactional(func() error {
		if sounds, getErr := r.GetSoundsByIds(soundIds); getErr != nil {
			return getErr
		} else if validationErr := validateSoundIds(sounds, soundIds); validationErr != nil {
			return validationErr
		}

		rows, insertErr := insertQuery.RunWith(r.Db()).Query()
		if insertErr != nil {
			return insertErr
		}
		defer rows.Close()

		ids, scanErr := rowsScan(rows, idScan)
		if scanErr != nil {
			return scanErr
		}

		for i := range prototypes {
			prototypes[i].Id = ids[i]
			prototypes[i].CreatedAt = ts[i]
			prototypes[i].UpdatedAt = ts[i]
		}
		return nil
	})
}

func (r *Repository) GetPlaylists() ([]Playlist, error) {
	rows, err := playlistSelectBuilder().
		OrderBy("created_at DESC").
		RunWith(r.Db()).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsScan(rows, playlistScan)
}

func (r *Repository) GetPlaylist(id int64) (*Playlist, error) {
	row := playlistSelectBuilder().
		Where(sq.Eq{"playlist.id": id}).
		RunWith(r.Db()).QueryRow()
	return playlistScan(row)
}

func (r *Repository) GetPlaylistGenres(playlistId int64) ([]string, error) {
	// SELECT sound.genres FROM sounds AS sound JOIN playlists AS playlist ON playlist.id = 11 WHERE sound.id = ANY(playlist.sound_ids);
	rows, qErr := soundSelectBuilder().
		Join("playlists AS playlist ON playlist.id = ?", playlistId).
		Where("sound.id = ANY(playlist.sound_ids)").
		RunWith(r.Db()).Query()
	if qErr != nil {
		return nil, qErr
	}
	defer rows.Close()

	sounds, sErr := rowsScan(rows, soundScan)
	if sErr != nil {
		return nil, sErr
	}

	var genres []string
	cache := make(map[string]struct{})
	for _, sound := range sounds {
		for _, genre := range sound.Genres {
			cache[genre] = struct{}{}
		}
	}
	for genre := range cache {
		genres = append(genres, genre)
	}

	return genres, nil
}
