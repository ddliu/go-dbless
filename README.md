# go-dbless

[![Travis](https://img.shields.io/travis/ddliu/go-dbless.svg?style=flat-square)](https://travis-ci.org/ddliu/go-dbless)
[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/ddliu/go-dbless)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ddliu/go-dbless)](https://goreportcard.com/report/github.com/ddliu/go-dbless)
[![cover.run](https://cover.run/go/github.com/ddliu/go-dbless.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com%2Fddliu%2Fgo-dbless)

Database library with less complexity and less modeling.

## Install

```
go get -u github.com/ddliu/go-dbless
```

## Usage

### Basic

- DBInsert
- DBUpdate
- DBDelete
- DBGetRows
- DBGetRow
- DBGetScalar

### Schema

- List databases
- List tables
- List columns

### Resource

```go
package main
import github.com/ddliu/go-dbless

func main() {
    db := somedb()

    resource := dbless.Resource{
        Name: "user",
        Timestamp: true,
        DB: db,
    }

    id, err := resource.Save(dbless.Record{
        "username": "ddliu",
        "password": "123456",
    })

    user, err := record.Get(id)

    listing := resource.List(dbless.ListInput{
        Pagination: dbless.Pagination {
            PageSize: 20,
            Page: 1,
        },
    })
}
```