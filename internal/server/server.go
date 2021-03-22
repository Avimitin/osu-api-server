package server

import (
	"errors"
	"fmt"

	"github.com/avimitin/osu-api-server/internal/api"
	"github.com/avimitin/osu-api-server/internal/config"
	"github.com/avimitin/osu-api-server/internal/database"
)

var (
	serverErr    = errors.New("unexpected server error")
	nullInputErr = errors.New("invalid null input")
)

// PrepareServer initialized all the service
func PrepareServer() (*OsuServer, error) {
	// load user setting
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("server prepare config: %v", err)
	}
	// initialize query key
	err = api.KeyInit(cfg.Key)
	if err != nil {
		return nil, err
	}
	// initialize database connection
	dsn := cfg.DatabaseSettings.EncodeDSN(cfg.DBType)
	db, err := database.Connect(cfg.DBType, dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to %s at %s:%v", cfg.DBType, dsn, err)
	}
	// check database connection
	if err = db.CheckUserDataStoreHealth(); err != nil {
		return nil, fmt.Errorf("check %s health:%v", cfg.DBType, err)
	}
	return NewOsuServer(NewOsuPlayerData(db)), nil
}
