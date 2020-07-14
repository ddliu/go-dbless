package test

import (
	"errors"
	"fmt"
	"time"

	"github.com/ddliu/go-dbless"
	"github.com/spf13/cast"
)

func doTestDB(db *dbless.DB) error {
	columns, err := db.ListColumns("", "test")
	if err != nil {
		return err
	}

	if columns[1].Name() != "name" {
		return errors.New("list columns error")
	}

	if err := doTestUtils(db); err != nil {
		return err
	}

	if err := doTestNull(db); err != nil {
		return err
	}

	return nil
}

func doTestUtils(db *dbless.DB) error {
	for i := 1; i <= 100; i++ {
		id, err := db.InsertGetID("test", map[string]interface{}{
			"name":       fmt.Sprintf("record%d", i),
			"created_at": time.Now(),
			"updated_at": time.Now(),
		})

		if err != nil || cast.ToInt(id) != i {
			return errors.New("insert error: " + err.Error())
		}
	}

	result, err := db.GetPaging(10, 2, "select * from test order by id asc")
	if err != nil || result.Pagination.PageTotal != 10 || len(result.List) != 10 || cast.ToString(result.List[0]["name"]) != "record11" {
		return errors.New("paging error")
	}

	// unmarshal
	type Test struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}

	m := Test{}
	err = result.List[0].Unmarshal(&m)

	if err != nil || m.Name != "record11" {
		return errors.New("unmarshal error " + err.Error())
	}

	return nil
}

func doTestNull(db *dbless.DB) error {
	row, err := db.GetRow("select 1 as num, null as name")
	if err != nil || row["name"] != nil {
		return errors.New("nullable error")
	}

	return nil
}
