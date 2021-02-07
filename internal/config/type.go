package config

import "fmt"

type Configuration struct {
	Key   string         `json:"key"`
	DBSec DatabaseSecret `json:"dbs"`
}

type DatabaseSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

// EncodeDSN return links to osu database
func (dbs DatabaseSecret) EncodeDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/osu?charset=utf8mb4", dbs.Username, dbs.Password, dbs.Host)
}
