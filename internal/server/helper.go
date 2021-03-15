package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/avimitin/osu-api-server/internal/api"
	"github.com/avimitin/osu-api-server/internal/database"
)

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

func assertIsGetMethod(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(NewJsonMsg().Set("error", "invalid method"))
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

type JsonMsg map[string]interface{}

func NewJsonMsg() JsonMsg {
	return make(JsonMsg)
}

func (jm JsonMsg) Set(key string, content interface{}) JsonMsg {
	jm[key] = content
	return jm
}

func fmtJsonMsg(key string, content interface{}) JsonMsg {
	return JsonMsg{key: content}
}

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func getFormValue(r *http.Request, key string) string {
	err := r.ParseForm()
	if err != nil {
		return ""
	}
	return r.PostForm.Get(key)
}

func isNullString(input string) bool {
	return input == ""
}
