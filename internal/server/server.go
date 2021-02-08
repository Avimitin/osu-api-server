package server

import (
	"fmt"
	"net/http"
	"strings"
)

func OsuServer(w http.ResponseWriter, r *http.Request) {
	players := strings.TrimPrefix(r.URL.Path, "/api/v1/players/")
	fmt.Fprint(w, GetUserStat(players))
}

func GetUserStat(username string) string {
	if username == "avimitin" {
		return `{"username": "avimitin"}`
	}
	if username == "coooool" {
		return `{"username": "coooool"}`
	}
	return ""
}
