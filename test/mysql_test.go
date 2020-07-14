package test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/ddliu/go-dbless"
	_ "github.com/ddliu/go-dbless/driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func newMysqlDB(url string) *dbless.DB {
	db, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	DROP table if exists test;
	CREATE table test (
		id  int PRIMARY KEY,
		name varchar(100) NOT NULL,
		created_at timestamp NOT NULL,
		updated_at timestamp NOT NULL
	);
	`)

	if err != nil {
		panic(err)
	}

	return dbless.New("mysql", db)
}

func TestMysql(t *testing.T) {
	dbUrl := os.Getenv("MYSQL_DATABASE_URL")
	if dbUrl != "" {
		if err := doTestDB(newMysqlDB(dbUrl)); err != nil {
			t.Error(err)
		}
	}
}
