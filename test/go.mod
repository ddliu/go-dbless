module github.com/ddliu/go-dbless/test

go 1.14

require (
	github.com/ddliu/go-dbless v0.2.1
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jackc/pgx/v4 v4.9.0
	github.com/lib/pq v1.7.0
	github.com/mattn/go-sqlite3 v1.14.0
	github.com/spf13/cast v1.3.1
	google.golang.org/appengine v1.6.6 // indirect
)

replace github.com/ddliu/go-dbless v0.2.1 => ../
