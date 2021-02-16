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
		if !assertSameUser(t, "shigetora") {
			t.Errorf("cookiezi not found")
		}
	})
	t.Run("get avimitin", func(t *testing.T) {
		if !assertSameUser(t, "avimitin") {
			t.Errorf("avimitin not found")
		}
	})
	t.Run("try to get tuna", func(t *testing.T) {
		if assertSameUser(t, "flyingtuna") {
			t.Errorf("unexpected user flyingtuna")
		}
	})
}

func assertSameUser(t testing.TB, want string) bool {
	t.Helper()
	u, err := db.GetPlayer(want)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false
		}
		t.Fatal(err)
	}

	var got = u.Username
	return got == want
}
