package mysql

import (
	"database/sql"

	"github.com/ddliu/go-dbless"
)

type MysqlDriver struct{}

func (m *MysqlDriver) Name() string {
	return "mysql"
}

func (m *MysqlDriver) QuoteIdentifier(v string) string {
	return "`" + v + "`"
}

func (m *MysqlDriver) ListDatabases(db *dbless.DB) ([]string, error) {
	return dbless.ListDatabaseByQuery(db,
		"show databases",
		"Database",
		[]string{
			"information_schema", "performance_schema", "mysql", "sys",
		})
}

func (m *MysqlDriver) ListTables(db *dbless.DB, dbname string) ([]string, error) {
	return dbless.ListTableByQuery(
		db,
		"SELECT TABLE_NAME AS name FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ?",
		"name",
		dbname)
}

func (m *MysqlDriver) ListColumns(db *dbless.DB, dbname string, tablename string) ([]*sql.ColumnType, error) {
	return dbless.ListColumnsByQuery(db, dbname, tablename)
}

func init() {
	dbless.RegisterDriver(&MysqlDriver{})
}
