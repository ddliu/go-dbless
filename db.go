package dbless

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cast"
)

type DB struct {
	DB     *sql.DB
	Driver Driver
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
	stmt, err := db.DB.Prepare(sqlStr)
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
		case "INT", "BIGINT", "INTEGER", "TINYINT":
			var v sql.NullInt64
			receiver[index] = &v
		case "DECIMAL":
			var v sql.NullFloat64
			receiver[index] = &v
		case "TIMESTAMP", "DATETIME":
			var v mysql.NullTime
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

	defer rows.Close()
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

func (db *DB) Insert(table string, row map[string]interface{}) (uint64, error) {
	var columns []string
	var values []interface{}
	var placeHolders []string

	for k, v := range row {
		columns = append(columns, db.Driver.QuoteIdentifier(k))
		values = append(values, v)
		placeHolders = append(placeHolders, "? ")
	}
	sql := "INSERT INTO " + db.Driver.QuoteIdentifier(table) + " (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(placeHolders, ", ") + ")"

	rst, err := db.DB.Exec(sql, values...)
	if err != nil {
		return 0, err
	}

	id, err := rst.LastInsertId()
	return uint64(id), err
}

func (db *DB) Update(table string, row map[string]interface{}, where string, args ...interface{}) (uint64, error) {
	var columns []string
	var values []interface{}

	for k, v := range row {
		columns = append(columns, db.Driver.QuoteIdentifier(k)+" = ?")
		values = append(values, v)
	}

	sql := "UPDATE " + db.Driver.QuoteIdentifier(table) + " SET " + strings.Join(columns, ", ")
	if where != "" {
		sql += " WHERE " + where
		values = append(values, args...)
	}

	rst, err := db.DB.Exec(sql, values...)
	if err != nil {
		return 0, err
	}

	affected, err := rst.RowsAffected()

	return uint64(affected), err
}

func (db *DB) Delete(table string, where string, args ...interface{}) (uint64, error) {
	sql := "DELETE FROM " + db.Driver.QuoteIdentifier(table) + " WHERE " + where
	rst, err := db.DB.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	affected, err := rst.RowsAffected()

	return uint64(affected), err
}