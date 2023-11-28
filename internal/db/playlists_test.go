package db

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/suite"
)

type playlistsTestSuite struct {
	Resources
}

func TestPlaylistsTestSuite(t *testing.T) {
	suite.Run(t, &playlistsTestSuite{})
}

func (s *playlistsTestSuite) TestAddGetPlaylist() {
	r := NewRepository(s.DB)

	playlists := []Playlist{
		{
			Title:    "Mellow Collection",
			SoundIds: []int64{},
		},
		{
			Title:    "Dance Moves",
			SoundIds: []int64{},
		},
	}

	// Add playlists
	err := r.AddPlaylists(playlists)
	s.Nil(err)
	for _, playlist := range playlists {
		s.Greater(playlist.Id, int64(0))
	}

	// Get playlists
	result, getErr := r.GetPlaylists()
	s.Nil(getErr)
	s.Equal(len(playlists), len(result))
	for i := range playlists {
		s.Equal(playlists[len(playlists)-i-1].Id, result[i].Id)
	}

	// Get specific playlist
	playlist, getErr2 := r.GetPlaylist(playlists[0].Id)
	s.Nil(getErr2)
	s.NotNil(playlist)
	s.Equal(playlist.Id, playlists[0].Id)
}

func (s *playlistsTestSuite) TestGetPlaylistGenres() {
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
			Genres:   []string{"electronic", "pop"},
		},
		{
			Title:    "Around the Corner",
			Bpm:      95,
			Credits:  []Credit{},
			Duration: 160,
			Genres:   []string{"jazz"},
		},
	}

	addErr := r.AddSounds(sounds)
	s.Nil(addErr)
	for _, sound := range sounds {
		s.Greater(sound.Id, int64(0))
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
	genres, recErr := r.GetPlaylistGenres(playlists[0].Id)
	s.Nil(recErr)
	s.Equal(3, len(genres))
	s.True(slices.Contains(genres, "pop"), "expected pop genre")
	s.True(slices.Contains(genres, "reggae"), "expected reggae genre")
	s.True(slices.Contains(genres, "electronic"), "expected electronic genre")
	s.False(slices.Contains(genres, "jazz"), "did not expect jazz genre")
}
