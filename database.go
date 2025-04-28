package dblite

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type FileBasedDB struct {
	Protocol string // i.e sqlite, duckdb
	FilePath string // i.e "api.db"
}

func (fd *FileBasedDB) IsEmpty() bool {
	return fd.FilePath == ""
}

type ServerBasedDB struct {
	Protocol        string // i.e postgresql
	UserName        string
	Password        string
	Host            string
	Port            int64
	DBName          string
	ProtocolOptions ProtocolOptions
}

func (sd *ServerBasedDB) IsEmpty() (empty bool) {
	empty = false
	if sd.Protocol == "" || sd.UserName == "" || sd.Password == "" || sd.Host == "" || sd.Port == 0 || sd.DBName == "" || sd.ProtocolOptions.IsEmpty() {
		empty = true
	}
	return empty
}

func (sd *ServerBasedDB) ConnectionString() string {
	if sd.ProtocolOptions.IsEmpty() {
		return fmt.Sprintf("%v://%v:%v@%v:%v/%v", sd.Protocol, sd.UserName, sd.Password, sd.Host, sd.Port, sd.DBName)
	}
	return fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=%v&connect_timeout=%v", sd.Protocol, sd.UserName, sd.Password, sd.Host, sd.Port, sd.DBName, sd.ProtocolOptions.SSLMode, sd.ProtocolOptions.ConnectionTimeout)
}

type ProtocolOptions struct {
	SSLMode           string
	ConnectionTimeout string
}

func (po *ProtocolOptions) IsEmpty() bool {
	if po.SSLMode == "" || po.ConnectionTimeout == "" {
		return true
	}
	return false
}

type DB struct {
	Protocol   string // i.e postgres, mysql, mongodb, etc..
	ConnString string
}

func NewDB(serverBasedDB ServerBasedDB, fileBasedDB FileBasedDB) *DB {
	var protocol, connstring string
	if !fileBasedDB.IsEmpty() {
		protocol = fileBasedDB.Protocol
		connstring = fileBasedDB.FilePath
	} else {
		protocol = serverBasedDB.Protocol
		connstring = serverBasedDB.ConnectionString()
	}

	return &DB{Protocol: protocol, ConnString: connstring}

}

type OmniDB struct {
	DB   DB
	Conn *sql.DB
}

func NewOmniDatabase(db DB) *OmniDB {
	var conn *sql.DB
	return &OmniDB{DB: db, Conn: conn}
}

func (omnidb *OmniDB) OpenConnection() *sql.DB {
	db, err := sql.Open(omnidb.DB.Protocol, omnidb.DB.ConnString)
	if err != nil {
		checkError(err)
	}
	return db
}

func (omnidb *OmniDB) Close() {
	if omnidb.Conn != nil {
		checkError(omnidb.Conn.Close())
	}
}

func (omnidb *OmniDB) Exec(query string, args ...any) (sql.Result, error) {
	return Exec(omnidb.Conn, query, args...)
}
func (omnidb *OmniDB) ExecMany(query string, records [][]any) (error, error) {
	return ExecMany(omnidb.Conn, query, records)
}
func (omnidb *OmniDB) Query(query string, args ...any) (*sql.Rows, error) {
	return omnidb.Conn.Query(query, args...)
}
