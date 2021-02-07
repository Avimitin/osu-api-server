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
	err = db.Conn.Ping()
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
	rows, err := db.Conn.Query("SHOW TABLES;")
	if err != nil {
		t.Fatal(err)
	}
	var tables []string
	for rows.Next() {
		var table string
		rowErr := rows.Scan(&table)
		if rowErr != nil {
			t.Fatal(err)
		}
		tables = append(tables, table)
	}

	var want = "users"
	exist := false
	if len(tables) == 0 {
		t.Errorf("no table")
	} else {
		for _, t := range tables {
			if t == want {
				exist = true
			}
		}
	}
	if !exist {
		t.Errorf("%s not found, got %v", want, tables)
	}
}
