package dblite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var Instance *Database

func Init(fileName string) {
	var err error
	Instance, err = NewDatabase(fileName)
	CheckError(err)
}

func DeInit() {
	if Instance != nil {
		Instance.Close()
	}
}

type Database struct {
	file string
	Conn *sql.DB
}

func NewDatabase(dbpath string) (*Database, error) {
	conn, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	return &Database{file: dbpath, Conn: conn}, nil
}

func (db *Database) Close() {
	if db.Conn != nil {
		CheckError(db.Conn.Close())
	}
}

func (db *Database) Exec(query string, args ...any) (sql.Result, error) {
	return db.Conn.Exec(query, args...)
}

func (db *Database) ExecMany(query string, records [][]any) (error, error) {
	tx, err := db.Conn.Begin()
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

func (db *Database) Query(query string, args ...any) (*sql.Rows, error) {
	return db.Conn.Query(query, args...)
}
