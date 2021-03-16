package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// GetConfig search config at
// $HOME/.config/osuapi/config.json or
// $OSU_CONF_PATH/config.json.
// If $OSU_CONF_PATH is specific, and no file found at the
// given path, program will try to read config field at
// environment variable. return error if no env set.
func GetConfig() (*Configuration, error) {
	if confPath := os.Getenv("OSU_CONF_PATH"); confPath != "" {
		var cfg *Configuration
		var err error
		cfg, err = getConfigFromPath(path.Join(confPath, "config.json"))
		if err == nil {
			return cfg, nil
		}

		// if error is about file not found, get those fields from env
		switch {
		case errorHas(err, "no such file or directory"):
		case errorHas(err, "The system cannot find the file specified"):
			cfg, err = getConfigFromEnv()
			if err != nil {
				return nil, err
			}
			err = SaveConfig(confPath, cfg)
			if err != nil {
				return nil, fmt.Errorf("writing config: %v", err)
			}
			return cfg, nil
		}

		return nil, err
	}

	if confPath := os.Getenv("HOME"); confPath != "" {
		return getConfigFromPath(path.Join(confPath, ".config", "osuapi", "config.json"))
	}

	return nil, errors.New("no config found")
}

func errorHas(err error, content string) bool {
	return strings.Contains(err.Error(), content)
}

func nilAndError(content string) (*Configuration, error) {
	return nil, errors.New(content)
}

func getConfigFromEnv() (*Configuration, error) {
	var cfg = new(Configuration)
	if cfg.Key = os.Getenv("OSU_API_KEY"); cfg.Key == "" {
		return nilAndError("no key is given")
	}
	if cfg.DBType = os.Getenv("OSU_DB_TYPE"); cfg.DBType == "" {
		return nilAndError("no database type is given")
	}

	var dbs = DatabaseSettings{
		Username: os.Getenv("OSU_DB_USERNAME"),
		Password: os.Getenv("OSU_DB_PASSWORD"),
	}
	if dbs.Host = os.Getenv("OSU_DB_HOST"); dbs.Host == "" {
		return nilAndError("no database host is given")
	}

	cfg.DatabaseSettings = dbs

	return cfg, nil
}

func getConfigFromPath(path string) (*Configuration, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s:%v", path, err)
	}
	var cfg *Configuration
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("json decode %s: %v", file, err)
	}
	return cfg, nil
}

// SaveConfig store configuration to as config.json to
// given path
func SaveConfig(where string, cfg *Configuration) error {
	var data, err = json.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("json encode config: %v", err)
	}
	err = ioutil.WriteFile(path.Join(where, "config.json"), data, 0644)
	if err != nil {
		return fmt.Errorf("write config to %s: %v", where, err)
	}
	return nil
}
