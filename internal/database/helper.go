package database

import (
	"errors"
	"fmt"
)

var NoRowAff = errors.New("no row affected")

func (db *OsuDB) modify(query string, value ...interface{}) error {
	stmtIn, err := db.Conn.Prepare(query)
	if err != nil {
		return fmt.Errorf("query %s : %v", query, err)
	}
	defer stmtIn.Close()
	res, err := stmtIn.Exec(value...)
	if err != nil {
		return fmt.Errorf("exec %s: %v", query, err)
	}
	af, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if af < 1 {
		return fmt.Errorf("%s:%w", query, NoRowAff)
	}
	return nil
}
