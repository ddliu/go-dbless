package dbless

import "database/sql"

type SqlServerDriver struct {
}

func (m *SqlServerDriver) Name() string {
	return "sqlserver"
}

func (m *SqlServerDriver) QuoteIdentifier(v string) string {
	return "\"" + v + "\""
}

func (m *SqlServerDriver) ListDatabases(db *sql.DB) ([]string, error) {
	return listDatabaseByQuery(db,
		"SELECT name FROM master.sys.databases",
		"name",
		[]string{
			"master", "tempdb", "model", "msdb",
		})
}

func (m *SqlServerDriver) ListTables(db *sql.DB, dbname string) ([]string, error) {
	return listTableByQuery(
		db,
		"SELECT TABLE_NAME FROM databaseName.INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_CATALOG = ?",
		"TABLE_NAME",
		dbname)
}

func (m *SqlServerDriver) ListColumns(db *sql.DB, dbname string, tablename string) ([]*sql.ColumnType, error) {
	return listColumnsByQuery(db, dbname, tablename)
}

func init() {
	registerDriver(&SqlServerDriver{})
}
