package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type OsuDB struct {
	Conn *sql.DB
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

	return &OsuDB{Conn: db}, nil
}

func (db *OsuDB) InitTable() error {
	return initUserTable(db.Conn)
}

type User struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	PlayCount    string `json:"play_count"`
	Rank         string `json:"rank"`
	PP           string `json:"pp"`
	Acc          string `json:"acc"`
	TotalPlay    string `json:"total_play"`
	PcYtd        string `json:"play_count_ytd"`
	RankYtd      string `json:"rank_ytd"`
	PpYtd        string `json:"pp_ytd"`
	AccYtd       string `json:"acc_ytd"`
	TotalPlayYtd string `json:"ttp_ytd"`
}

func (db *OsuDB) GetUserRecent(username string) (*User, error) {
	const query = "SELECT username, playcount, rank, pp, acc FROM users WHERE username = ? OR user_id = ?"
	var u *User
	stmtOut, err := db.Conn.Prepare(query)
	err = stmtOut.QueryRow(username).Scan(
		&u.Username, &u.PlayCount, &u.PlayCount, &u.Rank, &u.PP, &u.Acc)
	if err != nil {
		return nil, fmt.Errorf("query %s : %v", query, err)
	}
	if u == nil {
		return nil, fmt.Errorf("user %s not found", username)
	}
	return u, nil
}

func (db *OsuDB) InsertNewUser(
	userID string, username string, pc string, rank string, pp string, acc string,
) error {

	const query = `
INSERT INTO users (
	user_id, username, playcount, rank, pp, acc
) VALUES (
	?,?,?,?,?,?
)
`
	stmtIn, err := db.Conn.Prepare(query)
	if err != nil {
		return fmt.Errorf("query %s : %v", query, err)
	}
	res, err := stmtIn.Exec(userID, username, pc, rank, pp, acc)
	if err != nil {
		return fmt.Errorf("insert %s: %v", query, err)
	}
	af, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if af < 1 {
		return errors.New("not row affected")
	}
	return nil
}

func initUserTable(db *sql.DB) error {
	const userTable = `
CREATE TABLE IF NOT EXISTS users(
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
