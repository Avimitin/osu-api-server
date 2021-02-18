package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/avimitin/osu-api-server/internal/api"
	"github.com/avimitin/osu-api-server/internal/config"
	"github.com/avimitin/osu-api-server/internal/database"
)

var (
	db *database.OsuDB
)

// PrepareServer initialized all the service
func PrepareServer() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("server prepare config: %v", err)
	}
	db, err = database.Connect("mysql", cfg.DBSec.EncodeDSN())
	if err != nil {
		return fmt.Errorf("connect to %s:%v", cfg.DBSec.EncodeDSN(), err)
	}
	if err = db.CheckUserDataStoreHealth(); err != nil {
		return fmt.Errorf("check database health:%v", err)
	}
	return nil
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
		diff, e = getUserDiff(user, "recent", &lu.Date)
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
