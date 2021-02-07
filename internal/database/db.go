package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//OsuDB contain sql.DB field
type OsuDB struct {
	Conn *sql.DB
}

// Connect return database connection by given DSN
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

// InitTable initialize table at setup
func (db *OsuDB) InitTable() error {
	return initUserTable(db.Conn)
}

// User type contain user field
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

// GetUserRecent return user data with given name
func (db *OsuDB) GetUserRecent(username string) (*User, error) {
	const query = "SELECT username, playcount, rank, pp, acc, total_play FROM users WHERE username = ? OR user_id = ?"
	u := &User{}
	stmtOut, err := db.Conn.Prepare(query)
	defer stmtOut.Close()
	err = stmtOut.QueryRow(username, username).Scan(
		&u.Username, &u.PlayCount, &u.Rank, &u.PP, &u.Acc, &u.TotalPlay)
	if err != nil {
		return nil, fmt.Errorf("query %s : %v", query, err)
	}
	if u == nil {
		return nil, fmt.Errorf("user %s not found", username)
	}
	return u, nil
}

// GetUserYtd return a user's yesterday data with given name
func (db *OsuDB) GetUserYtd(username string) (*User, error) {
	const query = `
SELECT username, playcount_ytd, rank_ytd, pp_ytd, acc_ytd, total_play_ytd 
FROM users 
WHERE username = ? OR user_id = ?
`

	stmtOut, err := db.Conn.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("query %s : %v", query, err)
	}
	defer stmtOut.Close()

	u := &User{}
	err = stmtOut.QueryRow(username, username).Scan(
		&u.Username, &u.PcYtd, &u.RankYtd, &u.PpYtd, &u.AccYtd, &u.TotalPlayYtd,
	)
	if err != nil {
		return nil, fmt.Errorf("scan %s : %v", query, err)
	}
	return u, nil
}

// InsertNewUser insert user data into database
func (db *OsuDB) InsertNewUser(
	userID string, username string, pc string, rank string, pp string, acc string, total_play string,
) error {

	const query = `
INSERT INTO users (
	user_id, username, playcount, rank, pp, acc, total_play
) VALUES (
	?,?,?,?,?,?,?
)
`
	stmtIn, err := db.Conn.Prepare(query)
	defer stmtIn.Close()
	if err != nil {
		return fmt.Errorf("query %s : %v", query, err)
	}
	res, err := stmtIn.Exec(userID, username, pc, rank, pp, acc, total_play)
	if err != nil {
		return fmt.Errorf("insert %s: %v", query, err)
	}
	af, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if af < 1 {
		return errors.New("no row affected")
	}
	return nil
}

// UpdateUser update user data with given data
func (db *OsuDB) UpdateUser(
	username string, pc string, rank string, pp string, acc string, total_play string,
) error {
	const query = `
UPDATE users
SET playcount=?, rank=?, pp=?, acc=?, total_play=?
WHERE username=?
`
	stmtUp, err := db.Conn.Prepare(query)
	if err != nil {
		return fmt.Errorf("query %s: %v", query, err)
	}
	res, err := stmtUp.Exec(pc, rank, pp, acc, total_play, username)
	if err != nil {
		return fmt.Errorf("update %s: %v", query, err)
	}
	rows, err := res.RowsAffected()
	if rows < 1 {
		return errors.New("no row affected")
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
