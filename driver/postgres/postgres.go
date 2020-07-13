package postgres

import (
	"database/sql"

	"github.com/ddliu/go-dbless"
)

type PostgresDriver struct {
}

func (m *PostgresDriver) Name() string {
	return "postgres"
}

func (m *PostgresDriver) QuoteIdentifier(v string) string {
	return "\"" + v + "\""
}

func (m *PostgresDriver) ListDatabases(db *dbless.DB) ([]string, error) {
	return dbless.ListDatabaseByQuery(db,
		"SELECT datname FROM pg_catalog.pg_database",
		"datname",
		[]string{
			"postgres", "template1", "template0",
		})
}

func (m *PostgresDriver) ListTables(db *dbless.DB, dbname string) ([]string, error) {
	return dbless.ListTableByQuery(
		db,
		"SELECT tablename FROM pg_catalog.pg_tables where schemaname = ?",
		"name",
		dbname)
}

func (m *PostgresDriver) ListColumns(db *dbless.DB, dbname string, tablename string) ([]*sql.ColumnType, error) {
	return dbless.ListColumnsByQuery(db, dbname, tablename)
}

func init() {
	dbless.RegisterDriver(&PostgresDriver{})
}
