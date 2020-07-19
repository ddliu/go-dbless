package test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/ddliu/go-dbless"
	_ "github.com/ddliu/go-dbless/driver/postgres"
	_ "github.com/lib/pq"
)

func newPostgresDB(url string) *dbless.DB {
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	DROP table if exists test;
	CREATE table test (
		id  serial PRIMARY KEY,
		name varchar(100) NOT NULL,
		created_at timestamp NOT NULL,
		updated_at timestamp NOT NULL
	);
	`)

	if err != nil {
		panic(err)
	}

	return dbless.New("postgres", db)
}

func TestPostgres(t *testing.T) {
	dbUrl := os.Getenv("POSTGRES_DATABASE_URL")
	if dbUrl != "" {
		doTestDB(t, func() *dbless.DB {
			return newPostgresDB(dbUrl)
		})
	}
}
