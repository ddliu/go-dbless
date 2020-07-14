package dbless

import (
	"database/sql"

	"github.com/spf13/cast"
)

func ListDatabaseByQuery(db *DB, query string, name string, exclude []string) ([]string, error) {
	var result []string
	records, err := db.GetRows(query)
	if err != nil {
		return nil, err
	}

	for _, r := range records {
		name := cast.ToString(r[name])
		isExclude := false
		for _, x := range exclude {
			if x == name {
				isExclude = true
				break
			}
		}

		if isExclude {
			continue
		}

		result = append(result, name)
	}

	return result, nil
}

func ListTableByQuery(db *DB, query string, name string, params ...interface{}) ([]string, error) {
	var result []string
	records, err := db.GetRows(query, params...)
	if err != nil {
		return nil, err
	}

	for _, r := range records {
		name := cast.ToString(r[name])

		result = append(result, name)
	}

	return result, nil
}

func ListColumnsByQuery(db *DB, dbname, tableName string) ([]*sql.ColumnType, error) {
	tableName = db.Driver.QuoteIdentifier(tableName)
	if dbname != "" {
		dbname = db.Driver.QuoteIdentifier(dbname)
		tableName = dbname + "." + tableName
	}
	rows, err := db.DB.Query("select * from " + tableName + " limit 0")
	if err != nil {
		return nil, err
	}

	return rows.ColumnTypes()
}

type NoDriverError struct {
}

func (e NoDriverError) Error() string {
	return "No driver"
}

type DriverNotImplementedError struct {
}

func (e DriverNotImplementedError) Error() string {
	return "Not implemented"
}

type Driver interface {
	Name() string
	QuoteIdentifier(string) string
	Placeholder(values []interface{}) []string
	ListDatabases(db *DB) ([]string, error)
	ListTables(db *DB, dbname string) ([]string, error)
	ListColumns(db *DB, dbname string, tablename string) ([]*sql.ColumnType, error)
	ScanReceiver(t *sql.ColumnType) (interface{}, error)
}

var drivers = map[string]Driver{}

func getDriver(name string) (Driver, bool) {
	v, ok := drivers[name]
	return v, ok
}

func RegisterDriver(driver Driver) {
	drivers[driver.Name()] = driver
}
