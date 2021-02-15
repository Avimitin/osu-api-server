package database

import (
	"database/sql"

	"github.com/go-testfixtures/testfixtures/v3"
)

var (
	db       *sql.DB
	fixtures *testfixtures.Loader
)

func InitTestfixtures(dir string) error {
	var err error
	testFiles := testfixtures.Files(dir)
	fixtures, err = testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("mysql"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testFiles,
	)
	return err
}
