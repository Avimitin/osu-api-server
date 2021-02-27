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

func TestAddPlayer(t *testing.T) {
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
		t.Errorf("add new user: %v", err)
	}

	got, err := rdb.GetPlayer("test:recent")

	if err != nil {
		t.Errorf("get user: %v", err)
	}

	if !reflect.DeepEqual(&want, got) {
		t.Errorf("got %+v \nwant %+v", got, want)
	}

}
