package dblite

import (
	"fmt"
	"github.com/franela/goblin"
	"testing"
	"time"
)

const sqlOModel = `
DROP TABLE IF EXISTS omodel;
CREATE TABLE IF NOT EXISTS omodel (
	id            		 INTEGER NOT NULL PRIMARY KEY,
	email         		 TEXT NOT NULL UNIQUE,
	name          		 TEXT DEFAULT '',
	address   			 TEXT DEFAULT '',
	active        		 INTEGER DEFAULT 1
);
`

type OModel struct {
	Id      int64  `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Active  int    `json:"active"`
}

func NewOModel(id int64) *OModel {
	return &OModel{Id: id}
}

func (omodel *OModel) New() *OModel {
	return NewOModel(-1)
}

func (omodel *OModel) Clone() *OModel {
	var o = *omodel
	return &o
}

func (omodel *OModel) TableName() string {
	return "omodel"
}

var dbSource *DatabaseSource

func initPostgresDB() {
	var dbType = "postgres"
	var uri = "user=postgres password=1234 dbname=postgres host=localhost port=5432 sslmode=disable"
	var err error
	dbSource, err = NewDatabaseSource(dbType, uri)
	fmt.Printf("dbSource: %+v\n", dbSource)
	checkError(err)
	_, err = dbSource.Exec(sqlOModel)
	checkError(err)
}

func deInitPostgresDB() {
	if dbSource != nil {
		dbSource.Close()
	}
}

func TestODB(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Test ODB connection", func() {
		g.It("", func() {
			g.Timeout(1 * time.Hour)
			initPostgresDB()
			defer deInitPostgresDB()

			om := NewOModel(1)
			om.Name = "omodel1"
			om.Email = "omodel@odb.com"
			om.Address = "2221 454 343"
			cols := []string{"id", "email", "name", "address"}
			on := On{}
			bln, _, err := Insert(dbSource.Conn, om, cols, on, "postgres")
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()
		})

	})
}
