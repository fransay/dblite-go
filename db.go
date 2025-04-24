package dblite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

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
		checkError(db.Conn.Close())
	}
}

func (db *Database) Exec(query string, args ...any) (sql.Result, error) {
	return Exec(db.Conn, query, args...)
}

func (db *Database) ExecMany(query string, records [][]any) (error, error) {
	return ExecMany(db.Conn, query, records)
}

func (db *Database) Query(query string, args ...any) (*sql.Rows, error) {
	return db.Conn.Query(query, args...)
}
