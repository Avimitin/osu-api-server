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
	return nil
}

func initUserTable(db *sql.DB) error {
	result, err := db.Exec(`
CREATE TABLE IF NOT EXIST users(
id INT AUTO_INCREMENT,
user_id BIGINT ,
username VARCHAR(255) ,
playcount BIGINT ,
pp_rank INT ,
pp_raw INT ,
acc DOUBLE ,
)
	`)
	return nil
}
