package config

import "fmt"

type Configuration struct {
	Port             int              `json:"port"`
	Key              string           `json:"key"`
	DBType           string           `json:"db_type"`
	DatabaseSettings DatabaseSettings `json:"database_settings"`
}

type DatabaseSettings struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

// EncodeDSN return links to osu database
func (dbs DatabaseSettings) EncodeDSN(dbType string) string {
	switch dbType {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s)/osu?charset=utf8mb4", dbs.Username, dbs.Password, dbs.Host)
	case "redis":
		return fmt.Sprintf("redis://%s:%s@%s", dbs.Username, dbs.Password, dbs.Host)
	default:
		return ""
	}
}
