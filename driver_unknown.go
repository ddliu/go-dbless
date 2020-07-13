package dbless

import (
	"database/sql"
)

type UnknownDriver struct{}

func (m *UnknownDriver) Name() string {
	return "known"
}

func (m *UnknownDriver) QuoteIdentifier(v string) string {
	return v
}

func (m *UnknownDriver) ListDatabases(db *DB) ([]string, error) {
	return nil, DriverNotImplementedError{}
}

func (m *UnknownDriver) ListTables(db *DB, dbname string) ([]string, error) {
	return nil, DriverNotImplementedError{}
}

func (m *UnknownDriver) ListColumns(db *DB, dbname string, tablename string) ([]*sql.ColumnType, error) {
	return listColumnsByQuery(db, dbname, tablename)
}

func init() {
	registerDriver(&UnknownDriver{})
}
