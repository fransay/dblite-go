package dblite

import (
	"database/sql"
	"fmt"
	ref "github.com/intdxdt/goreflect"
)

func UpdateByExclusion[T ITable[T]](conn *sql.DB, model T, excludeCols []string, wc WhereClause) (bool, error) {
	var fields, err = ref.Fields(model)
	if err != nil {
		return false, err
	}

	fields, _, err = ref.FilterFieldReferences(fields, model)
	if err != nil {
		return false, err
	}
	var updateCols = make([]string, 0, len(fields))
	var excludeDict = KeysToMap(excludeCols, true)

	for _, field := range fields {
		if excludeDict[field] {
			continue
		}
		updateCols = append(updateCols, field)
	}
	return Update(conn, model, updateCols, wc)
}

func Update[T ITable[T]](conn *sql.DB, model T, updateCols []string, wc WhereClause) (bool, error) {
	var fields, err = ref.Fields(model)
	if err != nil {
		return false, err
	}

	fields, colRefs, err := ref.FilterFieldReferences(fields, model)
	if err != nil {
		return false, err
	}

	var cols = make([]string, 0, len(fields))
	var values = make([]any, 0, len(fields))

	var dict = KeysToMap(updateCols, true)

	for i, field := range fields {
		if dict[field] {
			cols = append(cols, field)
			values = append(values, colRefs[i])
		}
	}

	var holders = UpdatePlaceholders(cols)
	for _, arg := range wc.Arguments {
		values = append(values, arg)
	}

	var query = fmt.Sprintf(
		`UPDATE %v SET %v WHERE %v;`,
		model.TableName(), holders, wc.Where)

	res, err := Exec(conn, query, values...)

	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return count == 1, nil
}
