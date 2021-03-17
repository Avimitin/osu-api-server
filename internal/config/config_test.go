package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func prepareFile(filepath string) (*Configuration, error) {
	var config = prepareConfigToTest()
	data, _ := json.Marshal(config)
	err := ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		return nil, fmt.Errorf("write to %s with %s: %v", filepath, data, err)
	}
	return config, nil
}

func prepareConfigToTest() *Configuration {
	return &Configuration{
		Key:    "hash",
		DBType: "redis",
		DatabaseSettings: DatabaseSettings{
			Host: "localhost",
		},
	}
}

func cleanFile(filepath string) error {
	return os.Remove(filepath)
}

func TestGetConfig(t *testing.T) {
	t.Run("get file from OSU_CONF_PATH", func(t *testing.T) {
		testPath := "."
		var testConfig *Configuration
		var err error
		testConfig, err = prepareFile(path.Join(testPath, "config.json"))
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err = cleanFile(path.Join(testPath, "config.json")); err != nil {
				t.Fatal(err)
			}
			os.Unsetenv("OSU_CONF_PATH")
		}()

		os.Setenv("OSU_CONF_PATH", testPath)
		got, err := GetConfig()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, testConfig) {
			t.Errorf("want %+v got %+v", testConfig, got)
		}
	})
	t.Run("get config from env", func(t *testing.T) {
		var testPath = "."
		var err error

		var config = prepareConfigToTest()
		os.Setenv("OSU_CONF_PATH", testPath)
		os.Setenv("OSU_API_KEY", config.Key)
		os.Setenv("OSU_DB_TYPE", config.DBType)
		os.Setenv("OSU_DB_HOST", config.DatabaseSettings.Host)

		defer func() {
			os.Unsetenv("OSU_CONF_PATH")
			os.Unsetenv("OSU_API_KEY")
			os.Unsetenv("OSU_DB_TYPE")
			os.Unsetenv("OSU_DB_HOST")
			os.Remove(path.Join(testPath, "config.json"))
		}()

		got, err := GetConfig()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, config) {
			t.Errorf("want %+v got %+v", config, got)
		}

		var tempF []byte
		tempF, err = ioutil.ReadFile(path.Join(testPath, "config.json"))
		if err != nil {
			t.Fatal(err)
		}

		if len(tempF) == 0 {
			t.Errorf("config.json is null")
		}
	})
}
