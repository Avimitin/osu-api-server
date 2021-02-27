package intergration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/avimitin/osu-api-server/internal/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
)

var (
	db          *database.MySQLDataStore
	fixtures    *testfixtures.Loader
	projectPath = os.Getenv("osuapi_project_root")
)

func InitTestfixtures(dir string) error {
	var err error
	var dsn = "osu_test:osu_test@tcp(%s)/osu_test?charset=utf8mb4"
	if host := os.Getenv("database_host"); host != "" {
		dsn = fmt.Sprintf(dsn, host)
	} else {
		dsn = fmt.Sprintf(dsn, "localhost:3306")
	}

	db, err = database.NewMySQLStore(dsn)
	if err != nil {
		return fmt.Errorf("connect to database: %v", err)
	}
	err = db.CheckHealth()
	if err != nil {
		return fmt.Errorf("check health: %v", err)
	}
	test_db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("prepare test database: %v", err)
	}
	fixtures, err = testfixtures.New(
		testfixtures.Database(test_db),
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

func init() {
	err := InitTestfixtures(filepath.Join(projectPath, "intergration/fixtures"))
	if err != nil {
		fatalF("init test fixtures: %v", err)
	}
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
	var u *database.User
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
