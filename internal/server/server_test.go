package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/avimitin/osu-api-server/internal/api"
	"github.com/avimitin/osu-api-server/internal/database"
)

type playerDataTest struct {
	stat map[string]string
}

func (pdt *playerDataTest) GetPlayerStat(name string) (*Player, error) {
	if name == "" {
		return nil, errors.New("null user input")
	}
	username, ok := pdt.stat[name]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &Player{Data: &api.User{Username: username}}, nil
}

func (pdt *playerDataTest) GetRecent(name, mapID string, perf bool) (*api.RecentPlay, error) {
	if name == "" {
		return nil, errors.New("null input")
	}
	username, ok := pdt.stat[name]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &api.RecentPlay{
		UserID: username,
	}, nil
}

func (pdt *playerDataTest) GetBeatMaps(setID, mapID string) (*api.Beatmap, error) {
	return &api.Beatmap{BeatmapID: mapID}, nil
}

func TestGetPlayer(t *testing.T) {
	t.Run("Get avimitin score", func(t *testing.T) {
		response := httptest.NewRecorder()
		got := makeGetUserStatRequest("avimitin", response)
		want := "avimitin"

		assertSameUser(t, got, want)
		assertStatus(t, response.Code, http.StatusOK)
		assertJsonHeader(t, response)
	})

	t.Run("Get error", func(t *testing.T) {
		req := makeUserRequest("jixun")
		response := httptest.NewRecorder()
		ser := newSer()
		ser.ServeHTTP(response, req)
		assertErrMsg(t, response.Body.String(), serverErr.Error())
		assertStatus(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("nil input", func(t *testing.T) {
		req := makeUserRequest("")
		response := httptest.NewRecorder()
		ser := newSer()
		ser.ServeHTTP(response, req)
		assertErrMsg(t, response.Body.String(), nullInputErr.Error())
		assertStatus(t, response.Code, http.StatusInternalServerError)
	})
}

func TestGetDiff(t *testing.T) {
	t.Run("get positive diff", func(t *testing.T) {
		diff, err := getUserDiff(
			&api.User{
				PpRaw:              "3426.48",
				Accuracy:           "97.31963348388672",
				Playcount:          "14085",
				TotalSecondsPlayed: "1041175",
				PpRank:             "111254",
			}, "recent",
			&database.Date{
				Recent: database.Data{
					PP:        "4000.50",
					Acc:       "97.89052605628967",
					PlayTime:  "2478401",
					PlayCount: "25000",
					Rank:      "57",
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		if diff.Rank != "-111197" {
			t.Errorf("get %+v", diff)
		}
	})
}

func TestCastFloat64(t *testing.T) {
	i, err := parseFloat("3426.48")
	if err != nil {
		t.Fatalf("cast 3426.48 got %v", err)
	}
	if i != 3426.48 {
		t.Errorf("want 3426.48 got %f", i)
	}
}

func TestPanicNewOsuServer(t *testing.T) {
	defer func() {
		err := recover()
		if err != "nil data store" {
			t.Errorf("recover a panic failed")
		}
	}()

	NewOsuServer(nil)
}

func TestGetUserRecent(t *testing.T) {
	response := httptest.NewRecorder()
	form := url.Values{"player": {"avimitin"}}
	request, _ := http.NewRequest(
		http.MethodPost, "http://localhost:11451/api/v1/recent", strings.NewReader(form.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ser := newSer()
	ser.ServeHTTP(response, request)
	var score *api.RecentPlay
	err := json.Unmarshal(response.Body.Bytes(), &score)
	if err != nil {
		t.Fatalf("Unmarshal %s:%v", response.Body.String(), err)
	}
	got := score.UserID
	want := "avimitin"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestGetBeatMaps(t *testing.T) {
	response := httptest.NewRecorder()
	form := url.Values{"map_id": {"126645"}}
	request, _ := http.NewRequest(
		http.MethodPost,
		"/api/v1/beatmap",
		strings.NewReader(form.Encode()),
	)
	ser := newSer()
	ser.ServeHTTP(response, request)
	assertStatus(t, response.Code, http.StatusOK)
	assertJsonHeader(t, response)
	var bmaps *api.Beatmap
	err := json.Unmarshal(response.Body.Bytes(), &bmaps)
	if err != nil {
		t.Fatal(err)
	}
	got := bmaps.BeatmapID
	want := "126645"
	if got != want {
		t.Errorf("got %s want %s", bmaps.BeatmapID, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func makeUserRequest(username string) (request *http.Request) {
	form := url.Values{"player": {username}}
	request, _ = http.NewRequest(
		http.MethodPost,
		"http://example.com/api/v1/player",
		strings.NewReader(form.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
	return NewOsuServer(pdt)
}

func assertSameUser(t testing.TB, got, want string) {
	t.Helper()
	var u *Player
	err := json.Unmarshal([]byte(got), &u)
	if err != nil {
		t.Fatal(err)
	}
	if u == nil || u.Data == nil {
		t.Errorf("got %s want %s", got, want)
		return
	}
	if u.Data.Username != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func assertErrMsg(t testing.TB, got, want string) {
	t.Helper()

	var err JsonMsg
	e := json.Unmarshal([]byte(got), &err)
	if e != nil {
		t.Fatal(e)
	}

	if msg, ok := err["error"]; ok {
		if msg != want {
			t.Errorf("got %s want %s", got, want)
		}
		return
	}

	t.Errorf("error not exist")
}

func assertJsonHeader(t testing.TB, resp http.ResponseWriter) {
	t.Helper()
	var ct string
	if ct = resp.Header().Get("Content-Type"); ct != "" {
		want := "application/json"
		if ct != want {
			t.Errorf("got %s want %s", ct, want)
		}
		return
	}
	t.Errorf("no content-type,get %s header:%+v", ct, resp.Header())
}
