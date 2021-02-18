package server

import "github.com/avimitin/osu-api-server/internal/api"

type PlayerData interface {
	GetPlayerStat(name string) (string, error)
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
