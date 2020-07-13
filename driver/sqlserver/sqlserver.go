package sqlserver

import (
	"database/sql"

	"github.com/ddliu/go-dbless"
)

type SqlServerDriver struct {
}

func (m *SqlServerDriver) Name() string {
	return "sqlserver"
}

func (m *SqlServerDriver) QuoteIdentifier(v string) string {
	return "\"" + v + "\""
}

func (m *SqlServerDriver) ListDatabases(db *dbless.DB) ([]string, error) {
	return dbless.ListDatabaseByQuery(db,
		"SELECT name FROM master.sys.databases",
		"name",
		[]string{
			"master", "tempdb", "model", "msdb",
		})
}

func (m *SqlServerDriver) ListTables(db *dbless.DB, dbname string) ([]string, error) {
	return dbless.ListTableByQuery(
		db,
		"SELECT TABLE_NAME FROM databaseName.INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_CATALOG = ?",
		"TABLE_NAME",
		dbname)
}

func (m *SqlServerDriver) ListColumns(db *dbless.DB, dbname string, tablename string) ([]*sql.ColumnType, error) {
	return dbless.ListColumnsByQuery(db, dbname, tablename)
}

func init() {
	dbless.RegisterDriver(&SqlServerDriver{})
}
