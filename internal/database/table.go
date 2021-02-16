package database

const (
	tableUsers = `
CREATE TABLE IF NOT EXISTS users (
	id       INT NOT NULL AUTO_INCREMENT,
	user_id  VARCHAR(18) NOT NULL UNIQUE,
	username VARCHAR(255) NOT NULL,
	PRIMARY  KEY (id)
)CHARSET=utf8mb4
`

	tableRecentData = `
 CREATE TABLE IF NOT EXISTS recent_data ( 
 id         INT NOT NULL,
 play_count VARCHAR(18) NOT NULL,
 rank       VARCHAR(18) NOT NULL,
 pp         VARCHAR(18) NOT NULL,
 acc        VARCHAR(18) NOT NULL,
 play_time  VARCHAR(18) NOT NULL,
 FOREIGN    KEY(id) REFERRENCES users(id)
 )CHARSET=utf8mb4
`

	tableYesterdayData = `
 CREATE TABLE IF NOT EXISTS yesterday_data ( 
 id         INT NOT NULL,
 play_count VARCHAR(18) NOT NULL,
 rank       VARCHAR(18) NOT NULL,
 pp         VARCHAR(18) NOT NULL,
 acc        VARCHAR(18) NOT NULL,
 play_time  VARCHAR(18) NOT NULL,
 FOREIGN    KEY(id) REFERRENCES users(id)
 )CHARSET=utf8mb4
`
)
