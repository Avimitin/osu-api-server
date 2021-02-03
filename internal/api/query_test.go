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
