package server

import (
	"encoding/json"
	"log"
	"net/http"
)

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
	if !isWantedMethod(w, r, http.MethodPost) {
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
	if !isWantedMethod(w, r, http.MethodPost) {
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
	if !isWantedMethod(w, r, http.MethodPost) {
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
