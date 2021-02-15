package database

import (
	"fmt"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
)

var (
	db       *MySQLDataStore
	fixtures *testfixtures.Loader
)

func InitTestfixtures(dir string) error {
	var err error
	db, err = NewMySQLStore(
		"osu_test:osu_test@tcp(127.0.0.1:3306)/osu?charset=utf8mb4",
	)
	if err != nil {
		return fmt.Errorf("connect to database: %v", err)
	}
	testFiles := testfixtures.Files(dir)
	fixtures, err = testfixtures.New(
		testfixtures.Database(db.db),
		testfixtures.Dialect("mysql"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testFiles,
	)
	return err
}

func prepareTestDatabase(t testing.TB) {
	t.Helper()
	if err := fixtures.Load(); err != nil {
		t.Fatalf("loading fixtures: %v", err)
	}
}
