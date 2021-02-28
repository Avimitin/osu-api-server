package intergration

import (
	"testing"

	"github.com/avimitin/osu-api-server/internal/server"
)

func TestPrepareServer(t *testing.T) {
	server, err := server.PrepareServer()
	if err != nil {
		t.Fatal(err)
	}
	if server == nil {
		t.Errorf("got nil server")
	}
}
