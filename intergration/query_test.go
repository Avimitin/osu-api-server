package intergration

import (
	"fmt"
	"os"
	"testing"

	"github.com/avimitin/osu-api-server/internal/api"
)

func TestMain(t *testing.M) {
	err := api.KeyInit(os.Getenv("api_key"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(t.Run())
}

func TestGetBeatMaps(t *testing.T) {
	t.Run("get beatmap set", func(t *testing.T) {
		maps, err := api.GetBeatMaps("983911", "")
		if err != nil {
			t.Fatal(err)
		}
		for _, bmap := range maps {
			got := bmap.BeatmapsetID
			want := "983911"
			if got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		}
	})

	t.Run("get beatmap set", func(t *testing.T) {
		maps, err := api.GetBeatMaps("", "2118444")
		if err != nil {
			t.Fatal(err)
		}
		for _, bmap := range maps {
			got := bmap.BeatmapID
			want := "2118444"
			if got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		}
	})
}

func TestGetUsers(t *testing.T) {
	test := func(want string) {
		t.Run("get user", func(t *testing.T) {
			users, err := api.GetUsers(want)
			if err != nil {
				t.Fatal(err)
			}
			if users == nil {
				t.Errorf("failed to fetch users")
			}

			for _, user := range users {
				if user.Username != want {
					t.Errorf("got %s want %s", user.Username, want)
				}
			}
		})
	}

	test("avimitin")
}

func TestGetUserBest(t *testing.T) {
	userid := "16900842"
	mode := ""
	limit := 10
	maps, err := api.GetUserBest(userid, mode, limit)
	if err != nil {
		t.Fatal(err)
	}
	for _, bmap := range maps {
		if bmap.UserID != userid {
			t.Errorf("want %s score got %s", userid, bmap.UserID)
		}
	}
}

func TestGetUserRecent(t *testing.T) {
	username := ""
	recentMaps, err := api.GetUserRecent(username, 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, recentMap := range recentMaps {
		if recentMap == nil {
			t.Fatalf("fetch nil map")
		}
		if recentMap.Score == "" {
			t.Errorf("fetch no score")
		}
	}
}
