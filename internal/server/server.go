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
	// load user setting
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("server prepare config: %v", err)
	}
	// initialize query key
	err = api.KeyInit(cfg.Key)
	if err != nil {
		return nil, err
	}
	// initialize database connection
	dsn := cfg.DatabaseSettings.EncodeDSN(cfg.DBType)
	db, err := database.Connect(cfg.DBType, dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to %s at %s:%v", cfg.DBType, dsn, err)
	}
	// check database connection
	if err = db.CheckUserDataStoreHealth(); err != nil {
		return nil, fmt.Errorf("check %s health:%v", cfg.DBType, err)
	}
	return NewOsuServer(NewOsuPlayerData(db)), nil
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
	if !assertIsGetMethod(w, r) {
		return
	}
	setJsonHeader(w)
	player := getFormValue(r, "player")

	log.Printf("%q:%s:player:%q", r.RemoteAddr, r.Method, player)

	if isNullString(player) {
		responseServerError(w, nullInputErr)
		return
	}
	stat, err := osuSer.Data.GetPlayerStat(player)
	if err != nil {
		log.Printf("get %q data: %v", player, err)
		responseServerError(w, serverErr)
		return
	}
	err = json.NewEncoder(w).Encode(stat)
	if err != nil {
		log.Printf("encode %v: %v", stat, err)
		responseServerError(w, serverErr)
	}
}

func (osuSer *OsuServer) recentHandler(w http.ResponseWriter, r *http.Request) {
	if !assertIsGetMethod(w, r) {
		return
	}

	setJsonHeader(w)

	player := getFormValue(r, "player")
	if isNullString(player) {
		responseServerError(w, nullInputErr)
		return
	}

	perfect := getFormValue(r, "perfect")
	var perf = false
	if perfect == "true" {
		perf = true
	}

	mapID := getFormValue(r, "map")

	log.Printf("%q %s player %q recent map:%q perfect:%q", r.RemoteAddr, r.Method, player, mapID, perfect)

	score, err := osuSer.Data.GetRecent(player, mapID, perf)
	if err != nil {
		log.Printf("get recent data:%v", err)
		responseServerError(w, serverErr)
		return
	}

	err = json.NewEncoder(w).Encode(score)
	if err != nil {
		log.Printf("encode %v:%v", score, err)
		responseServerError(w, serverErr)
	}
}

func (osuSer *OsuServer) beatmapHandler(w http.ResponseWriter, r *http.Request) {
	if !assertIsGetMethod(w, r) {
		return
	}

	setJsonHeader(w)

	setID := getFormValue(r, "set_id")
	mapID := getFormValue(r, "map_id")

	log.Printf("%q %s beatmap [set:%q|bmap:%q]", r.RemoteAddr, r.Method, setID, mapID)

	bmap, err := osuSer.Data.GetBeatMaps(setID, mapID)
	if err != nil {
		log.Printf("get beatmap : %v", err)
		responseServerError(w, serverErr)
		return
	}
	err = json.NewEncoder(w).Encode(bmap)
	if err != nil {
		log.Printf("unmarshal beatmap %q: %v", bmap, err)
		responseServerError(w, serverErr)
	}
}

func responseServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
	e := json.NewEncoder(w).Encode(NewJsonMsg().Set("error", err.Error()))
	if e != nil {
		log.Printf("encoding error msg: %v", err)
	}
}
