package server

import "github.com/avimitin/osu-api-server/internal/api"

// OsuData handle osu data
type OsuData interface {
	GetPlayerStat(name string) (*Player, error)
	GetRecent(name, mapID string, perfect bool) (*api.RecentPlay, error)
	GetBeatMaps(setID, mapID string) (*api.Beatmap, error)
}

// Player store user information and different between
// each request
type Player struct {
	Data *api.User  `json:"latest_data"`
	Diff *Different `json:"diff"`
}

// Different specific user data different
type Different struct {
	PlayCount string `json:"play_count"`
	Rank      string `json:"rank"`
	PP        string `json:"pp"`
	Acc       string `json:"acc"`
	TotalPlay string `json:"total_play"`
}
