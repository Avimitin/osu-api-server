package main

import (
	"log"
	"net/http"

	"github.com/avimitin/osuapi/internal/server"
)

func main() {
	// cast
	handler := http.HandlerFunc(server.OsuServer)

	if err := http.ListenAndServe(":11451", handler); err != nil {
		log.Fatalf("handle %s : %v", ":11451", err)
	}
}
