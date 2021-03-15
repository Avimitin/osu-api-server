package server

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/avimitin/osu-api-server/internal/api"
	"github.com/avimitin/osu-api-server/internal/database"
)

// OsuPlayerData store database connection
type OsuPlayerData struct {
	db *database.OsuDB
}

// NewOsuPlayerData return a database controller which satisfy PlayerData interface
func NewOsuPlayerData(db *database.OsuDB) *OsuPlayerData {
	return &OsuPlayerData{db}
}

func (opd *OsuPlayerData) GetPlayerStat(name string) (*Player, error) {
	return getPlayerDataByName(name, opd.db)
}

func (opd *OsuPlayerData) GetRecent(name, mapID string, perf bool) (*api.RecentPlay, error) {
	if name == "" {
		return nil, errors.New("invalid name input")
	}
	// if map id and perfect is not specific, return latest play
	if mapID == "" && perf == false {
		scores, err := api.GetUserRecent(name, 1)

		if err != nil {
			return nil, fmt.Errorf("GetPlayerRecent: %v", err)
		}

		if len(scores) < 1 {
			return nil, errors.New("user not found")
		}
		return scores[0], nil
	}
	scores, err := api.GetUserRecent(name, 50)
	if err != nil {
		return nil, fmt.Errorf("GetPlayerRecent: %v", err)
	}

	switch {
	case mapID == "" && perf == true:
		for _, sc := range scores {
			if sc.Perfect == "1" {
				return sc, nil
			}
		}
	case perf == false:
		for _, sc := range scores {
			if sc.BeatmapID == mapID {
				return sc, nil
			}
		}
	case perf == true:
		for _, sc := range scores {
			if sc.BeatmapID == mapID && sc.Perfect == "1" {
				return sc, nil
			}
		}
	}
	return nil, fmt.Errorf("target not found")
}

func (*OsuPlayerData) GetBeatMaps(setID, mapID string) (*api.Beatmap, error) {
	maps, err := api.GetBeatMaps(setID, mapID)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, errors.New("map not found")
	}
	return maps[0], nil
}

func getPlayerDataByName(name string, db *database.OsuDB) (*Player, error) {
	u, e := api.GetUsers(name)
	if e != nil {
		return nil, e
	}
	if len(u) <= 0 {
		return nil, errors.New("user " + name + " not found")
	}
	user := u[0]
	lu, e := db.GetUserRecent(user.Username)
	if e != nil {
		if strings.Contains(e.Error(), "not found") {
			e = db.InsertNewUser(
				user.UserID, user.Username, user.Playcount, user.PpRank,
				user.PpRaw, user.Accuracy, user.TotalSecondsPlayed,
			)
			if e != nil {
				return nil, fmt.Errorf("insert user %s : %v", name, e)
			}
			log.Printf("inserted %s into database", user.Username)
		} else {
			return nil, fmt.Errorf("query user %s: %v", name, e)
		}
	}
	var diff *Different
	if lu != nil {
		diff, e = getUserDiff(user, "recent", &lu.Date)
	}
	if e != nil {
		return nil, e
	}

	return &Player{
		Data: user,
		Diff: diff,
	}, nil
}
