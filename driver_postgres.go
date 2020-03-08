package dbless

import "database/sql"

type PostgresDriver struct {
}

func (m *PostgresDriver) Name() string {
	return "postgres"
}

func (m *PostgresDriver) QuoteIdentifier(v string) string {
	return "\"" + v + "\""
}

func (m *PostgresDriver) ListDatabases(db *sql.DB) ([]string, error) {
	return listDatabaseByQuery(db,
		"SELECT datname FROM pg_catalog.pg_database",
		"datname",
		[]string{
			"postgres", "template1", "template0",
		})
}

func (m *PostgresDriver) ListTables(db *sql.DB, dbname string) ([]string, error) {
	return listTableByQuery(
		db,
		"SELECT tablename FROM pg_catalog.pg_tables where schemaname = ?",
		"name",
		dbname)
}

func (m *PostgresDriver) ListColumns(db *sql.DB, dbname string, tablename string) ([]*sql.ColumnType, error) {
	return listColumnsByQuery(db, dbname, tablename)
}

func init() {
	registerDriver(&PostgresDriver{})
}
