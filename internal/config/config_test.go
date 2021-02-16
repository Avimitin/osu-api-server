package config

import (
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	t.Run("get file from env", func(t *testing.T) {
		os.Setenv("osu_conf_path", "/home/avimitin")
		got, err := GetConfig()
		if err != nil {
			t.Fatal(err)
		}
		want := "hash"
		if got.Key != want {
			t.Errorf("got %s want %s", got, want)
		}
		_ = os.Unsetenv("osu_conf_path")
	})

		os.Unsetenv("osu_conf_path")
	})
}
