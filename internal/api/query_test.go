package api

import "testing"

func TestGetBeatMaps(t *testing.T) {
	t.Run("get beatmap set", func(t *testing.T) {
		maps, err := GetBeatMaps("983911", "")
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
		maps, err := GetBeatMaps("", "2118444")
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
			users, err := GetUsers(want)
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
