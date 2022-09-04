package sqlite

import (
	"database/sql"

	"github.com/ddliu/go-dbless"
)

type Sqlite3Driver struct {
}

func (m *Sqlite3Driver) Name() string {
	return "sqlite3"
}

func (m *Sqlite3Driver) QuoteIdentifier(v string) string {
	return "\"" + v + "\""
}

func (m *Sqlite3Driver) ListDatabases(db *dbless.DB) ([]string, error) {
	return nil, dbless.DriverNotImplementedError{}
}

func (m *Sqlite3Driver) ListTables(db *dbless.DB, dbname string) ([]string, error) {
	return dbless.ListTableByQuery(
		db,
		"SELECT name FROM sqlite_master WHERE type ='table' AND name NOT LIKE 'sqlite_%'",
		"name")
}

func (m *Sqlite3Driver) ListColumns(db *dbless.DB, dbname string, tablename string) ([]*sql.ColumnType, error) {
	return dbless.ListColumnsByQuery(db, dbname, tablename)
}

func (m *Sqlite3Driver) Placeholder(values []interface{}) []string {
	var result []string
	for range values {
		result = append(result, "?")
	}

	return result
}

func (m *Sqlite3Driver) ScanReceiver(t *sql.ColumnType) (interface{}, error) {
	return nil, nil
}

func init() {
	dbless.RegisterDriver(&Sqlite3Driver{})
}
