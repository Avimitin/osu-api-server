package config

type Configuration struct {
	Key   string         `json:"key"`
	DBSec DatabaseSecret `json:"dbs"`
}

type DatabaseSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}
