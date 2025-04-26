package dblite

import (
	"database/sql"
	"fmt"
)

func Count[T ITable[T]](conn *sql.DB, model T, refCol string, wc WhereClause) (int64, error) {
	var count int64
	var query = fmt.Sprintf(`SELECT COUNT(%v) FROM %v WHERE %v;`, refCol, model.TableName(), wc.Where)
	var rows, err = Query(conn, query, wc.Arguments...)
	if err != nil {
		return count, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return count, err
		}
		break
	}
	if rows.Err() != nil {
		return count, rows.Err()
	}

	return count, nil
}
