package server

import (
	"fmt"
	"net/http"
)

func OsuServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"username": "avimitin"}`)
}
