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

	queryUser = `
SELECT
	users.username,
	recent_data.play_count,
	recent_data.rank,
	recent_data.pp,
	recent_data.acc,
	recent_data.play_time
FROM
	users, recent_data
WHERE
	users.username=?
AND
	recent_data.id=(
	SELECT
		id
	FROM
		users
	WHERE
		username=?
)
`
)

// replace replace all the double quote to backtick
func replace(s string) string {
	return strings.ReplaceAll(s, `"`, "`")
}
