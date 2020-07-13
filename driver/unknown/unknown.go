package unknown

import (
	"database/sql"

	"github.com/ddliu/go-dbless"
)

type UnknownDriver struct{}

func (m *UnknownDriver) Name() string {
	return "known"
}

func (m *UnknownDriver) QuoteIdentifier(v string) string {
	return v
}

func (m *UnknownDriver) ListDatabases(db *dbless.DB) ([]string, error) {
	return nil, dbless.DriverNotImplementedError{}
}

func (m *UnknownDriver) ListTables(db *dbless.DB, dbname string) ([]string, error) {
	return nil, dbless.DriverNotImplementedError{}
}

func (m *UnknownDriver) ListColumns(db *dbless.DB, dbname string, tablename string) ([]*sql.ColumnType, error) {
	return dbless.ListColumnsByQuery(db, dbname, tablename)
}

func (m *UnknownDriver) Placeholder(values []interface{}) []string {
	var result []string
	for range values {
		result = append(result, "?")
	}

	return result
}

func init() {
	dbless.RegisterDriver(&UnknownDriver{})
}
