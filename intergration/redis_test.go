package intergration

import (
	"reflect"
	"testing"

	"github.com/avimitin/osu-api-server/internal/database"
)

var (
	rdb *database.RedisDataStore
)

func prepareRedisDB() error {
	rdb = database.NewRedisDataStore()
	if err := rdb.CheckHealth(); err != nil {
		return err
	}
	return nil
}

func init() {
	err := prepareRedisDB()
	if err != nil {
		fatalF("init redis connection: %v", err)
	}
}

func TestRedis(t *testing.T) {
	want := database.User{
		UserID:   "1",
		Username: "test",
		Date: database.Date{
			Recent: database.Data{
				PlayCount: "10",
				Rank:      "11",
				PP:        "12",
				Acc:       "13.33",
				PlayTime:  "1234567",
			},
		},
	}

	var err error
	err = rdb.AddPlayer(want)
	if err != nil {
		t.Fatalf("add new user: %v", err)
	}

	got, err := rdb.GetPlayer("test")

	if err != nil {
		t.Fatalf("get user: %v", err)
	}

	if !reflect.DeepEqual(&want, got) {
		t.Errorf("got %+v \nwant %+v", got, want)
	}

	want.Date.Recent.PlayCount = "20"
	want.Date.Recent.Rank = "21"
	err = rdb.Update(want)
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	got, err = rdb.GetPlayer("test")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}

	if !reflect.DeepEqual(&want, got) {
		t.Errorf("updated: got %+v \nwant %+v", got, want)
	}
}
