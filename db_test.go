package dblite

import (
	"fmt"
	"github.com/franela/goblin"
	util "github.com/intdxdt/goreflect"
	"os"
	"testing"
)

// DROP TABLE IF EXISTS model;
const UserSQLModel = `
CREATE TABLE IF NOT EXISTS model (
	id            		 INTEGER NOT NULL PRIMARY KEY,
	email         		 TEXT NOT NULL UNIQUE,
	name          		 TEXT DEFAULT '',
	address   			 TEXT DEFAULT '',
	active        		 INTEGER DEFAULT 0
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

func (user *Model) Clone() *Model {
	var o = *user
	return &o
}

func (user *Model) TableName() string {
	return "user"
}

func (user *Model) FieldRefMap() map[string]any {
	var fields = user.Fields()
	var refs, err = util.GetFieldReferences(user, fields)
	if err != nil {
		panic(err)
	}
	var dict = make(map[string]any, len(fields))
	for i, field := range fields {
		dict[field] = refs[i]
	}
	return dict
}

func (user *Model) FilterFieldReferences(fields []string) ([]string, []any, error) {
	return util.FilterFieldReferences(fields, user.FieldRefMap())
}

func (user *Model) Fields() []string {
	var fields, err = util.GetJSONTaggedFields(user)
	if err != nil {
		panic(err)
	}
	return fields
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
			g.Assert(12).Equal(12)
		})

	})

}
