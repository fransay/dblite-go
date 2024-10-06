package dblite

import (
	"database/sql"
	"errors"
	"fmt"
	ref "github.com/intdxdt/goreflect"
)

func Query(conn *sql.DB, query string, args ...any) (*sql.Rows, error) {
	return conn.Query(query, args...)
}

func QueryModel[T ITable[T]](conn *sql.DB, model T, where ...WhereClause) (T, error) {
	var fields, err = ref.Fields(model)
	if err != nil {
		return model.New(), err
	}
	return QueryModelByColumnNames(conn, model, fields, where...)
}

func QueryModelByColumnNames[T ITable[T]](conn *sql.DB, model T, fieldNames []string, where ...WhereClause) (T, error) {
	var tableName = model.TableName()
	var cols, colRefs, err = ref.FilterFieldReferences(fieldNames, model)
	if err != nil {
		return model, err
	}
	var fields = ColumnNames(cols)

	var args = make([]any, 0)
	var sqlStatement = fmt.Sprintf("SELECT %v FROM %v LIMIT 1;", fields, tableName)
	if len(where) > 0 {
		var wc = where[0]
		args = wc.Arguments
		sqlStatement = fmt.Sprintf("SELECT %v FROM %v WHERE %v LIMIT 1;", fields, tableName, wc.Where)
	}

	rows, err := Query(conn, sqlStatement, args...)
	if err != nil {
		return model, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(colRefs...)
		if err != nil {
			return model, err
		}
		break
	}

	if rows.Err() != nil {
		return model, rows.Err()
	}
	return model, nil
}

func QueryModels[T ITable[T]](conn *sql.DB, model T, where ...WhereClause) ([]T, error) {
	var fields, err = ref.Fields(model)
	if err != nil {
		return []T{}, err
	}
	return QueriesByColumnNames(conn, model, fields, where...)
}

func QueriesByColumnNames[T ITable[T]](conn *sql.DB, model T, fieldNames []string, where ...WhereClause) ([]T, error) {
	var results = make([]T, 0)
	var tableName = model.TableName()
	var cols, colRefs, err = ref.FilterFieldReferences(fieldNames, model)
	if err != nil {
		return nil, err
	}
	var fields = ColumnNames(cols)

	var args = make([]any, 0)
	var sqlStatement = fmt.Sprintf("SELECT %v FROM %v;", fields, tableName)
	if len(where) > 0 {
		var wc = where[0]
		args = wc.Arguments
		if len(args) == 0 {
			return results, errors.New("invalid number arguments in where clause")
		}
		sqlStatement = fmt.Sprintf("SELECT %v FROM %v WHERE %v;", fields, tableName, wc.Where)
	}

	rows, err := Query(conn, sqlStatement, args...)
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(colRefs...)
		if err != nil {
			return results, err
		}
		results = append(results, model.Clone())
	}

	if rows.Err() != nil {
		return results, rows.Err()
	}
	return results, nil
}
