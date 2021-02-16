package database

import "strings"

var (
	tableUsers = replace(`
CREATE TABLE IF NOT EXISTS users (
	"id"       INT NOT NULL AUTO_INCREMENT,
	"user_id"  TINYTEXT NOT NULL UNIQUE,
	"username" TINYTEXT NOT NULL,
	PRIMARY  KEY (id)
)CHARSET=utf8mb4
`)

	tableRecentData = replace(`
CREATE TABLE IF NOT EXISTS recent_data ( 
	 "id"         INT NOT NULL,
	 "play_count" TINYTEXT NOT NULL,
	 "rank"       TINYTEXT NOT NULL,
	 "pp"         TINYTEXT NOT NULL,
	 "acc"        TINYTEXT NOT NULL,
	 "play_time"  TINYTEXT NOT NULL,
	 FOREIGN    KEY(id) REFERENCES users(id)
 )CHARSET=utf8mb4
`)

	tableYesterdayData = replace(`
 CREATE TABLE IF NOT EXISTS yesterday_data ( 
	 "id"         INT NOT NULL,
	 "play_count" TINYTEXT NOT NULL,
	 "rank"       TINYTEXT NOT NULL,
	 "pp"         TINYTEXT NOT NULL,
	 "acc"        TINYTEXT NOT NULL,
	 "play_time"  TINYTEXT NOT NULL,
	 FOREIGN    KEY(id) REFERENCES users(id)
 )CHARSET=utf8mb4
`)
)

// replace replace all the double quote to backtick
func replace(s string) string {
	return strings.ReplaceAll(s, `"`, "`")
}
