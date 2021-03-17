package server

import (
	"errors"
	"fmt"

	"github.com/avimitin/osu-api-server/internal/api"
)

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

// Err contained error type
type Err struct {
	// Level defined error's level
	// 0 is normal user level error, can response directly.
	// 1 is critical system level error, cannot response.
	Level int32
	e     error
}

// Error is used for error interface
func (e *Err) Error() string {
	return e.Error()
}

// CanResp return a boolean value about this error
// can be response to user or not.
func (e *Err) CanResp() bool {
	return e.Level == 0
}

// NewErr return a new error struct pointer
func NewErr(level int32, err error) *Err {
	return &Err{level, err}
}

// MakeErr make a new Err struct with a new error
func MakeErr(level int32, text string) *Err {
	return &Err{level, errors.New(text)}
}

// FmtErr return a new Err with formatted string
func FmtErr(level int32, text string, args ...interface{}) *Err {
	return &Err{level, fmt.Errorf(text, args...)}
}
