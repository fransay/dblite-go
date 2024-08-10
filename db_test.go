package dblite

import (
	"fmt"
	"github.com/franela/goblin"
	"os"
	"testing"
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

type Model struct {
	Id      int64  `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Active  int    `json:"active"`
}

func NewModel() *Model {
	return &Model{Id: -1}
}

func (model *Model) New() *Model {
	return NewModel()
}

func (model *Model) Clone() *Model {
	var o = *model
	return &o
}

func (model *Model) TableName() string {
	return "model"
}

func (model *Model) Insert() (bool, error) {
	return Insert(Instance.Conn, model, []string{
		`email`, `name`, `address`,
	}, On{})
}

func initDB() {
	var dbDIR = "./bin"
	var dbPath = fmt.Sprintf("%v/test.db", dbDIR)
	err := os.MkdirAll(dbDIR, 0755)
	if err != nil {
		panic(err)
	}
	Init(dbPath)

	_, err = Instance.Exec(UserSQLModel)
	if err != nil {
		panic(err)
	}
}

func deInitDB() {
	DeInit()
}

func TestDBLite(t *testing.T) {
	g := goblin.Goblin(t)

	initDB()
	defer deInitDB()

	g.Describe("Tests Model Insert", func() {
		g.It("user insert", func() {
			var m = NewModel()
			m.Email = "email@db.com"
			m.Name = "model"
			m.Address = "123 db street"
			m.Address = "123 db street"
			bln, err := m.Insert()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()
		})

	})

}
