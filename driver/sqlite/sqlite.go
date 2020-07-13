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
	return nil, dbless.DriverNotImplementedError{}
	// return listTableByQuery(
	// 	db,
	// 	"SELECT TABLE_NAME FROM databaseName.INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_CATALOG = ?",
	// 	"TABLE_NAME",
	// 	dbname)
}

func (m *Sqlite3Driver) ListColumns(db *dbless.DB, dbname string, tablename string) ([]*sql.ColumnType, error) {
	return dbless.ListColumnsByQuery(db, dbname, tablename)
}

func init() {
	dbless.RegisterDriver(&Sqlite3Driver{})
}
