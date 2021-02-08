package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPlayer(t *testing.T) {
	t.Run("Get avimitin score", func(t *testing.T) {
		got := makeGetUserStatRequest("avimitin")
		want := `{"username": "avimitin"}`

		assertGetUser(t, got, want)
	})

	t.Run("Get coooool score", func(t *testing.T) {
		got := makeGetUserStatRequest("coooool")
		want := `{"username": "coooool"}`

		assertGetUser(t, got, want)
	})
}

func makeGetUserStatRequest(username string) string {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/players/%s", username), nil)
	response := httptest.NewRecorder()
	OsuServer(response, request)
	return response.Body.String()
}

func assertGetUser(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
