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

func (m *MysqlDriver) Placeholder(values []interface{}) []string {
	var result []string
	for range values {
		result = append(result, "?")
	}

	return result
}

func (m *MysqlDriver) ScanReceiver(t *sql.ColumnType) (interface{}, error) {
	switch t.DatabaseTypeName() {
	case "TIMESTAMP", "DATETIME":
		var v sql.NullTime
		return &v, nil
	}

	return nil, nil
}

func init() {
	dbless.RegisterDriver(&MysqlDriver{})
}
