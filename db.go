package dbless

import (
	"database/sql"
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

	return db.Driver.ListDatabases(db.DB)
}

func (db *DB) ListTables(dbname string) ([]string, error) {
	if db.Driver == nil {
		return nil, NoDriverError{}
	}

	return db.Driver.ListTables(db.DB, dbname)
}

func (db *DB) ListColumns(dbname, tablename string) ([]*sql.ColumnType, error) {
	if db.Driver == nil {
		return nil, NoDriverError{}
	}

	return db.Driver.ListColumns(db.DB, dbname, tablename)
}

func (db *DB) Scalar(sql string, args ...interface{}) (uint64, error) {
	return DBScalar(db.DB, sql, args...)
}

func (db *DB) GetPaging(pageSize uint, page uint, sql string, args ...interface{}) (*PagedRows, error) {
	return DBGetPaging(db.DB, pageSize, page, sql, args...)
}

func (db *DB) GetRows(sqlStr string, args ...interface{}) ([]Record, error) {
	return DBGetRows(db.DB, sqlStr, args...)
}

func (db *DB) GetRow(sql string, args ...interface{}) (Record, error) {
	return DBGetRow(db.DB, sql, args...)
}

func (db *DB) Insert(table string, row map[string]interface{}) (uint64, error) {
	return DBInsert(db.DB, table, row)
}

func (db *DB) Update(table string, row map[string]interface{}, where string, args ...interface{}) (uint64, error) {
	return DBUpdate(db.DB, table, row, where, args...)
}

func (db *DB) Delete(table string, where string, args ...interface{}) (uint64, error) {
	return DBDelete(db.DB, table, where, args...)
}
