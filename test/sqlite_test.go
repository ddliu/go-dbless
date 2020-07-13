package test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/ddliu/go-dbless"
	_ "github.com/ddliu/go-dbless/driver/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func newSqliteDB(url string) *dbless.DB {
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE table test (
		id  integer PRIMARY KEY AUTOINCREMENT,
		name text NOT NULL,
		created_at text NOT NULL,
		updated_at NOT NULL
	)
	`)

	if err != nil {
		panic(err)
	}

	return dbless.New("sqlite3", db)
}

func TestSqlite(t *testing.T) {
	dbUrl := os.Getenv("SQLITE_DATABASE_URL")
	if dbUrl != "" {
		err := doTestDB(newSqliteDB(dbUrl))
		if err != nil {
			t.Error(err)
		}
	}
}
