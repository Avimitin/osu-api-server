package database

import (
	"testing"

	"github.com/avimitin/osuapi/internal/config"
)

func connect(t *testing.T) *OsuDB {
	cfg, err := config.GetConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DBSec.Username == "" {
		t.Fatal("no database config")
	}
	db, err := Connect(cfg.DBSec.EncodeDSN())
	if err != nil {
		t.Fatal(err)
	}
	err = db.DB.Ping()
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestConnect(t *testing.T) {
	connect(t)
}

func TestUserTable(t *testing.T) {
	db := connect(t)
	err := db.InitTable()
	if err != nil {
		t.Fatal(err)
	}
	user, err := db.GetUser("avimitin")
	if err != nil {
		t.Fatal(err)
	}
	want := "avimitin"
	if user.Username != want {
		t.Fatalf("get wrong user %s", user.Username)
	}
	if currentName := user.GetCurrentStats().Username; currentName != want {
		t.Fatalf("want %s, get %s", want, currentName)
	}
}
