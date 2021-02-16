package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

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

	var want []string = []string{"users", "recent_data", "yesterday_data"}
	wantMatch := len(want)
	var exist = false

	if len(tables) == 0 {
	} else {
		var match int
		for _, t := range tables {
			for _, w := range want {
				if t == w {
					match++
					break
				}
			}
		}

		if wantMatch == match {
			exist = true
		}
	}

	if !exist {
		log.Println("database is setting up table...")
		return MySQLInitTable(mds.db)
	}
	return nil
}

func MySQLInitTable(db *sql.DB) error {
	_, err := db.Exec(tableUsers)
	if err != nil {
		return fmt.Errorf("init table: %s:%v", tableUsers, err)
	}
	_, err = db.Exec(tableRecentData)
	if err != nil {
		return fmt.Errorf("init table: %s:%v", tableRecentData, err)
	}
	_, err = db.Exec(tableYesterdayData)
	if err != nil {
		return fmt.Errorf("init table: %s:%v", tableYesterdayData, err)
	}
	return nil
}

func (mds *MySQLDataStore) GetPlayer(username string) (*User, error) {
	u := &User{}
	query := queryUserRecentData
	stmtOut, err := mds.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare %s: %v", query, err)
	}
	defer stmtOut.Close()
	err = stmtOut.QueryRow(username, username).Scan(
		&u.Username, &u.Recent.PlayCount, &u.Recent.Rank, &u.Recent.PP, &u.Recent.Acc, &u.Recent.PlayTime)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, fmt.Errorf("user %s not found", username)
		}
		return nil, fmt.Errorf("query %s with param %s: %v", query, username, err)
	}
	return u, nil
}

func (mds *MySQLDataStore) GetPlayerOld(username string) (*User, error) {
	query := queryUserOldData
	stmtOut, err := mds.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("query %s : %v", query, err)
	}
	defer stmtOut.Close()

	u := &User{}
	err = stmtOut.QueryRow(username, username).Scan(
		&u.Username,
		&u.Yesterday.PlayCount, &u.Yesterday.Rank, &u.Yesterday.PP, &u.Yesterday.Acc, &u.Yesterday.PlayTime,
	)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, fmt.Errorf("user %s not found", username)
		}
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
	err := mds.modify(
		query, u.UserID, u.Username,
		u.Recent.PlayCount, u.Recent.Rank, u.Recent.PP, u.Recent.Acc, u.Recent.PlayTime)
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
	err := mds.modify(query,
		u.Recent.PlayCount, u.Recent.Rank, u.Recent.PP, u.Recent.Acc, u.Recent.PlayTime, u.Username)
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
	err := mds.modify(
		query,
		u.Yesterday.PlayCount, u.Yesterday.Rank, u.Yesterday.PP, u.Yesterday.Acc, u.Yesterday.PlayTime,
		u.Username)
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
