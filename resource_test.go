// Copyright 2019 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package dbless

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func newDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
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

	return db
}

func newResource() Resource {
	db := newDB()
	resource := Resource{
		Name:      "test",
		Timestamp: true,
		DB:        db,
	}

	return resource
}

func TestResource(t *testing.T) {
	resource := newResource()

	// save
	id, err := resource.Save(Record{
		"name": "record1",
	})

	if err != nil || id != 1 {
		t.Error("save record failed")
	}

	// get
	record, err := resource.Get(1)
	if err != nil || record["name"] != "record1" {
		t.Error("get record failed")
	}

	for i := 2; i <= 100; i++ {
		resource.Save(Record{
			"name": fmt.Sprintf("record%d", i),
		})
	}
	// list
	rst, err := resource.List(ListInput{
		Pagination: Pagination{
			PageSize: 10,
			Page:     2,
		},
	})

	if err != nil || rst.Pagination.Total != 100 || rst.PageTotal != 10 || rst.Page != 2 || rst.PageSize != 10 {
		t.Error("list failed")
	}

	// delete
	err = resource.Delete(3)
	if err != nil {
		t.Error(err)
	}

	_, err = resource.Get(3)
	if !IsRecordNotFound(err) {
		t.Error("delete failed")
	}
}
