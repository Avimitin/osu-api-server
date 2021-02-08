package server

import (
	"fmt"
	"net/http"
	"strings"
)

func OsuServer(w http.ResponseWriter, r *http.Request) {
	players := strings.TrimPrefix(r.URL.Path, "/api/v1/players/")
	if players == "avimitin" {
		fmt.Fprint(w, `{"username": "avimitin"}`)
		return
	}
	if players == "coooool" {
		fmt.Fprint(w, `{"username": "coooool"}`)
		return
	}
}
