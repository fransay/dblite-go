package dblite

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/franela/goblin"
)

const UserSQLModel = `
DROP TABLE IF EXISTS model;
CREATE TABLE IF NOT EXISTS model (
	id            		 INTEGER NOT NULL PRIMARY KEY,
	email         		 TEXT NOT NULL UNIQUE,
	name          		 TEXT DEFAULT '',
	address   			 TEXT DEFAULT '',
	active        		 INTEGER DEFAULT 1
);
`

var dbInstance *Database

type Model struct {
	Id      int64  `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Active  int    `json:"active"`
}

func NewModel(id int64) *Model {
	return &Model{Id: id}
}

func (model *Model) New() *Model {
	return NewModel(-1)
}

func (model *Model) Clone() *Model {
	var o = *model
	return &o
}

func (model *Model) TableName() string {
	return "model"
}

func (model *Model) InsertWithArgs() (bool, error) {
	return Insert(dbInstance.Conn, model, []string{
		`id`, `email`, `name`, `address`,
	}, On{
		On:        "CONFLICT(id) DO UPDATE SET email=?, name=?, address=?",
		Arguments: []any{model.Email, model.Name, model.Address},
	})
}

func (model *Model) InsertOnConflictDoNothing() (bool, error) {
	return Insert(dbInstance.Conn, model, []string{
		`id`, `email`, `name`, `address`,
	}, On{
		On: "CONFLICT(id) DO NOTHING",
	})
}

func (model *Model) Upsert() (bool, error) {
	var columns = []string{`id`, `email`, `name`, `address`}
	return Insert(dbInstance.Conn, model, columns, On{
		On:            "CONFLICT(id)",
		UpsertColumns: columns[1:],
	})
}

func (model *Model) Update() (bool, error) {
	var updateColumns = []string{`email`, `name`, `address`}
	return Update(dbInstance.Conn, model, updateColumns, WhereClause{
		Where:     ``,
		Arguments: []any{},
	})
}
func (model *Model) Delete() (int64, error) {
	return Delete(dbInstance.Conn, model, WhereClause{
		Where:     ``,
		Arguments: []any{},
	})
}

func (model *Model) Count() (int64, error) {
	var refCol = ``
	return Count(dbInstance.Conn, model, refCol, WhereClause{
		Where:     ``,
		Arguments: []any{},
	})
}

func initDB() {
	var dbDIR = "./bin"
	var dbPath = fmt.Sprintf("%v/test.db", dbDIR)
	err := os.MkdirAll(dbDIR, 0755)
	checkError(err)
	dbInstance, err = NewDatabase(dbPath)
	checkError(err)
	_, err = dbInstance.Exec(UserSQLModel)
	checkError(err)
}

func deInitDB() {
	if dbInstance != nil {
		dbInstance.Close()
	}
}

func TestDBLite(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Tests Model Insert", func() {
		g.It("user insert", func() {
			g.Timeout(1 * time.Hour)
			initDB()
			defer deInitDB()

			var m = NewModel(1)
			m.Email = "email@db.com"
			m.Name = "model"
			m.Address = "123 db street"
			m.Address = "123 db street"
			bln, err := m.InsertOnConflictDoNothing()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()
			bln, err = m.InsertOnConflictDoNothing()
			g.Assert(bln).IsFalse()
			g.Assert(err).IsNil()
		})

		g.It("user upsert", func() {
			g.Timeout(1 * time.Hour)
			initDB()
			defer deInitDB()

			var m = NewModel(1)
			m.Email = "email@db.com"
			m.Name = "model"
			m.Address = "123 db street"
			m.Address = "123 db street"
			bln, err := m.Upsert()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()
			bln, err = m.Upsert()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()
		})

		g.It("user upsert with sql args", func() {
			g.Timeout(1 * time.Hour)
			initDB()
			defer deInitDB()

			var m = NewModel(1)
			m.Email = "email@db.com"
			m.Name = "model"
			m.Address = "123 db street"
			m.Address = "123 db street"
			bln, err := m.InsertWithArgs()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()
			bln, err = m.InsertWithArgs()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()
		})

		g.It("user update with args", func() {})
		g.It("user delete with args", func() {})
		g.It("user count with args", func() {})

	})

}
