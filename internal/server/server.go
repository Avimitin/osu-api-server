package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/avimitin/osu-api-server/internal/config"
	"github.com/avimitin/osu-api-server/internal/database"
)

// PrepareServer initialized all the service
func PrepareServer() (*OsuServer, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("server prepare config: %v", err)
	}
	db, err := database.Connect("mysql", cfg.DBSec.EncodeDSN())
	if err != nil {
		return nil, fmt.Errorf("connect to %s:%v", cfg.DBSec.EncodeDSN(), err)
	}
	if err = db.CheckUserDataStoreHealth(); err != nil {
		return nil, fmt.Errorf("check database health:%v", err)
	}
	opd := NewOsuPlayerData(db)
	return NewOsuServer(opd), nil
}

// OsuServer is a http handler and it store player data
type OsuServer struct {
	Data PlayerData
	http.Handler
}

// NewOsuServer return a OsuServer pointer
func NewOsuServer(store PlayerData) *OsuServer {
	if store == nil {
		panic("nil data store")
	}
	os := new(OsuServer)
	os.Data = store

	router := http.NewServeMux()
	router.Handle("/api/v1/player", http.HandlerFunc(os.playerHandler))
	router.Handle("/api/v1/recent", http.HandlerFunc(os.playerRecentHandler))
	os.Handler = router
	return os
}

func getUserFromForm(w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Content-Type", "application/json")
	player := r.PostForm.Get("player")
	log.Printf("%s:%s:player:%s", r.RemoteAddr, r.Method, player)
	return player
}

func (osuSer *OsuServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	player := getUserFromForm(w, r)
	if player == "" {
		returnServerError(w, errors.New("null user input"))
		return
	}
	stat, err := osuSer.Data.GetPlayerStat(player)
	if err != nil {
		log.Printf("get %s data: %v", player, err)
		returnServerError(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(stat)
	if err != nil {
		log.Printf("encode %v: %v", stat, err)
		returnServerError(w, err)
	}
}

func (osuSer *OsuServer) playerRecentHandler(w http.ResponseWriter, r *http.Request) {
	player := getUserFromForm(w, r)
	if player == "" {
		returnServerError(w, errors.New("no user specific"))
		return
	}
	score, err := osuSer.Data.GetPlayerRecent(player)
	if err != nil {
		returnServerError(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(score)
	if err != nil {
		log.Printf("encode %v:%v", score, err)
		returnServerError(w, fmt.Errorf("parse data failed %v", err))
	}
}

func returnServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
	json.NewEncoder(w).Encode(NewJsonMsg().Set("error", err.Error()))
}
