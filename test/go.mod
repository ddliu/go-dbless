module github.com/ddliu/go-dbless/test

go 1.14

require (
	github.com/ddliu/go-dbless v0.2.1
	github.com/ddliu/go-dbless/driver/mysql v0.2.1
	github.com/ddliu/go-dbless/driver/postgres v0.2.1
	github.com/ddliu/go-dbless/driver/sqlite v0.2.1
	github.com/go-sql-driver/mysql v1.4.1
	github.com/lib/pq v1.7.0
	github.com/mattn/go-sqlite3 v1.14.0
	github.com/spf13/cast v1.3.1
	google.golang.org/appengine v1.6.6 // indirect
)

replace (
	github.com/ddliu/go-dbless v0.2.1 => ../
	github.com/ddliu/go-dbless/driver/mysql v0.2.1 => ../driver/mysql
	github.com/ddliu/go-dbless/driver/postgres v0.2.1 => ../driver/postgres
	github.com/ddliu/go-dbless/driver/sqlite v0.2.1 => ../driver/sqlite
)
