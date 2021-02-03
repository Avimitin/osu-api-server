package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/avimitin/osuapi/internal/config"
)

const (
	APIURL = "https://osu.ppy.sh/api/"
)

var (
	key string
)

func init() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	key = conf.Key
}

// GetBeatMaps return list of beatmap by specific beatmap set ID
// or specific beatmap id. Should notice that if beatmap set and
// beatmap id is both given, beatmap id is foremore.
func GetBeatMaps(setID string, mapID string) ([]*Beatmap, error) {
	// if beatmap_id is specific, request it first
	if mapID != "" {
		body, err := request(
			buildURL("get_beatmaps",
				map[string]string{
					"k": key,
					"b": mapID,
				}),
		)
		if err != nil {
			return nil, err
		}
		return unmarshallBeatMaps(body)
	}

	if setID != "" {
		body, err := request(
			buildURL("get_beatmaps",
				map[string]string{
					"k": key,
					"b": mapID,
				}),
		)
		if err != nil {
			return nil, err
		}
		return unmarshallBeatMaps(body)
	}

	return nil, errors.New("invalid query parameters")
}

func unmarshallBeatMaps(body []byte) ([]*Beatmap, error) {
	var beatmaps []*Beatmap
	err := unmarshal(body, &beatmaps)
	if err != nil {
		return nil, err
	}
	return beatmaps, nil
}

// GetUsers return information about given users
func GetUsers(username string) ([]*User, error) {
	resp, err := request(buildURL("get_user", map[string]string{"k": key, "u": username}))
	if err != nil {
		return nil, err
	}
	return unmarshallUsers(resp)
}

func unmarshallUsers(body []byte) ([]*User, error) {
	var users []*User
	err := unmarshal(body, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func unmarshal(body []byte, object interface{}) error {
	err := json.Unmarshal(body, object)
	if err != nil {
		// handle not RESTful response
		if strings.Contains(err.Error(), "invalid character") {
			return fmt.Errorf("%s\n\nis not json format", body)
		}
		// handle error response
		if strings.Contains(err.Error(), "cannot unmarshal object") {
			var respErr APIResponseError
			err = json.Unmarshal(body, &respErr)
			// handle other error (seldom appear, may remove someday)
			if err != nil {
				return fmt.Errorf("unknown body: %s", body)
			}
			return fmt.Errorf(respErr.Error)
		}
		return fmt.Errorf("unmarshal beatmaps: %v", err)
	}
	return nil
}

func buildURL(method string, params map[string]string) string {
	if method == "" || params == nil {
		return ""
	}
	prefix := APIURL + method + "?"
	val := url.Values{}
	for k, v := range params {
		val.Set(k, v)
	}
	return prefix + val.Encode()
}

func request(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request %s: %v", url, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read %v: %v", resp.Body, err)
	}
	return body, nil
}
