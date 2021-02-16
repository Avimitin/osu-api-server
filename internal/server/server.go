package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/avimitin/osuapi/internal/api"
	"github.com/avimitin/osuapi/internal/config"
	"github.com/avimitin/osuapi/internal/database"
)

var (
	db *database.OsuDB
)

func init() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("config initialized")
	db, err = database.Connect("mysql", cfg.DBSec.EncodeDSN())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("database connected")
	if err = db.CheckUserDataStoreHealth(); err != nil {
		log.Fatalf("check database health:%v", err)
	}
}

type PlayerData interface {
	GetPlayerStat(name string) (string, error)
}

type OsuServer struct {
	Data PlayerData
}

func (osuSer *OsuServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s sent %s request to %s", r.RemoteAddr, r.Method, r.URL.Host+r.URL.Path)
	if r.Method != http.MethodGet {
		fmt.Fprint(w, `{"error":"invalid method"}`)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/api/v1/players") {
		fmt.Fprintf(w, `{"error":"page %s not found"}`, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	player := strings.TrimPrefix(r.URL.Path, "/api/v1/players/")
	stat, err := osuSer.Data.GetPlayerStat(player)
	if err != nil {
		log.Printf("get %s data: %v", player, err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
		return
	}
	fmt.Fprint(w, stat)
}

type OsuPlayerData struct{}

type Player struct {
	Data *api.User  `json:"latest_data"`
	Diff *Different `json:"diff"`
}

type Different struct {
	PlayCount string `json:"play_count"`
	Rank      string `json:"rank"`
	PP        string `json:"pp"`
	Acc       string `json:"acc"`
	TotalPlay string `json:"total_play"`
}

func (opd *OsuPlayerData) GetPlayerStat(name string) (string, error) {
	p, e := getPlayerDataByName(name)
	if e != nil {
		return "", e
	}
	data, e := json.Marshal(p)
	if e != nil {
		return "", e
	}
	return string(data), nil
}

func getPlayerDataByName(name string) (*Player, error) {
	u, e := api.GetUsers(name)
	if e != nil {
		return nil, e
	}
	if len(u) <= 0 {
		return nil, errors.New("user " + name + " not found")
	}
	user := u[0]
	lu, e := db.GetUserRecent(user.Username)
	if e != nil {
		if strings.Contains(e.Error(), "user") {
			e = db.InsertNewUser(
				user.UserID, user.Username, user.Playcount, user.PpRank,
				user.PpRaw, user.Accuracy, user.TotalSecondsPlayed,
			)
			if e != nil {
				return nil, fmt.Errorf("insert user %s : %v", name, e)
			}
			log.Printf("inserted %s into database", user.Username)
		} else {
			return nil, fmt.Errorf("query user %s: %v", name, e)
		}
	}
	var diff *Different
	if lu != nil {
		diff, e = getUserDiff(user, lu)
	}
	if e != nil {
		return nil, e
	}

	p := &Player{
		Data: u[0],
		Diff: diff,
	}
	return p, nil
}

// getUserDiff return data with given date and current data different
func getUserDiff(current *api.User, with string, local *database.Date) (*Different, error) {
	var data database.Data
	switch with {
	case "recent":
		data = local.Recent
	case "yesterday":
		data = local.Yesterday
	default:
		return nil, errors.New("invalid data date")
	}

	// cast latest data
	pc, err := parseInt64(current.Playcount)
	if err != nil {
		return nil, fmt.Errorf("cast playcount: %v", err)
	}
	totalp, err := parseInt64(current.TotalSecondsPlayed)
	if err != nil {
		return nil, fmt.Errorf("cast total_play: %v", err)
	}
	acc, err := parseFloat(current.Accuracy)
	if err != nil {
		return nil, fmt.Errorf("cast acc: %v", err)
	}
	rank, err := atoi(current.PpRank)
	if err != nil {
		return nil, fmt.Errorf("cast rank: %v", err)
	}
	pp, err := parseFloat(current.PpRaw)
	if err != nil {
		return nil, fmt.Errorf("cast pp: %v", err)
	}

	// cast local data
	pcLocal, err := parseInt64(data.PlayCount)
	if err != nil {
		return nil, fmt.Errorf("cast playcount: %v", err)
	}
	totalpLocal, err := parseInt64(data.PlayTime)
	if err != nil {
		return nil, fmt.Errorf("cast total_play: %v", err)
	}
	accLocal, err := parseFloat(data.Acc)
	if err != nil {
		return nil, fmt.Errorf("cast acc: %v", err)
	}
	rankLocal, err := atoi(data.Rank)
	if err != nil {
		return nil, fmt.Errorf("cast rank: %v", err)
	}
	ppLocal, err := parseFloat(data.PP)
	if err != nil {
		return nil, fmt.Errorf("cast pp: %v", err)
	}

	// get data different
	pcDiff := strconv.FormatInt(pc-pcLocal, 10)
	totalpDiff := strconv.FormatInt(totalp-totalpLocal, 10)
	accDiff := fmt.Sprintf("%.2f%%", acc-accLocal)
	rankDiff := strconv.Itoa(rankLocal - rank)
	ppDiff := fmt.Sprintf("%.3f", pp-ppLocal)
	return &Different{
		PlayCount: pcDiff,
		TotalPlay: totalpDiff,
		Acc:       accDiff,
		Rank:      rankDiff,
		PP:        ppDiff,
	}, nil
}

func atoi(s string) (int, error) {
	i, e := strconv.Atoi(s)
	if e != nil {
		return -1, fmt.Errorf("cast %s to int: %v", s, e)
	}
	return i, nil
}

func parseInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("cast %s to int64: %v", s, err)
	}
	return i, nil
}

func parseFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return -1.0, fmt.Errorf("cast %s to float32: %v", s, err)
	}
	return f, nil
}
