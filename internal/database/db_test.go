package database

import "testing"

func TestConnect(t *testing.T) {
	cfg := config.Configuration{
		DBSec: {}
	}
	db, err := Connect(cfg)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}
}
