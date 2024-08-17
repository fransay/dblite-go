package dblite

import (
	"database/sql"
	"fmt"
	ref "github.com/intdxdt/goreflect"
)

func Insert[T ITable[T]](conn *sql.DB, model T, insertCols []string, on On) (bool, error) {
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

	var dict = KeysToMap(insertCols, true)

	for i, field := range fields {
		if !(dict[field]) {
			continue
		}
		cols = append(cols, field)
		values = append(values, colRefs[i])
	}

	var columns = ColumnNames(cols)
	var holders = ColumnPlaceholders(cols)

	var sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v);`, model.TableName(), columns, holders)

	if len(on.On) > 0 {
		sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v)
		ON %v;`, model.TableName(), columns, holders, on.On)
		for _, v := range on.Arguments {
			values = append(values, v)
		}
	}

	res, err := Exec(conn, sqlStatement, values...)
	if err != nil {
		return false, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func InsertMany[T ITable[T]](conn *sql.DB, models []T, insertCols []string, on On) (error, error) {
	if len(models) == 0 {
		return nil, nil
	}

	var getColumnsValues = func(model T) ([]string, []any, error) {
		var fields, err = ref.Fields(model)
		if err != nil {
			return nil, nil, err
		}

		fields, colRefs, err := ref.FilterFieldReferences(fields, model)
		if err != nil {
			return nil, nil, err
		}
		var cols = make([]string, 0, len(fields))
		var values = make([]any, 0, len(fields))

		var dict = KeysToMap(insertCols, true)

		for i, field := range fields {
			if !dict[field] {
				continue
			}
			cols = append(cols, field)
			values = append(values, colRefs[i])
		}
		return cols, values, nil
	}
	var model = models[0]
	var cols, _, err = getColumnsValues(model)
	if err != nil {
		return err, nil
	}

	var columns = ColumnNames(cols)
	var holders = ColumnPlaceholders(cols)

	var sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v);`, model.TableName(), columns, holders)

	if len(on.On) > 0 {
		sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v)
		ON %v;`, model.TableName(), columns, holders, on.On)
	}

	var records = make([][]any, 0, len(models))
	for _, model = range models {
		_, values, err := getColumnsValues(model)
		if err != nil {
			return err, nil
		}
		if len(on.On) > 0 {
			for _, v := range on.Arguments {
				values = append(values, v)
			}
		}
		records = append(records, values)
	}

	return ExecMany(conn, sqlStatement, records)
}
