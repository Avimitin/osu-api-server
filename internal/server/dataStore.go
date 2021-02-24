package server

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/avimitin/osu-api-server/internal/api"
	"github.com/avimitin/osu-api-server/internal/database"
)

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

func (opd *OsuPlayerData) GetPlayerRecent(name string) (*api.RecentPlay, error) {
	scores, err := api.GetUserRecent(name, 1)
	if err != nil {
		return nil, fmt.Errorf("GetPlayerRecent: %v", err)
	}
	return scores[0], nil
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
		if strings.Contains(e.Error(), "user") {
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

	p := &Player{
		Data: u[0],
		Diff: diff,
	}
	return p, nil
}
