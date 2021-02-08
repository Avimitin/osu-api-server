package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPlayer(t *testing.T) {
	t.Run("Get avimitin score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/v1/players/avimitin", nil)
		response := httptest.NewRecorder()

		OsuServer(response, request)

		got := response.Body.String()
		want := `{"username": "avimitin"}`

		if got != want {
			t.Errorf("want %s got %s", want, got)
		}
	})

	t.Run("Get coooool score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/v1/players/coooool", nil)
		response := httptest.NewRecorder()

		OsuServer(response, request)

		got := response.Body.String()
		want := `{"username": "coooool"}`

		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}
