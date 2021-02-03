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
	err := json.Unmarshal(body, &beatmaps)
	if err != nil {
		// handle not RESTful response
		if strings.Contains(err.Error(), "invalid character") {
			return nil, fmt.Errorf("%s\n\nis not json format", body)
		}
		// handle error response
		if strings.Contains(err.Error(), "cannot unmarshal object") {
			var respErr APIResponseError
			err = json.Unmarshal(body, &respErr)
			// handle other error (seldom appear, may remove someday)
			if err != nil {
				return nil, fmt.Errorf("unknown body: %s", body)
			}
			return nil, fmt.Errorf(respErr.Error)
		}
		return nil, fmt.Errorf("unmarshal beatmaps: %v", err)
	}
	return beatmaps, nil
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
