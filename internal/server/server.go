package server

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerData interface {
	GetPlayerStat(name string) string
}

type OsuServer struct {
	Data PlayerData
}

func (osuSer *OsuServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/api/v1/players/")
	stat := osuSer.Data.GetPlayerStat(player)
	if stat == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, stat)
}
