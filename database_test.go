package dblite

import (
	"testing"

	"github.com/franela/goblin"
)

var omnidb *OmniDB
var db *DB
var serverBasedDB ServerBasedDB = ServerBasedDB{
	Protocol: "postgres",
	UserName: "",
	Password: "",
	Host:     "",
	Port:     5000,
	DBName:   "",
}
var fileBasedDB FileBasedDB = FileBasedDB{}

func TestOmniDB(t *testing.T) {
	db = NewDB(serverBasedDB, fileBasedDB)
	g := goblin.Goblin(t)
	g.Describe("Test Database Connection", func() {})

}
