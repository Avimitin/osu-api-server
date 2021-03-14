package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	pa "path"
)

// GetConfig search config at
// $HOME/.config/osuapi/config.json
// or $osu_conf_path/config.json.
// return error if no env set
func GetConfig() (*Configuration, error) {
	if path := os.Getenv("osu_conf_path"); path != "" {
		return getConfigFromPath(pa.Join(path, "config.json"))
	}

	if path := os.Getenv("HOME"); path != "" {
		return getConfigFromPath(pa.Join(path, ".config", "osuapi", "config.json"))
	}
	return nil, errors.New("no variable set")
}

func getConfigFromPath(path string) (*Configuration, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s:%v", path, err)
	}
	var cfg *Configuration
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("marshal %s: %v", file, err)
	}
	return cfg, nil
}
