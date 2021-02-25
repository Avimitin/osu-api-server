package server

import "github.com/avimitin/osu-api-server/internal/api"

type OsuData interface {
	GetPlayerStat(name string) (*Player, error)
	GetRecent(name, mapID string, perfect bool) (*api.RecentPlay, error)
	GetBeatMaps(setID, mapID string) (*api.Beatmap, error)
}

type Player struct {
	Data *api.User  `json:"latest_data"`
	Diff *Different `json:"diff"`
}

type Different struct {
	PlayCount string `json:"play_count"`
	Rank      string `json:"rank"`
	PP        string `json:"pp"`
	Acc       string `json:"acc"`
	TotalPlay string `json:"total_play"`
}
