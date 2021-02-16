package database

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
)

var (
	db          *MySQLDataStore
	fixtures    *testfixtures.Loader
	projectPath = os.Getenv("osuapi_project_root")
)

func InitTestfixtures(dir string) error {
	var err error
	db, err = NewMySQLStore(
		"osu_test:osu_test@tcp(127.0.0.1:3306)/osu_test?charset=utf8mb4",
	)
	if err != nil {
		return fmt.Errorf("connect to database: %v", err)
	}
	err = db.CheckHealth()
	if err != nil {
		return fmt.Errorf("check health: %v", err)
	}
	fixtures, err = testfixtures.New(
		testfixtures.Database(db.db),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory(dir),
	)
	return err
}

func prepareTestDatabase(t testing.TB) {
	t.Helper()
	if err := fixtures.Load(); err != nil {
		t.Fatalf("loading fixtures: %v", err)
	}
}

func TestMain(m *testing.M) {
	var err error

	err = InitTestfixtures(projectPath + "/internal/database/fixtures/")
	if err != nil {
		fatalF("init test fixtures: %v", err)
	}

	os.Exit(m.Run())
}

func fatalF(context string, args ...interface{}) {
	fmt.Printf(context, args...)
	os.Exit(1)
}

func TestGetPlayer(t *testing.T) {
	prepareTestDatabase(t)

	t.Run("get cookiezi", func(t *testing.T) {
		if !assertSameUser(t, "recent", "shigetora") {
			t.Errorf("cookiezi data mismatch")
		}
	})
	t.Run("get avimitin", func(t *testing.T) {
		if !assertSameUser(t, "recent", "avimitin") {
			t.Errorf("avimitin data mismatch")
		}
	})
	t.Run("try to get tuna", func(t *testing.T) {
		if assertSameUser(t, "recent", "flyingtuna") {
			t.Errorf("unexpected user flyingtuna")
		}
	})
}

func TestGetPlayerOldData(t *testing.T) {
	prepareTestDatabase(t)
	t.Run("get whitecat", func(t *testing.T) {
		if !assertSameUser(t, "yesterday", "whitecat") {
			t.Errorf("whitecat data mismatch")
		}
	})
	t.Run("get avimitin", func(t *testing.T) {
		if !assertSameUser(t, "yesterday", "avimitin") {
			t.Errorf("avimitin data mismatch")
		}
	})
	t.Run("try to get tuna", func(t *testing.T) {
		if assertSameUser(t, "yesterday", "flyingtuna") {
			t.Errorf("unexpected user flyingtuna")
		}
	})
}

func assertSameUser(t testing.TB, date string, want string) bool {
	t.Helper()
	var u *User
	var err error
	switch date {
	case "recent":
		u, err = db.GetPlayer(want)
		if exist := assertUserError(t, err); !exist {
			return false
		}

		var got = u.Username
		if !assertSameStr(t, got, want) {
			return false
		}

		switch want {
		case "avimitin":
			got = u.Recent.PlayCount
			want = "114514"
			return assertSameStr(t, got, want)
		case "shigetora":
			got = u.Recent.PlayTime
			want = "2478401"
			return assertSameStr(t, got, want)
		}

	case "yesterday":
		u, err = db.GetPlayerOld(want)
		if exist := assertUserError(t, err); !exist {
			return false
		}

		var got = u.Username
		if !assertSameStr(t, got, want) {
			return false
		}

		switch want {
		case "avimitin":
			got = u.Yesterday.PlayCount
			want = "114500"
			return assertSameStr(t, got, want)
		case "whitecat":
			got = u.Yesterday.PlayTime
			want = "1520200"
			return assertSameStr(t, got, want)
		}
	}
	return true
}

func assertSameStr(t testing.TB, got string, want string) bool {
	t.Helper()
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
	return true
}

func assertUserError(t testing.TB, err error) bool {
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false
		}
		t.Fatal(err)
	}
	return true
}
