package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/avimitin/osu-api-server/internal/api"
	"github.com/avimitin/osu-api-server/internal/config"
	"github.com/avimitin/osu-api-server/internal/database"
)

var (
	db *database.OsuDB
)

// PrepareServer initialized all the service
func PrepareServer() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("server prepare config: %v", err)
	}
	db, err = database.Connect("mysql", cfg.DBSec.EncodeDSN())
	if err != nil {
		return fmt.Errorf("connect to %s:%v", cfg.DBSec.EncodeDSN(), err)
	}
	if err = db.CheckUserDataStoreHealth(); err != nil {
		return fmt.Errorf("check database health:%v", err)
	}
	return nil
}

type OsuServer struct {
	Data PlayerData
	http.Handler
}

func NewOsuServer(store PlayerData) *OsuServer {
	if store == nil {
		panic("nil data store")
	}
	os := new(OsuServer)
	os.Data = store

	router := http.NewServeMux()
	router.Handle("/api/v1/player", http.HandlerFunc(os.playerHandler))
	os.Handler = router
	return os
}

func (osuSer *OsuServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	player := r.PostForm.Get("player")
	stat, err := osuSer.Data.GetPlayerStat(player)
	if err != nil {
		log.Printf("get %s data: %v", player, err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
		return
	}
	fmt.Fprint(w, stat)
}

type OsuPlayerData struct{}

func (opd *OsuPlayerData) GetPlayerStat(name string) (string, error) {
	p, e := getPlayerDataByName(name)
	if e != nil {
		return "", e
	}
	data, e := json.Marshal(p)
	if e != nil {
		return "", e
	}
	return string(data), nil
}

func getPlayerDataByName(name string) (*Player, error) {
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
