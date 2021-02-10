package main

import (
	"log"
	"net/http"

	"github.com/avimitin/osuapi/internal/server"
)

func main() {
	// cast
	OsuSer := &server.OsuServer{
		Data: &server.OsuPlayerData{},
	}

	log.Println("server listening on port 11451")
	if err := http.ListenAndServe(":11451", OsuSer); err != nil {
		log.Fatalf("handle %s : %v", ":11451", err)
	}
}
