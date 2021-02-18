package server

import "testing"

func TestHasPrefix(t *testing.T) {
	t.Run("test exist", func(t *testing.T) {
		test := "/api/v1/players/avimitin"
		ok, got := hasPrefix(test)
		if !ok {
			t.Errorf("%s don't match", test)
		}
		want := "players"
		if got != want {
			t.Errorf("Got %s want %s", got, want)
		}
	})

	t.Run("test invalid", func(t *testing.T) {
		test := "/api/players/avimitin"
		ok, _ := hasPrefix(test)
		if ok {
			t.Errorf("unexpected match")
		}
	})
}
