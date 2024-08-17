package dblite

import (
	"database/sql"
	"fmt"
	"time"
)

func Exec(conn *sql.DB, query string, args ...any) (sql.Result, error) {
	return conn.Exec(query, args...)
}

func ExecMany(conn *sql.DB, query string, records [][]any) (error, error) {
	tx, err := conn.Begin()
	if err != nil {
		return err, nil
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		var prepError = tx.Rollback() // rollback if prepare fails
		return err, prepError
	}
	defer stmt.Close()

	for _, record := range records {
		_, err = stmt.Exec(record...)
		if err != nil {
			var errRollback error
			for i := 1; i <= 5; i++ {
				errRollback = tx.Rollback()
				if errRollback == nil {
					break
				} else {
					errRollback = fmt.Errorf("failed to rollback after %d attempts: %v", i, errRollback)
				}
				time.Sleep(time.Second * 2)
			}
			return err, errRollback
		}
	}

	return tx.Commit(), nil // commit all changes at once
}
