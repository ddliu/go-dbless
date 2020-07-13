package dbless

import "database/sql"

type Sqlite3Driver struct {
}

func (m *Sqlite3Driver) Name() string {
	return "sqlite3"
}

func (m *Sqlite3Driver) QuoteIdentifier(v string) string {
	return "\"" + v + "\""
}

func (m *Sqlite3Driver) ListDatabases(db *DB) ([]string, error) {
	return nil, DriverNotImplementedError{}
}

func (m *Sqlite3Driver) ListTables(db *DB, dbname string) ([]string, error) {
	return nil, DriverNotImplementedError{}
	// return listTableByQuery(
	// 	db,
	// 	"SELECT TABLE_NAME FROM databaseName.INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_CATALOG = ?",
	// 	"TABLE_NAME",
	// 	dbname)
}

func (m *Sqlite3Driver) ListColumns(db *DB, dbname string, tablename string) ([]*sql.ColumnType, error) {
	return listColumnsByQuery(db, dbname, tablename)
}

func init() {
	registerDriver(&Sqlite3Driver{})
}
