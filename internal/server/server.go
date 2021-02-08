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
	players := strings.TrimPrefix(r.URL.Path, "/api/v1/players/")
	fmt.Fprint(w, osuSer.Data.GetPlayerStat(players))
}
