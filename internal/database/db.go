package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type OsuDB struct {
	DB *sql.DB
}

func Connect(dsn string) (*OsuDB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect %s:%v", dsn, err)
	}
	// Set limit
	if db != nil {
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
	}

	return &OsuDB{DB: db}, nil
}

func (db *OsuDB) InitTable() error {
	return initUserTable(db.DB)
}

type User struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	PlayCount string `json:"play_count"`
	Rank      string `json:"rank"`
	PP        string `json:"pp"`
	Acc       string `json:"acc"`
}

func (db *OsuDB) GetUserRecent(username string) (*User, error) {
	const query = "SELECT username, playcount, rank, pp, acc FROM users WHERE username = ? OR user_id = ?"
	var u *User
	stmtOut, err := db.DB.Prepare(query)
	err = stmtOut.QueryRow(username).Scan(&u)
	if err != nil {
		return nil, fmt.Errorf("query %s : %v", query, err)
	}
	if u == nil {
		return nil, fmt.Errorf("user %s not found", username)
	}
	return u, nil
}

func (db *OsuDB) InsertNewUser(userID string, username string, pc string, rank string, pp string, acc string) (*User, error) {
	const query = `INSERT INTO users (
	user_id, username, playcount, rank, pp, acc
) VALUES (
	?,?,?,?,?,?
)
`
	stmtIn, err := db.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("query %s : %v", query, err)
	}
	_, err = stmtIn.Exec(userID, username, pc, rank, pp, acc)
	if err != nil {
		return nil, fmt.Errorf("insert %s: %v", query, err)
	}
	return &User{
		UserID:    userID,
		Username:  username,
		PlayCount: pc,
		Rank:      rank,
		PP:        pp,
		Acc:       acc,
	}, nil
}

func initUserTable(db *sql.DB) error {
	const userTable = `
CREATE TABLE IF NOT EXIST users(
	id INT AUTO_INCREMENT,
	user_id VARCHAR(18) NOT NULL,
	username VARCHAR(255) NOT NULL,
	playcount VARCHAR(18),
	rank VARCHAR(18),
	pp VARCHAR(18),
	acc VARCHAR(18),
	total_play VARCHAR(18),
	playcount_ytd VARCHAR(18),
	rank_ytd VARCHAR(18),
	pp_ytd VARCHAR(18),
	acc_ytd VARCHAR(18),
	total_play_ytd VARCHAR(18),
	PRIMARY KEY (id)
)CHARSET=utf8mb4
	`
	_, err := db.Exec(userTable)
	if err != nil {
		return fmt.Errorf("init table: %s:%v", userTable, err)
	}

	return nil
}
