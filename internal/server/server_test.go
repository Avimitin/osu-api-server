package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/avimitin/osuapi/internal/api"
	"github.com/avimitin/osuapi/internal/database"
)

type playerDataTest struct {
	stat map[string]string
}

func (pdt *playerDataTest) GetPlayerStat(name string) (string, error) {
	username, ok := pdt.stat[name]
	if !ok {
		return "", errors.New("user not found")
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

	t.Run("monk actual curl", func(t *testing.T) {
		req := makeUserRequest("avimitin")
		resp := httptest.NewRecorder()
		ser := &OsuServer{
			Data: &OsuPlayerData{},
		}
		ser.ServeHTTP(resp, req)
		assertStatus(t, resp.Code, http.StatusOK)
		p := Player{}
		err := json.Unmarshal(resp.Body.Bytes(), &p)
		if err != nil {
			t.Fatalf("unmarshal %s : %v", resp.Body.Bytes(), err)
		}
		if p.Data == nil || p.Data.Accuracy == "" {
			t.Errorf("unexpected: %+v", p)
		}

		if p.Diff == nil || p.Diff.Acc == "" {
			t.Errorf("unexpected: %+v", p)
		}
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
	j, err := parseFloat("4000.50")
	if err != nil {
		t.Fatalf("cast 4000.5 got %v", err)
	}
	if i != 3426.48 {
		t.Errorf("want 3426.48 got %f", i)
	}
	if j != 4000.50 {
		t.Errorf("want 4000.5 got %f", j)
	}

	if i-j != -574.02 {
		t.Errorf("expect %f got %f", -574.02, i-j)
	}

	if fmt.Sprintf("%.3f", i-j) != "-574.020" {
		t.Errorf("expect %s got %s", "-574.020", fmt.Sprintf("%.3f", i-j))
	}
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
