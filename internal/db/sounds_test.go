package db

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type soundsTestSuite struct {
	Resources
}

func TestSoundsTestSuite(t *testing.T) {
	suite.Run(t, &soundsTestSuite{})
}

func (s *soundsTestSuite) TestAddGetSounds() {
	r := NewRepository(s.DB)

	sounds := []Sound{
		{
			Title:    "Party Time",
			Bpm:      140,
			Credits:  []Credit{},
			Duration: 240,
			Genres:   []string{},
		},
		{
			Title:    "Nervous Meditation",
			Bpm:      90,
			Credits:  []Credit{},
			Duration: 300,
			Genres:   []string{},
		},
	}

	// Add sounds
	addErr := r.AddSounds(sounds)
	s.Nil(addErr)
	for _, sound := range sounds {
		s.Greater(sound.Id, int64(0))
	}

	// Get sounds
	result, getErr := r.GetSounds()
	s.Nil(getErr)
	s.Equal(len(sounds), len(result))
	for i := range sounds {
		s.Equal(sounds[len(sounds)-i-1].Id, result[i].Id)
	}
}

func (s *playlistsTestSuite) TestGetRecommendedSounds() {
	r := NewRepository(s.DB)

	// Create sounds
	sounds := []Sound{
		{
			Title:    "Party Time",
			Bpm:      140,
			Credits:  []Credit{},
			Duration: 240,
			Genres:   []string{"pop", "reggae"},
		},
		{
			Title:    "Nervous Meditation",
			Bpm:      90,
			Credits:  []Credit{},
			Duration: 300,
			Genres:   []string{"electronic"},
		},
		{
			Title:    "Modern",
			Bpm:      110,
			Credits:  []Credit{},
			Duration: 200,
			Genres:   []string{"electronic"},
		},
		{
			Title:    "Hello",
			Bpm:      130,
			Credits:  []Credit{},
			Duration: 170,
			Genres:   []string{"blues", "pop"},
		},
		{
			Title:    "Excluded",
			Bpm:      90,
			Credits:  []Credit{},
			Duration: 130,
			Genres:   []string{"classic"},
		},
	}

	addErr := r.AddSounds(sounds)
	s.Nil(addErr)
	for _, sound := range sounds {
		s.Greater(sound.Id, int64(0))
	}

	// Assert sounds
	result, getErr := r.GetSounds()
	s.Nil(getErr)
	s.Equal(len(sounds), len(result))
	for i := range sounds {
		s.Equal(sounds[len(sounds)-i-1].Id, result[i].Id)
	}

	// Create playlist
	playlists := []Playlist{
		{
			Title:    "Hits",
			SoundIds: []int64{sounds[0].Id, sounds[1].Id},
		},
	}

	// Add playlists
	err := r.AddPlaylists(playlists)
	s.Nil(err)
	for _, playlist := range playlists {
		s.Greater(playlist.Id, int64(0))
	}

	// Get recommendations
	recommendations, recErr := r.GetRecommendedSounds(&playlists[0].Id)
	s.Nil(recErr)
	s.Equal(4, len(recommendations))
}
