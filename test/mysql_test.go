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

	db.Exec(`
	DROP table if exists test;
	`)
	_, err = db.Exec(`
	CREATE table test (
		id  int NOT NULL AUTO_INCREMENT,
		name varchar(100) NOT NULL,
		created_at timestamp NULL,
		updated_at timestamp NULL,
		PRIMARY KEY (id)
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
		doTestDB(t, func() *dbless.DB {
			return newMysqlDB(dbUrl)
		})
	}
}
