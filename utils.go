// Copyright 2019 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package dbless

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type Record map[string]interface{}

func (r Record) Unmarshal(input interface{}) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, input)
}

func DBScalar(db *sql.DB, sql string, args ...interface{}) (uint64, error) {
	row := db.QueryRow(sql, args...)

	var count uint64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

type PagedRows struct {
	Pagination
	List []Record `json:"list"`
}

func DBGetPaging(db *sql.DB, pageSize uint, page uint, sql string, args ...interface{}) (*PagedRows, error) {
	totalSql := "select count(*) from (" + sql + ") __t"
	total, err := DBScalar(db, totalSql, args...)
	if err != nil {
		return nil, err
	}

	pagination := Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	pagination.Valid()
	pagination.SetTotal(cast.ToUint(total))

	result := &PagedRows{
		Pagination: pagination,
	}
	if total > 0 {
		offset, limit := pagination.GetOffsetLimit()
		rowsSql := sql + " " + fmt.Sprintf("limit %d, %d", offset, limit)
		rows, err := DBGetRows(db, rowsSql, args...)
		if err != nil {
			return nil, err
		}

		result.List = rows
	}

	return result, nil
}

func DBGetRows(db *sql.DB, sql string, args ...interface{}) ([]Record, error) {
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	columnLength := len(columns)
	receiver := make([]interface{}, columnLength)
	for index, _ := range receiver {
		t := columnTypes[index]
		switch t.DatabaseTypeName() {
		case "INT", "BIGINT", "INTEGER":
			var v int64
			receiver[index] = &v
		case "DECIMAL":
			var v float64
			receiver[index] = &v
		case "TIMESTAMP", "DATETIME":
			var v time.Time
			receiver[index] = &v
		default:
			var v string
			receiver[index] = &v
		}

		// switch reflect.New(t.ScanType()).Interface().(type) {
		// case uint, uint8, uint32, uint64, int, int8, int32, int64:
		// 	var v int64
		// 	receiver[index] = &v
		// case float32, float64:
		// 	var v float64
		// 	receiver[index] = &v
		// case time.Time:
		// 	var v time.Time
		// 	receiver[index] = &v
		// default:
		// 	var v string
		// 	receiver[index] = &v
		// }
	}

	var result []Record

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(receiver...)
		if err != nil {
			return nil, err
		}

		item := make(map[string]interface{})
		for i, v := range receiver {
			item[columns[i]] = reflect.ValueOf(v).Elem().Interface()
		}

		result = append(result, Record(item))
	}

	return result, nil
}

func DBGetRow(db *sql.DB, sql string, args ...interface{}) (Record, error) {
	rows, err := DBGetRows(db, sql, args...)
	if err != nil {
		return nil, err
	}

	if len(rows) > 0 {
		return rows[0], nil
	}

	return nil, RecordNotFoundError{}
}

func DBInsert(db *sql.DB, table string, row map[string]interface{}) (uint64, error) {
	var columns []string
	var values []interface{}
	var placeHolders []string

	for k, v := range row {
		columns = append(columns, quote(k))
		values = append(values, v)
		placeHolders = append(placeHolders, "? ")
	}
	sql := "INSERT INTO " + quote(table) + " (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(placeHolders, ", ") + ")"

	rst, err := db.Exec(sql, values...)
	if err != nil {
		return 0, err
	}

	id, err := rst.LastInsertId()
	return uint64(id), err
}

func DBUpdate(db *sql.DB, table string, row map[string]interface{}, where string, args ...interface{}) (uint64, error) {
	var columns []string
	var values []interface{}

	for k, v := range row {
		columns = append(columns, quote(k)+" = ?")
		values = append(values, v)
	}

	sql := "UPDATE " + quote(table) + " SET " + strings.Join(columns, ", ")
	if where != "" {
		sql += " WHERE " + where
		values = append(values, args...)
	}

	rst, err := db.Exec(sql, values...)
	if err != nil {
		return 0, err
	}

	affected, err := rst.RowsAffected()

	return uint64(affected), err
}

func DBDelete(db *sql.DB, table string, where string, args ...interface{}) (uint64, error) {
	sql := "DELETE FROM " + quote(table) + " WHERE " + where
	rst, err := db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	affected, err := rst.RowsAffected()

	return uint64(affected), err
}

func quote(identifier string) string {
	return "`" + identifier + "`"
}
