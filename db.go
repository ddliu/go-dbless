package dbless

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

type DB struct {
	DB     *sql.DB
	Driver Driver
	Debug  bool
}

func Open(driverName, dsn string) (*DB, error) {
	db, err := sql.Open(driverName, dsn)
	driver, _ := getDriver(driverName)
	if err != nil {
		return nil, err
	}

	return &DB{
		DB:     db,
		Driver: driver,
	}, nil
}

func New(driverName string, db *sql.DB) *DB {
	rst := &DB{
		DB: db,
	}

	d, ok := getDriver(driverName)
	if ok {
		rst.Driver = d
	}

	return rst
}

func (db *DB) debug(message string, args ...interface{}) {
	if db.Debug {
		log.Print(append([]interface{}{message}, args...))
	}
}

func (db *DB) ListDatabases() ([]string, error) {
	if db.Driver == nil {
		return nil, NoDriverError{}
	}

	return db.Driver.ListDatabases(db)
}

func (db *DB) ListTables(dbname string) ([]string, error) {
	if db.Driver == nil {
		return nil, NoDriverError{}
	}

	return db.Driver.ListTables(db, dbname)
}

func (db *DB) ListColumns(dbname, tablename string) ([]*sql.ColumnType, error) {
	if db.Driver == nil {
		return nil, NoDriverError{}
	}

	return db.Driver.ListColumns(db, dbname, tablename)
}

func (db *DB) Scalar(sql string, args ...interface{}) (uint64, error) {
	db.debug(sql, args...)
	row := db.DB.QueryRow(sql, args...)

	var count uint64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (db *DB) GetPaging(pageSize uint, page uint, sql string, args ...interface{}) (*PagedRows, error) {
	totalSql := "select count(*) from (" + sql + ") __t"
	total, err := db.Scalar(totalSql, args...)
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
		rowsSql := sql + " " + fmt.Sprintf("limit %d offset %d", limit, offset)
		rows, err := db.GetRows(rowsSql, args...)
		if err != nil {
			return nil, err
		}

		result.List = rows
	}

	return result, nil
}

func (db *DB) GetRows(sqlStr string, args ...interface{}) ([]Record, error) {
	db.debug(sqlStr, args...)
	rows, err := db.DB.Query(sqlStr, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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
		r, err := db.Driver.ScanReceiver(t)
		if err != nil {
			return nil, err
		}

		if r != nil {
			receiver[index] = r
			continue
		}

		switch t.DatabaseTypeName() {
		case "INT", "BIGINT", "INTEGER", "TINYINT":
			var v sql.NullInt64
			receiver[index] = &v
		case "DECIMAL":
			var v sql.NullFloat64
			receiver[index] = &v
		case "TIMESTAMP", "DATETIME":
			var v sql.NullTime
			receiver[index] = &v
		default:
			var v sql.NullString
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

	for rows.Next() {
		err := rows.Scan(receiver...)
		if err != nil {
			return nil, err
		}

		item := make(map[string]interface{})
		for i, v := range receiver {
			if vv, ok := v.(driver.Valuer); ok {
				vvv, err := vv.Value()
				if err != nil {
					return nil, err
				}

				item[columns[i]] = vvv
			} else {
				item[columns[i]] = reflect.ValueOf(v).Elem().Interface()
			}
		}

		result = append(result, Record(item))
	}

	return result, nil
}

func (db *DB) GetRow(sql string, args ...interface{}) (Record, error) {
	rows, err := db.GetRows(sql, args...)
	if err != nil {
		return nil, err
	}

	if len(rows) > 0 {
		return rows[0], nil
	}

	return nil, RecordNotFoundError{}
}

func (db *DB) Insert(table string, row map[string]interface{}) error {
	var columns []string
	var values []interface{}
	var placeHolders []string

	for k, v := range row {
		columns = append(columns, db.Driver.QuoteIdentifier(k))
		values = append(values, v)
		placeHolders = db.Driver.Placeholder(values)
	}
	sql := "INSERT INTO " + db.Driver.QuoteIdentifier(table) + " (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(placeHolders, ", ") + ")"

	db.debug(sql, values...)
	_, err := db.DB.Exec(sql, values...)
	return err
}

func (db *DB) InsertGetID(table string, row map[string]interface{}) (string, error) {
	var columns []string
	var values []interface{}
	var placeHolders []string

	for k, v := range row {
		columns = append(columns, db.Driver.QuoteIdentifier(k))
		values = append(values, v)
		placeHolders = db.Driver.Placeholder(values)
	}
	sql := "INSERT INTO " + db.Driver.QuoteIdentifier(table) + " (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(placeHolders, ", ") + ")"

	if db.Driver.Name() == "postgres" {
		sql += " returning id"
		r, err := db.GetRow(sql, values...)
		if err != nil {
			return "", err
		}

		return cast.ToString(r["id"]), nil
	} else {
		db.debug(sql, values...)
		rst, err := db.DB.Exec(sql, values...)
		if err != nil {
			return "", err
		}
		id, err := rst.LastInsertId()
		return cast.ToString(id), err
	}
}

func (db *DB) Update(table string, row map[string]interface{}, where string, args ...interface{}) (uint64, error) {
	var columns []string
	var values []interface{}

	for k, v := range row {
		columns = append(columns, k)
		values = append(values, v)
	}

	placeholders := db.Driver.Placeholder(append(args, values...))[len(args):]
	for i, placeholder := range placeholders {
		columns[i] = db.Driver.QuoteIdentifier(columns[i]) + " = " + placeholder
	}

	sql := "UPDATE " + db.Driver.QuoteIdentifier(table) + " SET " + strings.Join(columns, ", ")
	if where != "" {
		sql += " WHERE " + where
		// fix param order, where params comes first
		if db.Driver.Name() == "postgres" {
			values = append(args, values...)
		} else {
			values = append(values, args...)
		}
	}

	db.debug(sql, values...)
	rst, err := db.DB.Exec(sql, values...)
	if err != nil {
		return 0, err
	}

	affected, err := rst.RowsAffected()

	return uint64(affected), err
}

func (db *DB) Delete(table string, where string, args ...interface{}) (uint64, error) {
	sql := "DELETE FROM " + db.Driver.QuoteIdentifier(table) + " WHERE " + where
	db.debug(sql, args...)
	rst, err := db.DB.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	affected, err := rst.RowsAffected()

	return uint64(affected), err
}
