package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type playerDataTest struct {
	stat map[string]string
}

func (pdt *playerDataTest) GetPlayerStat(name string) (string, error) {
	username, ok := pdt.stat[name]
	if !ok {
		return "", nil
	}
	return fmt.Sprintf(`{"username": "%s"}`, username), nil
}

func TestGetPlayer(t *testing.T) {
	t.Run("Get avimitin score", func(t *testing.T) {
		response := httptest.NewRecorder()
		got := makeGetUserStatRequest("avimitin", response)
		want := `{"username": "avimitin"}`

		assertGetUser(t, got, want)
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("Get coooool score", func(t *testing.T) {
		response := httptest.NewRecorder()
		got := makeGetUserStatRequest("coooool", response)
		want := `{"username": "coooool"}`

		assertGetUser(t, got, want)
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("Get 404", func(t *testing.T) {
		req := makeUserRequest("jixun")
		response := httptest.NewRecorder()
		ser := newSer()
		ser.ServeHTTP(response, req)
		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func makeUserRequest(username string) (request *http.Request) {
	request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/players/%s", username), nil)
	return request
}

func makeGetUserStatRequest(username string, response *httptest.ResponseRecorder) string {
	request := makeUserRequest(username)
	ser := newSer()
	ser.ServeHTTP(response, request)
	return response.Body.String()
}

func newSer() *OsuServer {
	pdt := &playerDataTest{
		map[string]string{
			"coooool":  "coooool",
			"avimitin": "avimitin",
		},
	}
	return &OsuServer{
		Data: pdt,
	}
}

func assertGetUser(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
