package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

type Credit struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type Credits []Credit

// Value JSON-encodes the Credits list.
func (c Credits) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan decodes a JSON-encoded Credits list.
func (c Credits) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

type Sound struct {
	Id        int64
	Title     string
	Duration  int32
	Bpm       int32
	Genres    []string
	Credits   Credits
	CreatedAt time.Time
	UpdatedAt time.Time
}

func validateSoundIds(sounds []Sound, ids []int64) error {
	present := make(map[int64]struct{})
	for _, sound := range sounds {
		present[sound.Id] = struct{}{}
	}
	var missing []int64
	for _, id := range ids {
		if _, ok := present[id]; ok {
			continue
		}
		missing = append(missing, id)
	}
	if len(missing) > 0 {
		return fmt.Errorf("sounds not found: %v", missing)
	}
	return nil
}

func soundSelectBuilder() sq.SelectBuilder {
	return sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select(
			"sound.id",
			"sound.title",
			"sound.duration",
			"sound.bpm",
			"sound.genres",
			"sound.credits",
			"sound.created_at",
			"sound.updated_at",
		).
		From("sounds AS sound")
}

func soundScan(row sq.RowScanner) (*Sound, error) {
	var sound Sound
	err := row.Scan(
		&sound.Id,
		&sound.Title,
		&sound.Duration,
		&sound.Bpm,
		pq.Array(&sound.Genres),
		&sound.Credits,
		&sound.CreatedAt,
		&sound.UpdatedAt,
	)
	return &sound, err
}

func (r *Repository) AddSound(prototype *Sound) error {
	return r.AddSounds([]Sound{*prototype})
}

func (r *Repository) AddSounds(prototypes []Sound) error {
	now := time.Now().UTC()

	insertQuery := pgStatementBuilder().Insert("sounds").
		Columns(
			"title",
			"duration",
			"bpm",
			"genres",
			"credits",
			"created_at",
			"updated_at")

	var ts []time.Time
	for _, prototype := range prototypes {
		ts = append(ts, now)
		insertQuery = insertQuery.Values(
			prototype.Title,
			prototype.Duration,
			prototype.Bpm,
			pq.Array(prototype.Genres),
			prototype.Credits,
			now,
			now)
		now = now.Add(time.Millisecond)
	}

	insertQuery = insertQuery.Suffix("RETURNING id")

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
}

func (r *Repository) GetSounds() ([]Sound, error) {
	rows, err := soundSelectBuilder().
		OrderBy("sound.created_at DESC").
		RunWith(r.Db()).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsScan(rows, soundScan)
}

func (r *Repository) GetSoundsByIds(ids []int64) ([]Sound, error) {
	rows, err := soundSelectBuilder().
		Where(sq.Eq{"sound.id": ids}).
		OrderBy("sound.created_at DESC").
		RunWith(r.Db()).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsScan(rows, soundScan)
}

func (r *Repository) GetSoundsByGenres(genres []string) ([]Sound, error) {
	// SELECT sound.* FROM sounds AS sound WHERE sound.genres <@ '{"jazz","pop"}'::TEXT[];
	rows, err := soundSelectBuilder().
		Where("sound.genres && ?", pq.Array(genres)).
		OrderBy("sound.created_at DESC").
		RunWith(r.Db()).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsScan(rows, soundScan)
}

func (r *Repository) GetRecommendedSounds(playlistId *int64) ([]Sound, error) {
	if playlistId != nil {
		var recommendations []Sound

		err := r.Transactional(func() error {
			if genres, pErr := r.GetPlaylistGenres(*playlistId); pErr != nil {
				return pErr
			} else if sounds, sErr := r.GetSoundsByGenres(genres); sErr != nil {
				return sErr
			} else {
				recommendations = append(recommendations, sounds...)
			}
			return nil
		})

		return recommendations, err
	} else {
		return r.GetSounds()
	}
}
