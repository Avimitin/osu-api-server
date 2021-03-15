package intergration

import (
	"testing"

	"github.com/avimitin/osu-api-server/internal/server"
)

func TestPrepareServer(t *testing.T) {
	TestServer, err := server.PrepareServer()
	if err != nil {
		t.Fatal(err)
	}
	if TestServer == nil {
		t.Errorf("got nil server")
	}
}
