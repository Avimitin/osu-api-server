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

func initUserTable(db *sql.DB) error {
	const userTable = `
CREATE TABLE IF NOT EXIST users(
id INT AUTO_INCREMENT,
user_id BIGINT,
username VARCHAR(255),
playcount BIGINT,
rank INT,
pp INT,
acc DOUBLE,
playcount_ytd BIGINT,
rank_ytd INT,
pp_ytd INT,
acc_ytd DOUBLE,
total_play_ytd BIGINT,
PRIMARY KEY (id)
)CHARSET=utf8mb4
	`
	_, err := db.Exec(userTable)
	if err != nil {
		return fmt.Errorf("init table: %s:%v", userTable, err)
	}

	return nil
}
