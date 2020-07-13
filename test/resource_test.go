// Copyright 2019 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package test

import (
	"errors"
	"fmt"

	"github.com/ddliu/go-dbless"
)

func doTestResource(resource dbless.Resource) error {
	// save
	id, err := resource.Save(dbless.Record{
		"name": "record1",
	})

	if err != nil || id != 1 {
		return errors.New("save record failed")
	}

	// get
	record, err := resource.Get(1)
	if err != nil || record["name"] != "record1" {
		return errors.New("get record failed" + err.Error())
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
		return errors.New("list failed")
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
		return errors.New("filter failed")
	}

	// delete
	err = resource.Delete(3)
	if err != nil {
		return err
	}

	_, err = resource.Get(3)
	if !dbless.IsRecordNotFound(err) {
		return errors.New("delete failed")
	}

	// update by name
	_, err = resource.Save(dbless.Record{
		"name": "updated",
	}, dbless.Filter{
		"name": "record2",
	})

	if err != nil {
		return err
	}

	r, err := resource.Get(dbless.Filter{
		"name": "updated",
	})

	if err != nil || r.ID() != 2 {
		return errors.New("Save & Get by name failed")
	}

	return nil
}
