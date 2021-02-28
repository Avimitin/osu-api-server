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
	db, err := database.Connect("redis", cfg.DBSec.EncodeRedisDSN())
	if err != nil {
		return nil, fmt.Errorf("connect to %s:%v", cfg.DBSec.EncodeRedisDSN(), err)
	}
	if err = db.CheckUserDataStoreHealth(); err != nil {
		return nil, fmt.Errorf("check database health:%v", err)
	}
	opd := NewOsuPlayerData(db)
	return NewOsuServer(opd), nil
}

// OsuServer is a http handler and it store player data
type OsuServer struct {
	Data OsuData
	http.Handler
}

// NewOsuServer return a OsuServer pointer
func NewOsuServer(store OsuData) *OsuServer {
	if store == nil {
		panic("nil data store")
	}
	os := new(OsuServer)
	os.Data = store

	router := http.NewServeMux()
	router.Handle("/api/v1/player", http.HandlerFunc(os.playerHandler))
	router.Handle("/api/v1/recent", http.HandlerFunc(os.recentHandler))
	router.Handle("/api/v1/beatmap", http.HandlerFunc(os.beatmapHandler))
	os.Handler = router
	return os
}

func (osuSer *OsuServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	setJsonHeader(w)
	player := getFormValue(r, "player")
	log.Printf("%s:%s:player:%s", r.RemoteAddr, r.Method, player)
	if player == "" {
		serErr(w, errors.New("null user input"))
		return
	}
	stat, err := osuSer.Data.GetPlayerStat(player)
	if err != nil {
		log.Printf("get %s data: %v", player, err)
		serErr(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(stat)
	if err != nil {
		log.Printf("encode %v: %v", stat, err)
		serErr(w, err)
	}
}

func (osuSer *OsuServer) recentHandler(w http.ResponseWriter, r *http.Request) {
	setJsonHeader(w)
	player := getFormValue(r, "player")
	if player == "" {
		serErr(w, errors.New("no user specific"))
		return
	}
	perfect := getFormValue(r, "perfect")
	var perf = false
	if perfect == "true" {
		perf = true
	}
	mapID := getFormValue(r, "map")

	log.Printf("recentHandler:[%s]:%s:%s", player, mapID, perfect)

	score, err := osuSer.Data.GetRecent(player, mapID, perf)
	if err != nil {
		serErr(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(score)
	if err != nil {
		log.Printf("encode %v:%v", score, err)
		serErr(w, fmt.Errorf("parse data failed %v", err))
	}
}

func (osuSer *OsuServer) beatmapHandler(w http.ResponseWriter, r *http.Request) {
	setJsonHeader(w)
	setID := getFormValue(r, "set_id")
	mapID := getFormValue(r, "map_id")
	bmap, err := osuSer.Data.GetBeatMaps(setID, mapID)
	if err != nil {
		serErr(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(bmap)
	if err != nil {
		log.Printf("unmarshal beatmap %s: %v", bmap, err)
		serErr(w, fmt.Errorf("parse data failed %v", err))
	}
}

func serErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
	json.NewEncoder(w).Encode(NewJsonMsg().Set("error", err.Error()))
}
