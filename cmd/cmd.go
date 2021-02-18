package main

import (
	"log"
	"net/http"

	"github.com/avimitin/osu-api-server/internal/server"
)

func main() {
	var err error
	var OsuSer *server.OsuServer
	if OsuSer, err = server.PrepareServer(); err != nil {
		log.Fatalf("preparing server:%v", err)
	}

	log.Println("server listening on port 11451")
	if err = http.ListenAndServe(":11451", OsuSer); err != nil {
		log.Fatalf("handle %s : %v", ":11451", err)
	}
}
