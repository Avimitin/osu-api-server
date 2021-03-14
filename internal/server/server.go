package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/avimitin/osu-api-server/internal/api"
	"github.com/avimitin/osu-api-server/internal/config"
	"github.com/avimitin/osu-api-server/internal/database"
)

var (
	serverErr    = errors.New("unexpected server error")
	nullInputErr = errors.New("invalid null input")
)

// PrepareServer initialized all the service
func PrepareServer() (*OsuServer, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("server prepare config: %v", err)
	}
	api.KeyInit(cfg.Key)
	db, err := database.Connect("redis", cfg.DatabaseSettings.EncodeRedisDSN())
	if err != nil {
		return nil, fmt.Errorf("connect to %s:%v", cfg.DatabaseSettings.EncodeRedisDSN(), err)
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

	log.Printf("%q:%s:player:%q", r.RemoteAddr, r.Method, player)

	if player == "" {
		serErr(w, nullInputErr)
		return
	}
	stat, err := osuSer.Data.GetPlayerStat(player)
	if err != nil {
		log.Printf("get %q data: %v", player, err)
		serErr(w, serverErr)
		return
	}
	err = json.NewEncoder(w).Encode(stat)
	if err != nil {
		log.Printf("encode %v: %v", stat, err)
		serErr(w, serverErr)
	}
}

func (osuSer *OsuServer) recentHandler(w http.ResponseWriter, r *http.Request) {
	setJsonHeader(w)
	player := getFormValue(r, "player")
	if player == "" {
		serErr(w, nullInputErr)
		return
	}
	perfect := getFormValue(r, "perfect")
	var perf = false
	if perfect == "true" {
		perf = true
	}
	mapID := getFormValue(r, "map")

	log.Printf("%q %s [%q]:%q:%q", r.RemoteAddr, r.Method, player, mapID, perfect)

	score, err := osuSer.Data.GetRecent(player, mapID, perf)
	if err != nil {
		log.Printf("get recent data:%v", err)
		serErr(w, serverErr)
		return
	}
	err = json.NewEncoder(w).Encode(score)
	if err != nil {
		log.Printf("encode %v:%v", score, err)
		serErr(w, serverErr)
	}
}

func (osuSer *OsuServer) beatmapHandler(w http.ResponseWriter, r *http.Request) {
	setJsonHeader(w)
	setID := getFormValue(r, "set_id")
	mapID := getFormValue(r, "map_id")

	log.Printf("%q %s beatmap [%q|%q]", r.RemoteAddr, r.Method, setID, mapID)

	bmap, err := osuSer.Data.GetBeatMaps(setID, mapID)
	if err != nil {
		log.Printf("get beatmap : %v", err)
		serErr(w, serverErr)
		return
	}
	err = json.NewEncoder(w).Encode(bmap)
	if err != nil {
		log.Printf("unmarshal beatmap %q: %v", bmap, err)
		serErr(w, serverErr)
	}
}

func serErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
	json.NewEncoder(w).Encode(NewJsonMsg().Set("error", err.Error()))
}
