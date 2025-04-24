package dblite

import (
	"database/sql"
	"fmt"
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

type OmniDatabase struct {
	DB   DB
	Conn *sql.DB
}

func NewOmniDatabase(db DB) *OmniDatabase {
	var conn *sql.DB
	return &OmniDatabase{DB: db, Conn: conn}
}
