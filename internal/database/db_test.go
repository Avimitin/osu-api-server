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

func TestGetUser(t *testing.T) {
	db := connect(t)
	_, err := db.Conn.Exec(`
INSERT INTO users (user_id, username, playcount, rank, pp, acc, total_play)
VALUES ('114514', 'avimitin', '100', '1', '999', '99', '1919810')
`)
	if err != nil {
		t.Fatal(err)
	}
	user, err := db.GetUserRecent("avimitin")
	if err != nil {
		t.Fatal(err)
	}
	if user.Acc != "99" {
		t.Errorf("get %+v is not wanted", user)
	}
	t.Log("clean")
	_, err = db.Conn.Exec(`
DELETE FROM users WHERE user_id = '114514';
	`)
	if err != nil {
		t.Errorf("Clean failed: %v", err)
	}
}
