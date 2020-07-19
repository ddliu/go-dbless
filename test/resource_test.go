// Copyright 2019 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package test

import (
	"fmt"
	"testing"

	"github.com/ddliu/go-dbless"
	"github.com/spf13/cast"
)

func doTestResource(t *testing.T, setup setupDB) {
	db := setup()
	resource := dbless.Resource{
		Name:      "test",
		DB:        db,
		Timestamp: true,
	}

	// save
	id, err := resource.Save(dbless.Record{
		"name": "record1",
	})

	if err != nil || cast.ToInt(id) != 1 {
		t.Fatal("save record failed", err, id)
	}

	// get
	record, err := resource.Get(1)
	if err != nil || record["name"] != "record1" {
		t.Fatal("get record failed" + err.Error())
	}

	for i := 2; i <= 100; i++ {
		resource.Save(dbless.Record{
			"name": fmt.Sprintf("record%d", i),
		})
	}
	// list
	rst, err := resource.List(dbless.ListInput{
		Pagination: dbless.Pagination{
			PageSize: 10,
			Page:     2,
		},
	})

	if err != nil || rst.Pagination.Total != 100 || rst.PageTotal != 10 || rst.Page != 2 || rst.PageSize != 10 {
		t.Fatal("list failed", err)
	}

	// filter
	rst, err = resource.List(dbless.ListInput{
		Filter: dbless.Filter{
			"name": "record9",
		},
		Pagination: dbless.Pagination{
			PageSize: 10,
			Page:     1,
		},
	})

	if err != nil || rst.Pagination.Total != 1 || rst.List[0]["name"] != "record9" {
		t.Fatal("filter failed")
	}

	// delete
	err = resource.Delete(3)
	if err != nil {
		t.Fatal("delete error", err)
	}

	_, err = resource.Get(3)
	if !dbless.IsRecordNotFound(err) {
		t.Fatal("delete failed")
	}

	// update by name
	_, err = resource.Save(dbless.Record{
		"name": "updated",
	}, dbless.Filter{
		"name": "record2",
	})

	if err != nil {
		t.Fatal(err)
	}

	r, err := resource.Get(dbless.Filter{
		"name": "updated",
	})

	if err != nil || r.ID() != "2" {
		t.Fatal("Save & Get by name failed")
	}
}
