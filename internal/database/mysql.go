package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLDataStore struct {
	db *sql.DB
}

func NewMySQLStore(dsn string) (mds *MySQLDataStore, err error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &MySQLDataStore{db}, nil
}

func (mds *MySQLDataStore) CheckHealth() error {
	rows, err := mds.db.Query("SHOW TABLES;")
	if err != nil {
		return fmt.Errorf("show tables : %v", err)
	}
	var tables []string
	for rows.Next() {
		var table string
		rowErr := rows.Scan(&table)
		if rowErr != nil {
			return fmt.Errorf("scan %v:%v", rows, rowErr)
		}
		tables = append(tables, table)
	}

	var want = "users"
	exist := false
	if len(tables) == 0 {
		return errors.New("no table")
	} else {
		for _, t := range tables {
			if t == want {
				exist = true
			}
		}
	}
	if !exist {
		return MySQLInitTable(mds.db)
	}
	return nil
}

func MySQLInitTable(db *sql.DB) error {
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

func (mds *MySQLDataStore) GetPlayer(username string) (*User, error) {
	const query = `
SELECT
	username, playcount, rank, pp, acc, total_play 
FROM
	users
WHERE
	username = ?
OR
	user_id = ?
`
	u := &User{}
	stmtOut, err := mds.db.Prepare(query)
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

func (mds *MySQLDataStore) GetPlayerOld(username string) (*User, error) {
	const query = `
SELECT 
	username, playcount_ytd, rank_ytd, pp_ytd, acc_ytd, total_play_ytd 
FROM 
	users 
WHERE 
	username = ? OR user_id = ?
`

	stmtOut, err := mds.db.Prepare(query)
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

func (mds *MySQLDataStore) AddPlayer(u User) error {
	const query = `
INSERT INTO users (
	user_id, username, playcount, rank, pp, acc, total_play
) VALUES (
	?,?,?,?,?,?,?
)
`
	err := mds.modify(query, u.UserID, u.Username, u.PlayCount, u.Rank, u.PP, u.Acc, u.TotalPlay)
	if err != nil {
		return err
	}
	return nil
}

func (mds *MySQLDataStore) Update(u User) error {
	const query = `
UPDATE 
	users
SET 
	playcount=?, rank=?, pp=?, acc=?, total_play=?
WHERE 
	username=?
`
	err := mds.modify(query, u.PlayCount, u.Rank, u.PP, u.Acc, u.TotalPlay, u.Username)
	if err != nil {
		return err
	}
	return nil
}

func (mds *MySQLDataStore) UpdateOld(u User) error {
	const query = `
UPDATE
	users
SET
	playcount_ytd=?, rank_ytd=?, pp_ytd=?, acc_ytd=?, total_play_ytd=?
WHERE
	username=?
`
	err := mds.modify(query, u.PcYtd, u.RankYtd, u.PpYtd, u.AccYtd, u.TotalPlayYtd, u.Username)
	if err != nil {
		return err
	}
	return nil
}

var NoRowAff = errors.New("no row affected")

func (mds *MySQLDataStore) modify(query string, value ...interface{}) error {
	stmtIn, err := mds.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("query %s : %v", query, err)
	}
	defer stmtIn.Close()
	res, err := stmtIn.Exec(value...)
	if err != nil {
		return fmt.Errorf("exec %s: %v", query, err)
	}
	af, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if af < 1 {
		return fmt.Errorf("%s:%w", query, NoRowAff)
	}
	return nil
}
