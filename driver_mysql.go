package dbless

import "database/sql"

type MysqlDriver struct{}

func (m *MysqlDriver) Name() string {
	return "mysql"
}

func (m *MysqlDriver) QuoteIdentifier(v string) string {
	return "`" + v + "`"
}

func (m *MysqlDriver) ListDatabases(db *DB) ([]string, error) {
	return listDatabaseByQuery(db,
		"show databases",
		"Database",
		[]string{
			"information_schema", "performance_schema", "mysql", "sys",
		})
}

func (m *MysqlDriver) ListTables(db *DB, dbname string) ([]string, error) {
	return listTableByQuery(
		db,
		"SELECT TABLE_NAME AS name FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ?",
		"name",
		dbname)
}

func (m *MysqlDriver) ListColumns(db *DB, dbname string, tablename string) ([]*sql.ColumnType, error) {
	return listColumnsByQuery(db, dbname, tablename)
}

func init() {
	registerDriver(&MysqlDriver{})
}
