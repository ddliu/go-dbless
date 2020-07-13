// Copyright 2019 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package dbless

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/spf13/cast"

	_ "github.com/mattn/go-sqlite3"
)

func newDB() *DB {
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

	return New("sqlite3", db)
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
		t.Error("get record failed", err)
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

	// filter
	rst, err = resource.List(ListInput{
		Filter: Filter{
			"name": "record9",
		},
		Pagination: Pagination{
			PageSize: 10,
			Page:     1,
		},
	})

	if err != nil || rst.Pagination.Total != 1 || rst.List[0]["name"] != "record9" {
		t.Error("filter failed")
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

	// update by name
	_, err = resource.Save(Record{
		"name": "updated",
	}, Filter{
		"name": "record2",
	})

	if err != nil {
		t.Error(err)
	}

	r, err := resource.Get(Filter{
		"name": "updated",
	})

	if err != nil || r.ID() != 2 {
		t.Error("Save & Get by name failed")
	}
}

func TestUtils(t *testing.T) {
	db := newResource().DB

	for i := 1; i <= 100; i++ {
		id, err := db.Insert("test", map[string]interface{}{
			"name":       fmt.Sprintf("record%d", i),
			"created_at": "",
			"updated_at": "",
		})

		if err != nil || cast.ToInt(id) != i {
			t.Error("insert error: ", err)
		}
	}

	result, err := db.GetPaging(10, 2, "select * from test order by id asc")
	if err != nil || result.Pagination.PageTotal != 10 || len(result.List) != 10 || cast.ToString(result.List[0]["name"]) != "record11" {
		t.Error("paging error:", result, result.Pagination.PageTotal, err)
	}

	// unmarshal
	type Test struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}

	m := Test{}
	err = result.List[0].Unmarshal(&m)

	if err != nil || m.Name != "record11" {
		t.Error("unmarshal error", err, m, result.List[0])
	}
}

func TestNull(t *testing.T) {
	db := newResource().DB
	row, err := db.GetRow("select 1 as num, null as name")
	if err != nil || row["name"] != nil {
		t.Errorf("nullable error")
	}
}
