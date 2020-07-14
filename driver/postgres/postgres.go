package postgres

import (
	"database/sql"
	"fmt"
	"time"

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

func (m *PostgresDriver) Placeholder(values []interface{}) []string {
	var result []string
	for i := range values {
		result = append(result, fmt.Sprintf("$%d", i+1))
	}

	return result
}

func (m *PostgresDriver) ScanReceiver(t *sql.ColumnType) (interface{}, error) {
	switch t.DatabaseTypeName() {
	case "INT2", "INT4", "INT8":
		var v sql.NullInt64
		return &v, nil
	case "DECIMAL", "FLOAT4", "FLOAT8":
		var v sql.NullFloat64
		return &v, nil
	case "TIMESTAMP", "DATETIME":
		var v time.Time
		return &v, nil
	}

	return nil, nil
}

func init() {
	dbless.RegisterDriver(&PostgresDriver{})
}
